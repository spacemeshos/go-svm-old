package svm

/*
#include "svm.h"
#include "memory.h"
*/
import "C"
import (
	"unsafe"
)

type State [32]byte

type Context struct {
	ByteArray
}

func NewContext() *Context {
	c := &Context{}
	c.byteArray = C.svm_context_alloc()
	return c
}

func (st State) Context() *Context {
	c := NewContext()
	C.memcpy(unsafe.Pointer(uintptr(unsafe.Pointer(c.bytes))+40), unsafe.Pointer(&st[0]), C.ulong(32))
	return c
}
