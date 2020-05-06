package svm

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"unsafe"
)

// TestBytes tests the C/Go bytes array/slice conversion methods.
// It creates and verifies functionality of 5 different types:
// 1) [x]byte array.
// 2) []byte slice from the [x]byte array.
// 3) C array clone from the []byte slice.
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

// TestString tests the C/Go string conversion methods.
// It creates and verifies functionality of 4 different types:
// 1) string.
// 2) C string clone from the string.
// 3) Go string clone from the C string.
// 4) Go string alias from the C string.
func TestString(t *testing.T) {
	req := require.New(t)

	// 1) Allocate a new string.
	str := "Nobility and honour are attached solely to otium and bellum"
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))

	// 2) Use the string to create a new C string.
	cStr := GoString(str).CStringClone()
	// Verify that the C string header data pointer does not equal the string data pointer,
	// but their len equal.
	req.NotEqual(uintptr(cStr.data), strHeader.Data)
	req.Equal(cStr.Len(), len(str))

	// 3) Use the C string to create a new Go string clone.
	goStringClone := cStr.GoStringClone()
	goStringCloneHeader := (*reflect.StringHeader)(unsafe.Pointer(&goStringClone))
	// Verify that the Go string clone header data pointer does not equal the C string data pointer,
	// but their len equal.
	req.NotEqual(uintptr(cStr.data), goStringCloneHeader.Data)
	req.Equal(cStr.Len(), goStringCloneHeader.Len)

	// 4) Use the C string to create a new Go string alias.
	goStringAlias := cStr.GoStringAlias()
	goStringAliasHeader := (*reflect.StringHeader)(unsafe.Pointer(&goStringAlias))
	// Verify that the Go string alias header data pointer and len equal the C array data pointer and len.
	req.Equal(uintptr(cStr.data), goStringAliasHeader.Data)
	req.Equal(cStr.Len(), goStringAliasHeader.Len)

	// Iterate over the original string bytes.
	for i, b := range []byte(str) {
		i := uintptr(i)

		// Verify that both C string and the Go strings i byte equal the original string i byte.
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(cStr.data) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goStringCloneHeader.Data)) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goStringAliasHeader.Data)) + i)))

		// Mutate the C string i byte.
		newVal := uint8(0)
		req.NotEqual(newVal, *(*byte)(unsafe.Pointer(uintptr(cStr.data) + i)))
		*(*byte)(unsafe.Pointer(uintptr(cStr.data) + i)) = newVal
		req.Equal(newVal, *(*byte)(unsafe.Pointer(uintptr(cStr.data) + i)))

		// Verify that the original string i byte isn't affected.
		req.NotEqual(newVal, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(strHeader.Data)) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(strHeader.Data)) + i)))

		// Verify that the Go []byte slice clone i byte isn't affected.
		req.NotEqual(newVal, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goStringCloneHeader.Data)) + i)))
		req.Equal(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goStringCloneHeader.Data)) + i)))

		// Verify that the Go []byte slice alias i byte is affected.
		req.Equal(newVal, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goStringAliasHeader.Data)) + i)))
		req.NotEqual(b, *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(goStringAliasHeader.Data)) + i)))
	}

	// Free the C string which was allocated on the C heap.
	cStr.Free()
}

// TestSvmByteArray tests the C/Go SVM byte array conversion methods.
func TestSvmByteArray(t *testing.T) {
	req := require.New(t)

	testAddress := func(val []byte) {
		ba := GoBytes(val).svmByteArrayClone()
		goBytes := ba.CBytes().GoBytesClone()
		ba.Free()
		req.Equal(val, goBytes)
	}

	testRange := 100
	b := make([]byte, testRange, testRange)
	for i := 0; i < testRange; i++ {
		b[i] = byte(i) // Assign some arbitrary value to the next additional byte we're testing.
		testAddress(b[:i])
	}
}

// TestAddress tests the C/Go SVM address conversion methods.
func TestSvmAddress(t *testing.T) {
	req := require.New(t)

	testAddress := func(in []byte, expectedOut []byte) {
		ba := GoBytes(in).svmAddressClone()
		goBytes := ba.CBytes().GoBytesClone()
		ba.Free()
		req.Equal(expectedOut, goBytes)
	}

	testRange := SvmAddressLen + 10
	b := make([]byte, testRange, testRange)
	for i := 0; i < testRange; i++ {
		b[i] = byte(i) // Assign some arbitrary value to the next additional byte we're testing.

		if i > SvmAddressLen {
			testAddress(b[:i], b[:SvmAddressLen])
		} else {
			testAddress(b[:i], b[:i])
		}
	}
}

// TestAddress tests the C/Go SVM error conversion method.
func TestSvmError(t *testing.T) {
	req := require.New(t)

	ba := GoBytes("some error").svmByteArrayClone()
	goError := ba.svmError()

	req.Equal("svm error: some error", goError.Error())
}
