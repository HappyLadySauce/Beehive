package libnet

import (
	"sync"

	sessionId "github.com/HappyLadySauce/Beehive/common/utils/sessionid"
	"github.com/HappyLadySauce/Beehive/common/utils/hash"
)

// 会话map数量
// 用于将会话根据token散列到不同的会话map中
// 减少锁竞争，提高性能
const sessionMapNum = 32

// Manager 管理器
type Manager struct {
	Name        string
	sessionMaps [sessionMapNum]sessionMap
	disposeFlag bool
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

// sessionMap 会话map
type sessionMap struct {
	sync.RWMutex
	sessions      map[sessionId.SessionId]*Session
	tokenSessions map[string][]sessionId.SessionId
}

// 创建管理器
func NewManager(name string) *Manager {
	manager := &Manager{
		Name: name,
	}
	for i := 0; i < sessionMapNum; i++ {
		manager.sessionMaps[i].sessions = make(map[sessionId.SessionId]*Session)
		manager.sessionMaps[i].tokenSessions = make(map[string][]sessionId.SessionId)
	}
	return manager
}

// 获取会话
func (m *Manager) GetSession(sessionId sessionId.SessionId) *Session {
	token := sessionId.Token()
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.RLock()
	defer smap.RUnlock()

	session, ok := smap.sessions[sessionId]
	if !ok {
		return nil
	}
	return session
}

// 添加会话
func (m *Manager) AddSession(session *Session) {
	sessionId := session.SessionId()
	token := session.token
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	smap.sessions[sessionId] = session
	smap.tokenSessions[token] = append(smap.tokenSessions[token], sessionId)
}

// 删除会话
func (m *Manager) removeSession(session *Session) {
	sessionId := session.SessionId()
	token := session.token
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	smap.sessions[sessionId] = nil
	smap.tokenSessions[token] = smap.tokenSessions[token][:len(smap.tokenSessions[token])-1]
	if len(smap.tokenSessions[token]) == 0 {
		delete(smap.tokenSessions, token)
	}
	delete(smap.sessions, sessionId)
}

// 关闭管理器
func (m *Manager) Close() {
	m.disposeOnce.Do(func() {
		m.disposeFlag = true

		for i := 0; i < sessionMapNum; i++ {
			smap := &m.sessionMaps[i]
			smap.Lock()
			for _, session := range smap.sessions {
				session.Close()
			}
			smap.Unlock()
		}
		m.disposeWait.Wait()
	})
}