package fastrand

import (
	"bytes"
	"testing"
	"unsafe"
)

func TestNewBuffer(t *testing.T) {
	b := newbuffer()
	if len(b.Data) != pageSize || cap(b.Data) != pageSize {
		t.Error("incorrect data region size")
	}
	if len(b.memory) != 3*pageSize || cap(b.memory) != 3*pageSize {
		t.Error("incorrect memory region size")
	}
	if len(b.preguard) != pageSize || cap(b.preguard) != pageSize {
		t.Error("incorrect preguard size")
	}
	if len(b.postguard) != pageSize || cap(b.postguard) != pageSize {
		t.Error("incorrect postguard size")
	}
	for i := range b.Data {
		if b.Data[i] != 0 {
			t.Error("memory should be zero")
		}
	}
	b.kill()
}

func TestKill(t *testing.T) {
	b := newbuffer()
	b.kill()
	if b.Data != nil {
		t.Error("data should be nil")
	}
	if b.memory != nil {
		t.Error("memory should be nil")
	}
	if b.preguard != nil || b.postguard != nil {
		t.Error("guard pages should be nil")
	}
}

func TestWipe(t *testing.T) {
	b := make([]byte, 4096)
	for i := range b {
		b[i] = 1
	}
	wipe(b)
	for i := range b {
		if b[i] != 0 {
			t.Error("wipe unsuccessful")
		}
	}
}

func TestGetBytes(t *testing.T) {
	buffer := make([]byte, 32)
	derived := getBytes(&buffer[0], len(buffer))
	if !bytes.Equal(buffer, derived) {
		t.Error("naive equality check failed")
	}
	buffer[0] = 1
	buffer[31] = 1
	if !bytes.Equal(buffer, derived) {
		t.Error("modified equality check failed")
	}
	if uintptr(unsafe.Pointer(&buffer[0])) != uintptr(unsafe.Pointer(&derived[0])) {
		t.Error("pointer values differ")
	}
	if len(buffer) != len(derived) || cap(buffer) != cap(derived) {
		t.Error("length or capacity values differ")
	}
}
