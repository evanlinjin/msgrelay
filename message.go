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
