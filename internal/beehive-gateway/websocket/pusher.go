package websocket

import (
	"encoding/json"
	"time"

	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/connection"
)

// Pusher 消息推送器
type Pusher struct {
	connMgr *connection.Manager
}

// NewPusher 创建新的推送器
func NewPusher(connMgr *connection.Manager) *Pusher {
	return &Pusher{
		connMgr: connMgr,
	}
}

// Message 消息结构
type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// PushToUser 向指定用户推送消息
func (p *Pusher) PushToUser(userID string, msg Message) error {
	conn, exists := p.connMgr.Get(userID)
	if !exists || conn.IsClosed() {
		return nil // 用户不在线，不返回错误
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := conn.Send(data); err != nil {
		klog.Errorf("Failed to send message to user %s: %v", userID, err)
		return err
	}

	return nil
}

// PushToUsers 向多个用户推送消息
func (p *Pusher) PushToUsers(userIDs []string, msg Message) {
	for _, userID := range userIDs {
		if err := p.PushToUser(userID, msg); err != nil {
			klog.Errorf("Failed to push message to user %s: %v", userID, err)
		}
	}
}

// PushOfflineMessages 推送离线消息
func (p *Pusher) PushOfflineMessages(userID string, messages []map[string]interface{}) error {
	conn, exists := p.connMgr.Get(userID)
	if !exists || conn.IsClosed() {
		return nil
	}

	msg := Message{
		Type: "offline_messages",
		Data: map[string]interface{}{
			"messages": messages,
			"count":    len(messages),
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return conn.Send(data)
}

// SendError 发送错误消息
func (p *Pusher) SendError(userID string, errorType string, message string) error {
	msg := Message{
		Type: "error",
		Data: map[string]interface{}{
			"error_type": errorType,
			"message":    message,
			"timestamp":  time.Now().Unix(),
		},
	}
	return p.PushToUser(userID, msg)
}

// SendSuccess 发送成功消息
func (p *Pusher) SendSuccess(userID string, msgType string, data map[string]interface{}) error {
	msg := Message{
		Type: msgType,
		Data: data,
	}
	return p.PushToUser(userID, msg)
}
