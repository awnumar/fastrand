package fastrand

import "testing"

var tp = newpool()

func TestNewPool(t *testing.T) {
	if cap(tp.queue) != 4096 {
		t.Error("incorrect buffer size")
	}
	if len(tp.queue) != pageSize/64 {
		t.Error("incorrect number of buffers added", len(tp.queue))
	}
}

func TestAdd(t *testing.T) {
	tp.add()
	if len(tp.queue) != pageSize/32 {
		t.Error("buffers not added", len(tp.queue))
	}
}

func TestCycle(t *testing.T) {
	l := len(tp.queue)
	tp.put(make([]byte, 64))
	if len(tp.queue) != l+1 {
		t.Error("buffer not added")
	}
	for i := 0; i < 8192; i++ {
		b := tp.get()
		if len(b) != 64 {
			t.Error("incorrect buffer size")
		}
		for j := range b {
			if b[j] != 0 {
				t.Error("buffer not zeroed")
			}
			b[j] = 1
		}
		tp.put(b)
	}
}

func TestEmpty(t *testing.T) {
	for i := 0; i < 256; i++ {
		b := tp.get()
		if len(b) != 64 {
			t.Error("incorrect buffer size")
		}
	}
}
