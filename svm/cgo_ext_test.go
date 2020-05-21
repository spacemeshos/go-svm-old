package svm

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"unsafe"
)

// TestBytes tests the C/Go bytes array/slice conversion methods.
// It creates and verifies functionality of 5 different types:
// 1) Go [x]byte array.
// 2) Go []byte slice from the Go [x]byte array.
// 3) C array clone from the Go []byte slice.
// 4) Go []byte slice clone from the C array.
// 5) Go []byte slice alias from the C array.
func TestBytes(t *testing.T) {
	req := require.New(t)

	// 1) Allocate a new array.
	arr := [3]byte{0x1, 0x02, 0x03}

	// 2) Create a []byte slice from the array.
	slice := arr[:]
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	// Verify that the slice header data pointer does not equal the array data pointer.
	req.Equal(uintptr(unsafe.Pointer(&arr[0])), sliceHeader.Data)

	// 3) Use the slice to create a new C array.
	cBytes := GoBytes(slice).CBytesClone()
	// Verify that the C array header data pointer does not equal the slice data pointer.
	req.NotEqual(uintptr(cBytes.data), sliceHeader.Data)

	// 4) Use the C array to create a new Go []byte slice clone.
	goSliceClone := cBytes.GoBytesClone()
	goSliceCloneHeader := (*reflect.SliceHeader)(unsafe.Pointer(&goSliceClone))
	// Verify that the Go []byte slice clone header data pointer does not equal the C array data pointer.
	req.NotEqual(uintptr(cBytes.data), goSliceCloneHeader.Data)

	// 5) Use the C array to create a new Go []byte slice alias.
	goSliceAlias := cBytes.GoBytesAlias()
	goSliceAliasHeader := (*reflect.SliceHeader)(unsafe.Pointer(&goSliceAlias))
	// Verify that the Go []byte slice alias header data pointer equal the C array data pointer.
	req.Equal(uintptr(cBytes.data), goSliceAliasHeader.Data)

	// Iterate over the original slice bytes.
	for i, b := range slice {
		i := uintptr(i)

		// Verify that both C array and the Go []byte slices i byte equal the original slice i byte.
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(cBytes.data) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goSliceCloneHeader.Data)) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goSliceAliasHeader.Data)) + i)))

		// Mutate the C array i byte.
		newVal := uint8(0)
		req.NotEqual(newVal, *(*byte)(unsafe.Pointer(uintptr(cBytes.data) + i)))
		*(*byte)(unsafe.Pointer(uintptr(cBytes.data) + i)) = newVal
		req.Equal(newVal, *(*byte)(unsafe.Pointer(uintptr(cBytes.data) + i)))

		// Verify that the original slice i byte isn't affected.
		req.NotEqual(newVal, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(sliceHeader.Data)) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(sliceHeader.Data)) + i)))

		// Verify that the Go []byte slice clone i byte isn't affected.
		req.NotEqual(newVal, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goSliceCloneHeader.Data)) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goSliceCloneHeader.Data)) + i)))

		// Verify that the Go []byte slice alias i byte is affected.
		req.Equal(newVal, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goSliceAliasHeader.Data)) + i)))
		req.NotEqual(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goSliceAliasHeader.Data)) + i)))
	}

	// Free the C array which was allocated on the C heap.
	cBytes.Free()
}
