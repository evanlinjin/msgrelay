package msgrelay

import "sync"

type UserID string

type User struct {
	sync.RWMutex
	id       UserID
	sessions []*Session
	addS     chan SessionID
	remS     chan SessionID
	quit     chan bool
}

func NewUser(id UserID) *User {
	u := User{
		id:   id,
		addS: make(chan SessionID),
		remS: make(chan SessionID),
		quit: make(chan bool),
	}
	u.runService()
	return &u
}

func (u *User) AddSession(id SessionID) {
	u.addS <- id
}

func (u *User) RemoveSession(id SessionID) {
	u.remS <- id
}

func (u *User) runService() {
	go func() {
		for {
		UserRunServiceSelect:
			select {
			case sid := <-u.addS:
				for _, s := range u.sessions {
					if s.id == sid {
						goto UserRunServiceSelect
					}
				}
				u.RLock()
				u.sessions = append(u.sessions, NewSession(sid))
				u.RUnlock()

			case sid := <-u.remS:
				for i, s := range u.sessions {
					if s.id == sid {
						u.RLock()
						// Removes element 'i'.
						u.sessions[i] = u.sessions[len(u.sessions)-1]
						u.sessions = u.sessions[:len(u.sessions)-1]
						break
						u.RUnlock()
					}
				}

			case <-u.quit:
				return
			}
		}
	}()
}

func (u *User) endService() {
	go func() {
		u.quit <- true
	}()
}
