package svm

//#cgo CFLAGS: -I${SRCDIR}
//#cgo LDFLAGS: -lsvm
//
//#include "svm.h"
//#include "memory.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type byteArray = C.svm_byte_array

// `ByteArray` is a wrapper for `svm_byte_array`
type ByteArray struct {
	byteArray
}

func ByteArrayToSlice(raw *C.svm_byte_array) []byte {
	bytes := make([]byte, raw.length, raw.capacity)
	C.memcpy(unsafe.Pointer(&bytes[0]), unsafe.Pointer(raw.bytes), C.ulong(raw.length))

	return bytes
}

// Creates a new instance with internal buffer content copied from the given `bytes` Go byte-slice.
func (ba *ByteArray) FromBytes(bytes []byte) error {
	if ba.capacity < C.uint(len(bytes)) {
		return fmt.Errorf("`svm_byte_array` is too small, required %v bytes but just %v are available", len(bytes), ba.capacity)
	}
	C.memcpy(unsafe.Pointer(ba.bytes), unsafe.Pointer(&bytes[0]), C.ulong(len(bytes)))
	ba.length = C.uint(len(bytes))
	return nil
}

// Copies the inner buffer into a new Go allocated memory
func (ba *ByteArray) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(ba.bytes), C.int(ba.length))
}

// Releases associated resources with `svm_byte_array`
func (ba *ByteArray) Destroy() {
	if ba.byteArray.capacity != 0 {
		C.svm_byte_array_destroy(ba.byteArray)
		ba.byteArray.capacity = 0
		ba.byteArray.length = 0
	}
}

// `Close` is the Go idiomatic way to destroy `svm_byte_array`
func (ba *ByteArray) Close() error {
	ba.Destroy()
	return nil
}
