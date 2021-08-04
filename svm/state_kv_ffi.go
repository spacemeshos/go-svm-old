package svm

// #include <stdint.h>
// #include "./svm.h"
//
// extern void kv_get(uint8_t*, uint32_t, uint8_t*, uint32_t*);
// extern void kv_set(uint8_t*, uint32_t, uint8_t*, uint32_t);
// extern void kv_discard();
// extern void kv_checkpoint(uint8_t*);
// extern void kv_head(uint8_t*);
//
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

//export kv_get
func kv_get(keyPtr *C.uint8_t, keyLen C.uint32_t, valuePtr *C.uint8_t, valueLen *C.uint32_t) {
	// Create []byte slice alias to the key.
	keyAlias := NewCBytes(unsafe.Pointer(keyPtr), int(keyLen)).GoBytesAlias()

	// Invoke handler.
	f := kvHandlers.get
	if f == nil {
		panic("go-svm: `get` handler wasn't registered for FFI state KV")
	}
	result := f(keyAlias)

	if result == nil {
		*valueLen = 0
		return
	}

	resultLen := len(result)
	if resultLen != KVValueSize {
		panic(fmt.Sprintf("go-svm: `get` returned an invalid value size; expected: %v, got: %v", KVValueSize, resultLen))
	}

	// Create []byte slice alias to the value buffer and copy the result into there.
	valueAlias := NewCBytes(unsafe.Pointer(valuePtr), int(resultLen)).GoBytesAlias()
	copy(valueAlias, result)
	runtime.KeepAlive(result)

	// Update the result value len.
	*valueLen = C.uint(resultLen)
}

//export kv_set
func kv_set(keyPtr *C.uint8_t, keyLen C.uint32_t, valuePtr *C.uint8_t, valueLen C.uint32_t) {
	// Create []byte slice aliases.
	keyAlias := NewCBytes(unsafe.Pointer(keyPtr), int(keyLen)).GoBytesAlias()
	valueAlias := NewCBytes(unsafe.Pointer(valuePtr), int(valueLen)).GoBytesAlias()

	// Invoke handler.
	f := kvHandlers.set
	if f == nil {
		panic("go-svm: `set` handler wasn't registered for FFI state KV")
	}
	f(keyAlias, valueAlias)
}

//export kv_discard
func kv_discard() {
	f := kvHandlers.discard
	if f == nil {
		panic("go-svm: `discard` handler wasn't registered for FFI state KV")
	}
	f()
}

//export kv_checkpoint
func kv_checkpoint(statePtr *C.uint8_t) {
	f := kvHandlers.checkpoint
	if f == nil {
		panic("go-svm: `checkpoint` handler wasn't registered for FFI state KV")
	}
	result := f()
	resultLen := len(result)
	if resultLen != StateSize {
		panic(fmt.Sprintf("go-svm: `checkpoint` returned an invalid state size; expected: %v, got: %v", StateSize, resultLen))
	}

	// Create []byte slice alias to the state buffer and copy the result into there.
	stateAlias := NewCBytes(unsafe.Pointer(statePtr), resultLen).GoBytesAlias()
	copy(stateAlias, result)
	runtime.KeepAlive(result)
}

//export kv_head
func kv_head(headPtr *C.uint8_t) {
	f := kvHandlers.head
	if f == nil {
		panic("go-svm: `head` handler wasn't registered for FFI state KV")
	}
	result := f()
	resultLen := len(result)
	if resultLen != StateSize {
		panic(fmt.Sprintf("go-svm: `head` returned an invalid state size; expected: %v, got: %v", StateSize, resultLen))
	}

	// Create []byte slice alias to the state buffer and copy the result into there.
	headAlias := NewCBytes(unsafe.Pointer(headPtr), resultLen).GoBytesAlias()
	copy(headAlias, result)
	runtime.KeepAlive(result)
}

func cSvmFFIStateKVCreate(p *unsafe.Pointer) cSvmResultT {
	return (cSvmResultT)(C.svm_ffi_state_kv_create(
		p,
		(*[0]byte)(C.kv_get),
		(*[0]byte)(C.kv_set),
		(*[0]byte)(C.kv_discard),
		(*[0]byte)(C.kv_checkpoint),
		(*[0]byte)(C.kv_head),
	))
}

// kvHandlers is a static container for the KV-ops handlers, written in Go,
// to be invoked from the unsafe, cgo-exported handlers.
var kvHandlers = struct {
	get        func([]byte) []byte
	set        func([]byte, []byte)
	discard    func()
	checkpoint func() []byte
	head       func() []byte
}{}

type StateKV_FFI struct {
	// _inner is a pointer to an SVM-managed heap allocation.
	_inner unsafe.Pointer
}

func NewStateKV_FFI() (StateKV_FFI, error) {
	var p unsafe.Pointer
	if res := cSvmFFIStateKVCreate(&p); res != cSvmSuccess {
		return StateKV_FFI{}, fmt.Errorf("failed to create FFI state KV store")
	}

	return StateKV_FFI{p}, nil
}

func (StateKV_FFI) RegisterGet(f func([]byte) []byte) {
	kvHandlers.get = f
}

func (StateKV_FFI) RegisterSet(f func([]byte, []byte)) {
	kvHandlers.set = f
}

func (StateKV_FFI) RegisterDiscard(f func()) {
	kvHandlers.discard = f
}

func (StateKV_FFI) RegisterCheckpoint(f func() []byte) {
	kvHandlers.checkpoint = f
}

func (StateKV_FFI) RegisterHead(f func() []byte) {
	kvHandlers.head = f
}

func (kv StateKV_FFI) Free() {
	cSvmStateKVDestroy(kv._inner)
}
