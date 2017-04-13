package msgrelay

import (
	"sync"
	"time"
)

type Message struct {
	sync.RWMutex
	from   UserID
	toUser UserID
	//toGroup
	data interface{}
	sent time.Time
}

func NewMessage(from, to UserID, data interface{}) *Message {
	return &Message{
		from:   from,
		toUser: to,
		data:   data,
	}
}
