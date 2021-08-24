package svm

import "C"

// #include "./svm.h"
//
import "C"
import (
	"fmt"
	"unsafe"
)

type (
	cSvmByteArray = C.svm_byte_array
)

func bytesCloneToSvmByteArray(b []byte) cSvmByteArray {
	var ba cSvmByteArray
	ba.FromBytesClone(b)
	return ba
}

func bytesAliasToSvmByteArray(b []byte) cSvmByteArray {
	var ba cSvmByteArray
	ba.FromBytesAlias(b)
	return ba
}

func svmByteArrayCloneToBytes(ba cSvmByteArray) []byte {
	return ba.AsCBytes().GoBytesClone()
}

func (ba *cSvmByteArray) FromBytesClone(b []byte) {
	if len(b) == 0 {
		return
	}

	cBytes := GoBytes(b).CBytesClone()
	ba.bytes = (*cUchar)(cBytes.data)
	ba.length = (cUint)(cBytes.len)
}

func (ba *cSvmByteArray) FromBytesAlias(b []byte) {
	if len(b) == 0 {
		return
	}

	cBytes := GoBytes(b).CBytesAlias()
	ba.bytes = (*cUchar)(cBytes.data)
	ba.length = (cUint)(cBytes.len)
}

// AsCBytes converts an cSvmByteArray struct to a C array.
func (ba cSvmByteArray) AsCBytes() CBytes {
	return CBytes{
		data: unsafe.Pointer(ba.bytes),
		len:  int(ba.length),
	}
}

// svmError converts an SVM byte array to a Go error.
func (ba cSvmByteArray) svmError() error {
	b := svmByteArrayCloneToBytes(ba)
	return newSvmError(b)
}

// String helps cSvmByteArray to implement the Stringer interface.
func (ba cSvmByteArray) String() string {
	return fmt.Sprintf("%s", ba.AsCBytes().GoBytesAlias())
}

// Free deallocate the cSvmByteArray struct from via C allocator.
func (ba cSvmByteArray) Free() {
	ba.AsCBytes().Free()
}

// SvmDestroy deallocate the cSvmByteArray struct via the SVM (Rust) allocator.
func (ba cSvmByteArray) SvmFree() {
	cSvmByteArrayDestroy(ba)
}
