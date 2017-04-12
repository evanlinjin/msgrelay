package msgrelay

import "sync"

type Distributor struct {
	users map[UserID]*User
	addU  chan *User
	remU  chan UserID
	getU  chan *UserList
	msgs  chan *Message
	quit  chan bool
}

func NewDistributor() *Distributor {
	d := Distributor{
		users: make(map[UserID]*User),
		addU:  make(chan *User),
		remU:  make(chan UserID),
		getU:  make(chan *UserList),
		msgs:  make(chan *Message),
		quit:  make(chan bool),
	}
	d.runService()
	return &d
}

func (d *Distributor) AddUser(uid UserID) *User {
	u := NewUser(uid)
	d.addU <- u
	return u
}

func (d *Distributor) RemoveUser(uid UserID) {
	d.remU <- uid
}

func (d *Distributor) GetUsers(ids ...UserID) []*User {
	users := make([]*User, len(ids))
	l := NewUserList(ids...)
	d.getU <- l
	for i := 0; i < len(ids); i++ {
		users[i] = <-l.c
	}
	return users
}

func (d *Distributor) runService() {
	for {
		select {
		case u := <-d.addU:
			d.users[u.id] = u

		case i := <-d.remU:
			delete(d.users, i)

		case l := <-d.getU:
			for _, id := range l.ids {
				u, _ := d.users[id]
				l.c <- u
			}

		case <-d.quit:
			return
		}
	}
}

func (d *Distributor) endService() {
	d.quit <- true
}
