package msgrelay

import "sync"

type SessionID string

type Session struct {
	sync.RWMutex
	id      SessionID
	online  bool
	channel chan *Message
}

func NewSession(id SessionID) *Session {
	return &Session{
		id:      id,
		channel: make(chan *Message),
	}
}

func (s *Session) SetOnline() {
	s.RLock()
	defer s.RUnlock()
	s.online = true
}

func (s *Session) SetOffline() {
	s.RLock()
	defer s.RUnlock()
	s.online = false
}

func (s *Session) SendMsg(msg *Message) {
	s.channel <- msg
}

func (s *Session) IsOnline() bool {
	s.Lock()
	defer s.Unlock()
	return s.online
}
