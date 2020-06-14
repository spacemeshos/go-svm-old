package svm

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDataLayout_Encode(t *testing.T) {
	req := require.New(t)

	dl := DataLayout([]uint32{10, 20, 30})
	b := dl.Encode()
	req.Len(b, 12)
	req.Equal(uint32(10), binary.BigEndian.Uint32(b[:]))
	req.Equal(uint32(20), binary.BigEndian.Uint32(b[4:]))
	req.Equal(uint32(30), binary.BigEndian.Uint32(b[8:]))
}
