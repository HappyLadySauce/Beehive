package libnet

import (
	"sync"

	sessionId "github.com/HappyLadySauce/Beehive/common/utils/sessionid"
	"github.com/HappyLadySauce/Beehive/common/utils/hash"
)

const sessionMapNum = 32

type Manager struct {
	Name        string
	sessionMaps [sessionMapNum]sessionMap
	disposeFlag bool
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

type sessionMap struct {
	sync.RWMutex
	sessions      map[sessionId.SessionId]*Session
	tokenSessions map[string][]sessionId.SessionId
}

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