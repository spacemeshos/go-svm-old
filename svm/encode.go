package svm

import "encoding/binary"

func Encode16be(bs []byte, val uint16) []byte {
	b := []byte{0, 0}
	binary.BigEndian.PutUint16(b[:], val)
	return append(bs, b...)
}

func Encode32be(bs []byte, val uint32) []byte {
	b := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(b[:], val)
	return append(bs, b...)
}

func Encode64be(bs []byte, val uint64) []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(b[:], val)
	return append(bs, b...)
}
