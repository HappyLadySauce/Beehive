package connection

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection WebSocket 连接封装
type Connection struct {
	conn        *websocket.Conn
	userID      string
	connectedAt time.Time
	lastPing    time.Time
	sendChan    chan []byte
	closeChan   chan struct{}
	closed      bool
	mu          sync.RWMutex
}

// NewConnection 创建新的连接
func NewConnection(conn *websocket.Conn, userID string) *Connection {
	return &Connection{
		conn:        conn,
		userID:      userID,
		connectedAt: time.Now(),
		lastPing:    time.Now(),
		sendChan:    make(chan []byte, 256),
		closeChan:   make(chan struct{}),
		closed:      false,
	}
}

// UserID 返回用户ID
func (c *Connection) UserID() string {
	return c.userID
}

// Conn 返回 WebSocket 连接
func (c *Connection) Conn() *websocket.Conn {
	return c.conn
}

// Send 发送消息
func (c *Connection) Send(data []byte) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return websocket.ErrCloseSent
	}
	c.mu.RUnlock()

	select {
	case c.sendChan <- data:
		return nil
	case <-c.closeChan:
		return websocket.ErrCloseSent
	}
}

// Close 关闭连接
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	close(c.closeChan)
	close(c.sendChan)
	return c.conn.Close()
}

// IsClosed 检查连接是否已关闭
func (c *Connection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// UpdateLastPing 更新最后心跳时间
func (c *Connection) UpdateLastPing() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastPing = time.Now()
}

// LastPing 返回最后心跳时间
func (c *Connection) LastPing() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastPing
}

// ConnectedAt 返回连接建立时间
func (c *Connection) ConnectedAt() time.Time {
	return c.connectedAt
}

// SendChan 返回发送通道
func (c *Connection) SendChan() <-chan []byte {
	return c.sendChan
}
