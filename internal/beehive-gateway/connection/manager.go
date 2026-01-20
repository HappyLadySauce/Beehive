package connection

import (
	"sync"
)

// Manager 连接管理器
type Manager struct {
	connections map[string]*Connection // userID -> Connection
	mu          sync.RWMutex
}

// NewManager 创建新的连接管理器
func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
	}
}

// Register 注册连接
func (m *Manager) Register(userID string, conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果用户已有连接，先关闭旧连接
	if oldConn, exists := m.connections[userID]; exists {
		oldConn.Close()
	}

	m.connections[userID] = conn
}

// Unregister 注销连接
func (m *Manager) Unregister(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.connections[userID]; exists {
		conn.Close()
		delete(m.connections, userID)
	}
}

// Get 获取连接
func (m *Manager) Get(userID string) (*Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[userID]
	return conn, exists
}

// GetAll 获取所有连接
func (m *Manager) GetAll() map[string]*Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*Connection)
	for userID, conn := range m.connections {
		result[userID] = conn
	}
	return result
}

// Count 返回连接数量
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// IsOnline 检查用户是否在线
func (m *Manager) IsOnline(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[userID]
	if !exists {
		return false
	}
	return !conn.IsClosed()
}
