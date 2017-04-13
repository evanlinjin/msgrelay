package msgrelay

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestDistributor_AddUser(t *testing.T) {
	d := NewDistributor()
	u := d.AddUser("test_uid")

	sessions := make([]*Session, 3)
	channels := make([]chan *Message, 3)

	for i := 0; i < len(sessions); i++ {
		sessions[i] = u.AddSession("session" + strconv.Itoa(i))
		channels[i] = sessions[i].SetOnline()
	}

	quitChan := make(chan struct{})

	go func() {
		for {
			select {
			case m := <-channels[0]:
				fmt.Println("[s1]", m.data)
			case m := <-channels[1]:
				fmt.Println("[s2]", m.data)
			case m := <-channels[2]:
				fmt.Println("[s3]", m.data)
			case <-quitChan:
				return
			}
		}
	}()

	m := u.NewMsg("test_uid", "Hello, this is a test")
	u.ReceiveMsg(m)

	timer := time.NewTimer(time.Second)
	<-timer.C

	u.endService()
	d.endService()

	timer = time.NewTimer(time.Second)
	<-timer.C
}
