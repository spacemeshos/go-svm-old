package svm

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSvmByteArray(t *testing.T) {
	req := require.New(t)

	testAddress := func(val []byte) {
		ba := bytesCloneToSvmByteArray(val)
		goBytes := svmByteArrayCloneToBytes(ba)
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

func TestBytesCloneToSvmByteArray(t *testing.T) {
	req := require.New(t)

	ba := bytesCloneToSvmByteArray(nil)
	req.Equal("<nil>", fmt.Sprintf("%v", ba.bytes))
	req.Equal("0", fmt.Sprintf("%v", ba.length))
	req.Equal("", fmt.Sprintf("%v", ba.String()))

	ba = bytesCloneToSvmByteArray(make([]byte, 0))
	req.Equal("<nil>", fmt.Sprintf("%v", ba.bytes))
	req.Equal("0", fmt.Sprintf("%v", ba.length))
	req.Equal("", fmt.Sprintf("%v", ba.String()))

	ba = bytesCloneToSvmByteArray(make([]byte, 1))
	req.Contains(fmt.Sprintf("%v", ba.bytes), "0x")
	req.Equal("1", fmt.Sprintf("%v", ba.length))
	req.Equal("\x00", fmt.Sprintf("%v", ba.String()))
	ba.Free()
}

func TestSvmByteArrayCloneToBytes(t *testing.T) {
	req := require.New(t)

	b := svmByteArrayCloneToBytes(cSvmByteArray{})
	req.Equal(make([]byte, 0), b)
}
