package svm

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func randomByteArray(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func Test_ByteArrayFromBytes(t *testing.T) {
	length := 100
	message := NewMessage(length)
	defer message.Destroy()

	bs1 := randomByteArray(length - 1)
	e := message.FromBytes(bs1)
	assert.NoError(t, e)
	assert.Equal(t, message.Bytes(), bs1)

	bs2 := randomByteArray(length)
	e = message.FromBytes(bs2)
	assert.NoError(t, e)
	assert.Equal(t, message.Bytes(), bs2)

	bs3 := randomByteArray(length + 1)
	e = message.FromBytes(bs3)
	assert.Error(t, e)
}
