package svm

/*
#cgo LDFLAGS: -lsvm
#include "svm.h"
#include "memory.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type byteArray = C.svm_byte_array

// ByteArray is the svm_byte_array wrapper
type ByteArray struct {
	byteArray
}

// FromBytes fills ByteArry with bytes
func (ba *ByteArray) FromBytes(bs []byte) error {
	if ba.capacity < C.uint(len(bs)) {
		return fmt.Errorf("bytearray is too small, required %v bytes but just %v is available", len(bs), ba.capacity)
	}
	C.memcpy(unsafe.Pointer(ba.bytes), unsafe.Pointer(&bs[0]), C.ulong(len(bs)))
	ba.length = C.uint(len(bs))
	return nil
}

// Bytes returns bytes from wrapped svm_byte_array
func (ba *ByteArray) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(ba.bytes), C.int(ba.length))
}

// Destroy releases associated resources with svm_byte_array
func (ba *ByteArray) Destroy() {
	if ba.byteArray.capacity != 0 {
		C.svm_byte_array_destroy(ba.byteArray)
		ba.byteArray.capacity = 0
		ba.byteArray.length = 0
	}
}

// Close is the go idiomatic way to destory svm_byte_array
func (ba *ByteArray) Close() error {
	ba.Destroy()
	return nil
}
