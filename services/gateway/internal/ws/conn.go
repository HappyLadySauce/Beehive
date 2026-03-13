package ws

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Connection 表示单条 WebSocket 连接，含 ConnID、登录后绑定的 UserID 等。
type Connection struct {
	ConnID    string
	UserID    string // 登录后由 auth 逻辑绑定
	DeviceID  string
	GatewayID string
	conn      *websocket.Conn
	writeMu   sync.Mutex
}

// Hub 管理当前进程内所有 WebSocket 连接，用于分配 ConnID 与按需推送。
type Hub struct {
	gatewayID string
	mu        sync.RWMutex
	conns     map[string]*Connection
}

// NewHub 创建连接管理器，gatewayID 用于多实例部署时区分本实例。
func NewHub(gatewayID string) *Hub {
	if gatewayID == "" {
		gatewayID = "gateway-1"
	}
	return &Hub{
		gatewayID: gatewayID,
		conns:     make(map[string]*Connection),
	}
}

// Register 将已升级的 WebSocket 注册到 Hub，返回带 ConnID 的 Connection。
func (h *Hub) Register(conn *websocket.Conn, deviceID string) *Connection {
	connID := uuid.Must(uuid.NewUUID()).String()
	c := &Connection{
		ConnID:    connID,
		GatewayID: h.gatewayID,
		DeviceID:  deviceID,
		conn:      conn,
	}
	h.mu.Lock()
	h.conns[connID] = c
	h.mu.Unlock()
	return c
}

// Unregister 从 Hub 移除连接（关闭时调用）。
func (h *Hub) Unregister(connID string) {
	h.mu.Lock()
	delete(h.conns, connID)
	h.mu.Unlock()
}

// Get 按 ConnID 获取连接。
func (h *Hub) Get(connID string) *Connection {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.conns[connID]
}

// BindUser 在登录成功后绑定 UserID。
func (c *Connection) BindUser(userID string) {
	c.UserID = userID
}

// WriteJSON 向该连接写入 JSON（加锁，供推送与响应使用）。
func (c *Connection) WriteJSON(v interface{}) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteJSON(v)
}

// Conn 返回底层 WebSocket 连接，供读循环使用。
func (c *Connection) Conn() *websocket.Conn {
	return c.conn
}

// Close 关闭底层 WebSocket 连接。
func (c *Connection) Close() error {
	return c.conn.Close()
}
