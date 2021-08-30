package svm

/*
#include "svm.h"
#include "memory.h"
*/
import "C"
import (
	"unsafe"
)

// State is a SVM state abstraction
type State [32]byte

// Context is a SVM context wrapper
type Context struct {
	ByteArray
}

// NewContext returns new clean context
func NewContext() *Context {
	c := &Context{}
	c.byteArray = C.svm_context_alloc()
	return c
}

// Context creates new context from the State
func (st State) Context() *Context {
	c := NewContext()
	C.memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(c.bytes))+40), unsafe.Pointer(&st[0]), C.ulong(32))
	return c
}
