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

type ByteArray struct {
	byteArray
}

func (ba *ByteArray) FromBytes(bs []byte) error {
	if ba.capacity < C.uint(len(bs)) {
		return fmt.Errorf("bytearray is too small, required %v bytes but just %v is available", len(bs), ba.capacity)
	}
	C.memcpy(unsafe.Pointer(ba.bytes), unsafe.Pointer(&bs[0]), C.ulong(len(bs)))
	ba.length = C.uint(len(bs))
	return nil
}

func (ba *ByteArray) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(ba.bytes), C.int(ba.length))
}

func (ba *ByteArray) Destroy() {
	if ba.byteArray.capacity != 0 {
		C.svm_byte_array_destroy(ba.byteArray)
		ba.byteArray.capacity = 0
		ba.byteArray.length = 0
	}
}

func (ba *ByteArray) Close() error {
	ba.Destroy()
	return nil
}
