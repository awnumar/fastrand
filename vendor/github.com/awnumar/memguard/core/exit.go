package core

import (
	"os"

	"github.com/awnumar/memguard/memcall"
)

/*
Purge wipes all sensitive data and keys before reinitialising the session with a fresh encryption key and secure values. Subsequent library operations will use these fresh values and the old data is assumed to be practically unrecoverable.

The creation of new Enclave objects should wait for this function to return since subsequent Enclave objects will use the newly created key.

This function should be called before the program terminates, or else the provided Exit or Panic functions should be used to terminate.
*/
func Purge() {
	// Generate a new encryption key, wiping the old.
	err := key.Initialise()
	if err != nil {
		key.Destroy()
		key = NewCoffer()
	}

	// Get a snapshot of existing Buffers.
	snapshot := buffers.flush()
	buffers.add(key.left, key.right, buf32)

	// Destroy them, performing the usual sanity checks.
	for _, b := range snapshot {
		// Don't destroy the key partitions.
		if b != key.left && b != key.right && b != buf32 {
			b.Destroy()
		}
	}
}

/*
Exit terminates the process with a specified exit code but securely wipes and cleans up sensitive data before doing so.
*/
func Exit(c int) {
	// Wipe the encryption key used to encrypt data inside Enclaves.
	key.Destroy()

	// Get a snapshot of existing Buffers.
	snapshot := buffers.flush()

	// Destroy them, performing the usual sanity checks.
	for _, b := range snapshot {
		b.Destroy()
	}

	// Exit with the specified exit code.
	os.Exit(c)
}

/*
Panic is identical to the builtin panic except it wipes all it can before calling panic.
*/
func Panic(v interface{}) {
	// Wipe both halves of the Enclave encryption key.
	Wipe(key.left.Data())
	Wipe(key.right.Data())

	// Wipe all of the currently active LockedBuffers.
	for _, b := range buffers.list {
		if !b.mutable {
			memcall.Protect(b.inner, memcall.ReadWrite)
		}
		Wipe(b.Data())
	}

	// Panic.
	panic(v)
}
