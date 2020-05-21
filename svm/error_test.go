package svm

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSvmError(t *testing.T) {
	req := require.New(t)

	err := bytesCloneToSvmByteArray([]byte("Mayday"))
	goError := err.svmError()

	req.Equal("svm error: Mayday", goError.Error())
}
