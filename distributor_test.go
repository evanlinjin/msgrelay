package msgrelay

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"sync"
)

type Counter struct {
	sync.Mutex
	n int64
}

func NewCounter() *Counter {
	return &Counter{n:0}
}

func (c *Counter) Increment() {
	c.Lock()
	defer c.Unlock()
	c.n += 1
}

func (c *Counter) GetCount() int64 {
	c.Lock()
	defer c.Unlock()
	return c.n
}

func makeNewUserWithSessions(
	distributor *Distributor,
	uid UserID,
	nSessions int,
	msgCounter *Counter,
) (
	user *User,
	quitFunc func(),
) {
	user = distributor.AddUser(uid)
	sessions := make([]*Session, nSessions)
	quitChan := make(chan struct{})

	f := func(s *Session, q chan struct{}, c *Counter) {
		for {
			select {
			case msg := <-s.C:
				fmt.Printf("[%s] Message: %v\n", s.id, msg.data)
				c.Increment()
			case <-q:
				fmt.Printf("[Listener] %s : listener ended.\n", s.id)
				return
			}
		}
	}

	for i := 0; i < nSessions; i++ {
		sessions[i] = user.AddSession(
			SessionID(string(uid) + "_session" + strconv.Itoa(i)),
		)
		sessions[i].SetStatus(Online)
		go f(sessions[i], quitChan, msgCounter)
	}

	quitFunc = func() {
		for {
			select {
			case quitChan<- struct{}{}:
			default:
				return
			}
		}
	}
	return
}

func TestDistributor_AddUser(t *testing.T) {
	d := NewDistributor()
	c := NewCounter()
	u, q := makeNewUserWithSessions(d, "user0", 10, c)
	m := u.NewMsg(u.id, "Test message.")
	u.ReceiveMsg(m)
	<-time.NewTimer(time.Second).C
	d.EndAll()
	q()
	<-time.NewTimer(time.Second).C
	if c.GetCount() != 10 {
		t.Error("c.GetCount() != 10", "[GOT]", c.GetCount())
	}
}

func TestDistributor_AddUser2(t *testing.T) {
	const nUsers = 10
	c := NewCounter()
	d := NewDistributor()
	users := make([]*User, nUsers)
	quits := make([]func(), nUsers)
	for i := 0; i < nUsers; i++ {
		users[i], quits[i] = makeNewUserWithSessions(
			d, UserID("user" + strconv.Itoa(i)), 1, c)
	}
	d.EndAll()
	for _, q := range quits {
		q()
	}
	<-time.NewTimer(2*time.Second).C
}