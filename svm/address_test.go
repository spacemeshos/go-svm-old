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

	testRange := AddressLen + 10
	b := make([]byte, testRange, testRange)
	for i := 0; i < testRange; i++ {
		b[i] = byte(i) // Assign some arbitrary value to the next additional byte we're testing.

		if i > AddressLen {
			testAddress(
				bytesToAddress(b[:i]),
				bytesToAddress(b[:AddressLen]),
			)
		} else {
			testAddress(
				bytesToAddress(b[:i]),
				bytesToAddress(b[:i]),
			)
		}
	}
}
