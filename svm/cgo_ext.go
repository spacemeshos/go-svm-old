package svm

import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// GoString is an alias type for string, used to define local methods.
type GoString string

// CStringClone is using the built-in CString cgo function to
// create a C string that is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be freed
// by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the underlying array,
// in oppose to creating a pointer which references it.
func (s GoString) CStringClone() CString {
	p := C.CString(string(s))
	return CString{data: unsafe.Pointer(p)}
}

// GoBytes is an alias type for []byte slice, used to define local methods.
type GoBytes []byte

// CBytesClone is using the built-in CBytes cgo function to
// create a C array that is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be freed
// by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the underlying array,
// in oppose to creating a pointer which references it.
func (b GoBytes) CBytesClone() CBytes {
	p := C.CBytes(b)
	return CBytes{data: p, len: len(b)}
}

// svmByteArrayClone converts []byte slice to SVM byte array.
// The byte array is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be freed
// by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the underlying array,
// in oppose to creating a pointer which references it.
func (b GoBytes) svmByteArrayClone() cSvmByteArray {
	cBytes := b.CBytesClone()
	return cSvmByteArray{
		bytes:  (*C.uchar)(cBytes.data),
		length: (C.uint)(cBytes.len),
	}
}

// svmAddressClone converts []byte slice to SVM address.
// The address byte array is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be freed
// by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the underlying array,
// in oppose to creating a pointer which references it.
func (b GoBytes) svmAddressClone() cSvmByteArray {
	if len(b) <= SvmAddressLen {
		return b.svmByteArrayClone()
	}

	dst := make([]byte, SvmAddressLen, SvmAddressLen)
	copy(dst, b[:SvmAddressLen])
	return GoBytes(dst).svmByteArrayClone()
}

// CString represents a C string allocated in the C heap.
type CString struct {
	data unsafe.Pointer // C pointer (allocated using malloc)
}

// GoStringClone is using the built-in GoString cgo function to
// create a new Go string from the C string.
// It is the caller's responsibility to arrange for the C string to
// eventually be freed, by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the C string,
// in oppose to creating a pointer which references it.
func (s CString) GoStringClone() GoString {
	return GoString(C.GoString((*C.char)(s.data)))
}

// GoBytesAlias create a new []byte slice backed by a C array, without copying the original data.
// Go garbage collector will not interact with this data, and if it is freed from
// the C side of things, the behavior of any Go code using the slice is non-deterministic.
func (s CString) GoStringAlias() string {
	var str string
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	strHeader.Data = uintptr(s.data)
	strHeader.Len = s.Len()

	return str
}

// Len returns the C string length (assuming null-termination byte).
func (s CString) Len() int {
	return cStrLen(s.data)
}

// Free deallocate the C string from the C heap.
func (s CString) Free() {
	cFree(s.data)
}

// CBytes represents a C array allocated in the C heap.
type CBytes struct {
	data unsafe.Pointer // C pointer (allocated using malloc)
	len  int
}

// GoBytesClone is using the built-in GoBytes cgo function to
// create a new Go []byte slice from the C array.
// It is the caller's responsibility to arrange for the C array to
// eventually be freed, by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the C array,
// in oppose to creating a pointer which references it.
func (s CBytes) GoBytesClone() []byte {
	return C.GoBytes(s.data, C.int(s.len))
}

// GoBytesAlias create a new []byte slice backed by a C array, without copying the original data.
// Go garbage collector will not interact with this data, and if it is freed from
// the C side of things, the behavior of any Go code using the slice is non-deterministic.
func (s CBytes) GoBytesAlias() []byte {
	// Arbitrary large-enough size for
	// the array type to hold any len.
	const size = 1 << 30

	p := s.data
	len := s.len

	return (*[size]byte)(p)[:len:len]
}

// Free deallocate the C array from the C heap.
func (s CBytes) Free() {
	cFree(s.data)
}

// CBytes converts an cSvmByteArray struct to a cgo-native C array.
func (ba cSvmByteArray) CBytes() CBytes {
	return CBytes{
		data: unsafe.Pointer(ba.bytes),
		len:  int(ba.length),
	}
}

// svmError converts an SVM byte array to an svmError error.
// The original SVM byte array would be deallocated.
func (ba cSvmByteArray) svmError() error {
	return newSvmError(ba, true)
}

// String helps the cSvmByteArray to implement the Stringer interface.
func (ba cSvmByteArray) String() string {
	return fmt.Sprintf("%s", ba.CBytes().GoBytesAlias())
}

// Free deallocate the cSvmByteArray struct from the C heap.
func (ba cSvmByteArray) Free() {
	cSvmByteArrayDestroy(ba)

	// Should work the same as if calling ba.CBytes().Free().
}

// svmError is error type which represent an error originated in the SVM runtime.
type svmError struct {
	s string
}

// newSvmError creates a new svmError instance from an SVM byte array clone.
// The free param determines whether the original SVM byte array will be de-allocated.
func newSvmError(ba cSvmByteArray, free bool) error {
	clone := ba.CBytes().GoBytesClone()
	if free {
		ba.Free()
	}

	return &svmError{s: string(clone)}
}

// Error helps svmError to implement the error interface.
func (e *svmError) Error() string {
	return fmt.Sprintf("svm error: %v", e.s)
}
