package msgrelay

import "sync"

type UserID string

type User struct {
	sync.RWMutex
	id       UserID
	sessions []*Session
	online   bool
	addS     chan *Session
	remS     chan SessionID
	msgs     chan *Message
	quit     chan bool
}

func NewUser(id UserID) *User {
	u := User{
		id:     id,
		online: false,
		addS:   make(chan *Session),
		remS:   make(chan SessionID),
		msgs:   make(chan *Message),
		quit:   make(chan bool),
	}
	go u.runService()
	return &u
}

func (u *User) AddSession(id SessionID) *Session {
	s := NewSession(id)
	u.addS <- s
	return s
}

func (u *User) RemoveSession(id SessionID) {
	u.remS <- id
}

func (u *User) ReceiveMsg(msg *Message) {
	u.msgs <- msg
}

func (u *User) NewMsg(to UserID, data interface{}) *Message {
	return NewMessage(u.id, to, data)
}

func (u *User) IsOnline() bool {
	u.Lock()
	defer u.Unlock()
	return u.online
}

func (u *User) runService() {
	for {
		select {
		case s := <-u.addS:
			u.RLock()
			u.sessions = append(u.sessions, s)
			u.online = true
			u.RUnlock()

		case sid := <-u.remS:
			for i, s := range u.sessions {
				if s.id == sid {
					u.RLock()
					// Removes element 'i'.
					u.sessions[i] = u.sessions[len(u.sessions)-1]
					u.sessions = u.sessions[:len(u.sessions)-1]
					u.online = len(u.sessions) > 0
					u.RUnlock()
					break
				}
			}

		case msg := <-u.msgs:
			u.Lock()
			for _, s := range u.sessions {
				s.ReceiveMessage(msg)
			}
			u.Unlock()

		case <-u.quit:
			return
		}
	}
}

func (u *User) endService() {
	u.quit <- true
}
