package fastrand

import (
	"os"
	"runtime"
	"sync"
	"unsafe"

	"github.com/awnumar/memcall"
)

var pageSize = os.Getpagesize()

// fixed size allocation with data region = page size
type buffer struct {
	sync.RWMutex

	Data   []byte
	memory []byte

	preguard  []byte
	postguard []byte
}

func newbuffer() (b *buffer) {
	var err error
	b = new(buffer)

	b.memory, err = memcall.Alloc(3 * pageSize)
	if err != nil {
		panic(err)
	}
	b.Data = getBytes(&b.memory[pageSize], pageSize)
	b.preguard = getBytes(&b.memory[0], pageSize)
	b.postguard = getBytes(&b.memory[2*pageSize], pageSize)

	if err := memcall.Lock(b.Data); err != nil {
		panic(err)
	}
	if err := memcall.Protect(b.preguard, memcall.NoAccess); err != nil {
		panic(err)
	}
	if err := memcall.Protect(b.postguard, memcall.NoAccess); err != nil {
		panic(err)
	}

	return
}

func (b *buffer) kill() {
	wipe(b.Data)

	if err := memcall.Unlock(b.Data); err != nil {
		panic(err)
	}
	if err := memcall.Free(b.memory); err != nil {
		panic(err)
	}

	b.Data = nil
	b.memory = nil
	b.preguard = nil
	b.postguard = nil
}

func wipe(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
	runtime.KeepAlive(buf)
}

func getBytes(ptr *byte, len int) []byte {
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(ptr)), len, len}
	return *(*[]byte)(unsafe.Pointer(&sl))
}
