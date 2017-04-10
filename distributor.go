package msgrelay

import "sync"

type Distributor struct {
	sync.Mutex
	users map[UserID]*User
	addU  chan *User
	remU  chan UserID
	msgs  chan *Message
	quit  chan bool
}

func NewDistributor() *Distributor {
	d := Distributor{
		users: make(map[UserID]*User),
		addU:  make(chan *User),
		remU:  make(chan UserID),
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

func (d *Distributor) GetUsers(uids ...UserID) (res chan *User) {
	res = make(chan *User)
	go func() {
		d.Lock()
		for _, uid := range uids {
			u, _ := d.users[uid]
			res <- u
		}
		d.Unlock()
	}()
	return res
}

func (d *Distributor) runService() {
	for {
		select {
		case u := <-d.addU:
			d.Lock()
			d.users[u.id] = u
			d.Unlock()

		case i := <-d.remU:
			d.Lock()
			delete(d.users, i)
			d.Unlock()

		case <-d.quit:
			return
		}
	}
}

func (d *Distributor) endService() {
	d.quit <- true
}
