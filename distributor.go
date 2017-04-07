package msgrelay

import "sync"

type Distributor struct {
	sync.Mutex
	users map[UserID]*User
}

func NewDistributor() *Distributor {
	d := Distributor{
		users: make(map[UserID]*User),
	}
	return &d
}
