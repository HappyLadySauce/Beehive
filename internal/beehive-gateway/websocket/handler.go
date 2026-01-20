package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/connection"
	messagepb "github.com/HappyLadySauce/Beehive/pkg/api/proto/message/v1"
	presencepb "github.com/HappyLadySauce/Beehive/pkg/api/proto/presence/v1"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handler WebSocket 处理器
type Handler struct {
	cfg            *config.Config
	grpcClient     *client.Client
	connMgr        *connection.Manager
	messageHandler *MessageHandler
	pusher         *Pusher
}

// NewHandler 创建新的 WebSocket 处理器
func NewHandler(cfg *config.Config, grpcClient *client.Client, connMgr *connection.Manager) *Handler {
	pusher := NewPusher(connMgr)
	messageHandler := NewMessageHandler(grpcClient, pusher)

	return &Handler{
		cfg:            cfg,
		grpcClient:     grpcClient,
		connMgr:        connMgr,
		messageHandler: messageHandler,
		pusher:         pusher,
	}
}

// HandleConnection 处理 WebSocket 连接
func (h *Handler) HandleConnection(c *gin.Context) {
	// 提取 Token
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		return
	}

	// 验证 Token
	ctx := c.Request.Context()
	validateResp, err := h.grpcClient.ValidateToken(ctx, token)
	if err != nil || !validateResp.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := validateResp.Id

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		klog.Errorf("WebSocket upgrade failed: %v", err)
		return
	}

	// 创建连接对象
	connection := connection.NewConnection(conn, userID)

	// 注册连接
	h.connMgr.Register(userID, connection)

	// 通知用户上线
	_, err = h.grpcClient.PresenceService().UserOnline(ctx, &presencepb.UserOnlineRequest{
		Id: userID,
	})
	if err != nil {
		klog.Errorf("Failed to notify user online: %v", err)
	}

	// 查询并推送离线消息
	go h.pushOfflineMessages(ctx, userID)

	// 启动连接处理
	go h.handleClient(connection, userID)

	klog.Infof("WebSocket connection established for user: %s", userID)
}

// handleClient 处理客户端连接
func (h *Handler) handleClient(conn *connection.Connection, userID string) {
	defer func() {
		// 清理连接
		h.connMgr.Unregister(userID)

		// 通知用户下线
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := h.grpcClient.PresenceService().UserOffline(ctx, &presencepb.UserOfflineRequest{
			Id: userID,
		})
		if err != nil {
			klog.Errorf("Failed to notify user offline: %v", err)
		}

		conn.Close()
		klog.Infof("WebSocket connection closed for user: %s", userID)
	}()

	// 设置读取超时
	conn.Conn().SetReadDeadline(time.Now().Add(h.cfg.WebSocket.ReadTimeout))
	conn.Conn().SetPongHandler(func(string) error {
		conn.UpdateLastPing()
		conn.Conn().SetReadDeadline(time.Now().Add(h.cfg.WebSocket.ReadTimeout))
		return nil
	})

	// 启动心跳
	go h.startHeartbeat(conn)

	// 启动消息发送协程
	go h.startMessageWriter(conn)

	// 读取消息
	for {
		_, message, err := conn.Conn().ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				klog.Errorf("WebSocket error: %v", err)
			}
			break
		}

		// 更新读取超时
		conn.Conn().SetReadDeadline(time.Now().Add(h.cfg.WebSocket.ReadTimeout))

		// 解析消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			klog.Errorf("Failed to parse message: %v", err)
			h.pusher.SendError(userID, "INVALID_MESSAGE_FORMAT", "Invalid JSON format")
			continue
		}

		// 处理消息
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		h.messageHandler.HandleMessage(ctx, userID, msg)
		cancel()
	}
}

// startHeartbeat 启动心跳
func (h *Handler) startHeartbeat(conn *connection.Connection) {
	ticker := time.NewTicker(h.cfg.WebSocket.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if conn.IsClosed() {
				return
			}

			// 检查最后心跳时间
			if time.Since(conn.LastPing()) > h.cfg.WebSocket.ReadTimeout*2 {
				klog.Warningf("Connection timeout for user: %s", conn.UserID())
				conn.Close()
				return
			}

			// 发送 ping
			conn.Conn().SetWriteDeadline(time.Now().Add(h.cfg.WebSocket.WriteTimeout))
			if err := conn.Conn().WriteMessage(websocket.PingMessage, nil); err != nil {
				klog.Errorf("Failed to send ping: %v", err)
				conn.Close()
				return
			}
		}
	}
}

// startMessageWriter 启动消息写入协程
func (h *Handler) startMessageWriter(conn *connection.Connection) {
	for {
		select {
		case message, ok := <-conn.SendChan():
			if !ok {
				return
			}

			conn.Conn().SetWriteDeadline(time.Now().Add(h.cfg.WebSocket.WriteTimeout))
			if err := conn.Conn().WriteMessage(websocket.TextMessage, message); err != nil {
				klog.Errorf("Failed to write message: %v", err)
				conn.Close()
				return
			}
		}
	}
}

// pushOfflineMessages 推送离线消息
func (h *Handler) pushOfflineMessages(ctx context.Context, userID string) {
	req := &messagepb.GetUnreadMessagesRequest{
		Id: userID,
	}

	resp, err := h.grpcClient.MessageService().GetUnreadMessages(ctx, req)
	if err != nil {
		klog.Errorf("Failed to get unread messages: %v", err)
		return
	}

	if len(resp.Messages) == 0 {
		return
	}

	// 转换消息格式
	messages := make([]map[string]interface{}, len(resp.Messages))
	for i, m := range resp.Messages {
		messages[i] = map[string]interface{}{
			"id":           m.Id,
			"type":         m.Type,
			"from_id":      m.FromId,
			"to_id":        m.ToId,
			"group_id":     m.GroupId,
			"content":      m.Content,
			"message_type": m.MessageType,
			"status":       m.Status,
			"created_at":   m.CreatedAt,
			"updated_at":   m.UpdatedAt,
		}
	}

	if err := h.pusher.PushOfflineMessages(userID, messages); err != nil {
		klog.Errorf("Failed to push offline messages: %v", err)
	}
}

// extractToken 从请求中提取 Token
func extractToken(c *gin.Context) string {
	// 从 Header 获取
	token := c.GetHeader("Authorization")
	if token != "" && strings.HasPrefix(token, "Bearer ") {
		return token[7:]
	}

	// 从 Query 参数获取（备用方案）
	token = c.Query("token")
	return token
}
