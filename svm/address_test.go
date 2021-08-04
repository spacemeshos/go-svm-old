package svm

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddress(t *testing.T) {
	req := require.New(t)

	testAddress := func(in, expectedOut Address) {
		ba := bytesCloneToSvmByteArray(in[:])
		out := svmByteArrayCloneToAddress(ba)
		ba.Free()
		req.Equal(expectedOut, out)
	}

	testRange := AddressSize + 10
	b := make([]byte, testRange, testRange)
	for i := 0; i < testRange; i++ {
		b[i] = byte(i) // Assign some arbitrary value to the next additional byte we're testing.

		if i > AddressSize {
			testAddress(
				BytesToAddress(b[:i]),
				BytesToAddress(b[:AddressSize]),
			)
		} else {
			testAddress(
				BytesToAddress(b[:i]),
				BytesToAddress(b[:i]),
			)
		}
	}
}
