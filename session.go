package msgrelay

import "sync"

const (
	Online  = true
	Offline = false
)

type SessionID string

type Session struct {
	sync.RWMutex
	id         SessionID
	status     bool
	msgChan    chan *Message
	statusChan chan bool
	quitChan   chan struct{}
	C          chan *Message
}

func NewSession(id SessionID) *Session {
	return &Session{
		id:         id,
		status:     false,
		msgChan:    make(chan *Message),
		statusChan: make(chan bool),
		quitChan:   make(chan struct{}),
		C:          make(chan *Message),
	}
}

func (s *Session) SetOnline() chan *Message {
	s.RLock()
	defer s.RUnlock()
	s.status = true
	return s.msgChan
}

func (s *Session) SetOffline() {
	s.RLock()
	defer s.RUnlock()
	s.status = false
}

func (s *Session) ReceiveMessage(msg *Message) {
	s.msgChan <- msg
}

func (s *Session) IsOnline() bool {
	s.Lock()
	defer s.Unlock()
	return s.status
}

func (s *Session) runService() {
	for {
		select {
		case status := <-s.statusChan:
			s.RLock()
			if s.status == Offline && status == Online {
				// TODO: When coming online.
			}
			s.status = status
			s.RUnlock()

		case msg := <-s.msgChan:
			if s.status == Online {
				go func() { s.C <- msg }()
			} else {
				// TODO: Offline storage.
			}

		case <-s.quitChan:
			return
		}
	}
}
