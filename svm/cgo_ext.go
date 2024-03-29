package svm

import "C"
import (
	"reflect"
	"unsafe"
)

// GoBytes is an alias type for []byte slice, used to define local methods.
type GoBytes []byte

// CBytesClone is using the built-in `CBytes` cgo function to
// create a C array that is allocated via the C allocator.
// It is the caller's responsibility to arrange for it to be freed
// via the C allocator by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the underlying array,
// in oppose to creating a pointer which references it.
func (b GoBytes) CBytesClone() CBytes {
	p := C.CBytes(b)
	return CBytes{data: p, len: len(b)}
}

// CBytesAlias creates a new C array backed by the []byte slice, without copying the original data.
//
// ⚠️ UNSAFE.
// Go garbage collector might interact with the []byte slice. Once it will be freed,
// the behavior of any code using the C array is non-deterministic.
func (b GoBytes) CBytesAlias() CBytes {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return CBytes{data: unsafe.Pointer(header.Data), len: len(b)}
}

// CBytes represents a C array allocated via the C allocator.
type CBytes struct {
	data unsafe.Pointer
	len  int
}

func NewCBytes(data unsafe.Pointer, len int) CBytes {
	return CBytes{data, len}
}

// GoBytesClone is using the built-in `GoBytes` cgo function to
// create a new Go []byte slice from the C array.
// It is the caller's responsibility to arrange for the C array to
// eventually be freed via the C allocator, by calling the Free method on it.
// The "Clone" name suffix is to explicitly clarify that it clones the C array,
// in oppose to creating a pointer which references it.
func (s CBytes) GoBytesClone() []byte {
	return C.GoBytes(s.data, C.int(s.len))
}

// GoBytesAlias create a new []byte slice backed by the C array, without copying the original data.
//
// ⚠️ UNSAFE.
// Go garbage collector will not interact with this data. Once it is freed by
// the C allocator, the behavior of any Go code using the slice is non-deterministic.
func (s CBytes) GoBytesAlias() []byte {
	// Arbitrary large-enough size for
	// the array type to hold any len.
	const size = 1 << 30

	p := s.data
	len := s.len

	if p == nil || len == 0 {
		return []byte(nil)
	}

	return (*[size]byte)(p)[:len:len]
}

// Free deallocate the C array via the C allocator.
func (s CBytes) Free() {
	cFree(s.data)
}
