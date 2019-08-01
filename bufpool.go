package fastrand

import (
	"sync"
)

// pool of buffers within securely allocated memory
var bp = newpool()

type pool struct {
	queue chan []byte

	sync.RWMutex
	allocs []*buffer // todo
}

func newpool() (p *pool) {
	p = &pool{queue: make(chan []byte, 4096)}
	p.add()
	return
}

func (p *pool) add() {
	b := newbuffer()
	for i := pageSize - 64; i >= 0; i -= 64 {
		p.put(getBytes(&b.Data[i], 64)) // construct slice to hide capacity
	}
}

func (p *pool) get() (b []byte) {
	select {
	case b = <-p.queue:
		// got buf from pool
	default:
		// pool empty, add more chunks
		p.add()
		b = p.get()
	}
	return
}

func (p *pool) put(buf []byte) {
	wipe(buf)
	select {
	case p.queue <- buf:
		// buffer went into pool
	default:
		// pool was full; drop value (todo: handle this? grow pool?)
	}
}
