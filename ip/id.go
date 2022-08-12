package ip

import "sync"

var (
	id   uint16 = 128
	idmu sync.Mutex
)

func generateID() uint16 {
	idmu.Lock()
	ret := id
	id += 1
	defer idmu.Unlock()
	return ret
}
