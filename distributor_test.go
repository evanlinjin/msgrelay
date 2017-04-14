package msgrelay

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func makeNewUserWithSessions(
	distributor *Distributor,
	uid UserID,
	nSessions int,
) (
	user *User,
	quitFunc func(),
) {
	user = distributor.AddUser(uid)
	sessions := make([]*Session, nSessions)
	quitChan := make(chan struct{})

	for i := 0; i < nSessions; i++ {
		sessions[i] = user.AddSession(
			SessionID(string(uid) + "_session" + strconv.Itoa(i)),
		)
		sessions[i].SetStatus(Online)

		go func(sid SessionID, c chan *Message, q chan struct{}){
			for {
				select {
				case msg := <-c:
					fmt.Printf("[%s] Message: %v\n", sid, msg.data)
				case <-q:
					fmt.Printf("[Listener] %s : listener ended.\n", sid)
					return
				}
			}
		}(sessions[i].id, sessions[i].C, quitChan)
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
	_, q := makeNewUserWithSessions(d, "user0", 10)
	d.EndAll()
	q()
	<-time.NewTimer(time.Second).C
}

func TestDistributor_AddUser2(t *testing.T) {
	const nUsers = 10
	d := NewDistributor()
	users := make([]*User, nUsers)
	quits := make([]func(), nUsers)
	for i := 0; i < nUsers; i++ {
		users[i], quits[i] = makeNewUserWithSessions(
			d, UserID("user" + strconv.Itoa(i)), 1)
	}
	d.EndAll()
	for _, q := range quits {
		q()
	}
	<-time.NewTimer(2*time.Second).C
}