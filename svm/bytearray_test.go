package svm

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func genbytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func Test_ByteArrayFromBytes(t *testing.T) {
	balen := 100
	ba := NewMessage(balen)
	defer ba.Destroy()
	bs1 := genbytes(balen - 1)
	e := ba.FromBytes(bs1)
	assert.NoError(t, e)
	assert.Equal(t, ba.Bytes(), bs1)
	bs2 := genbytes(balen)
	e = ba.FromBytes(bs2)
	assert.NoError(t, e)
	assert.Equal(t, ba.Bytes(), bs2)
	bs3 := genbytes(balen + 1)
	e = ba.FromBytes(bs3)
	assert.Error(t, e)
}
