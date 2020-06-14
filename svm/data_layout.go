package svm

import (
	"encoding/binary"
)

type DataLayout []uint32

func (dl DataLayout) Encode() []byte {
	buf := make([]byte, len(dl)*4)
	offset := 0

	for _, v := range dl {
		binary.BigEndian.PutUint32(buf[offset:], v)
		offset += 4
	}

	return buf
}
