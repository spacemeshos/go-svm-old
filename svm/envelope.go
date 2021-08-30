package svm

/*
#include "svm.h"
#include "memory.h"
*/
import "C"
import "encoding/binary"

type Envelope struct {
	ByteArray
}

func DefaultEnvelope(principal Address) *Envelope {
	return NewEnvelope(principal, 0, 0, 0)
}

func NewEnvelope(principal Address, amount, gasLimit, gasFee uint64) *Envelope {
	e := &Envelope{}
	e.byteArray = C.svm_envelope_alloc()
	bs := make([]byte, AddressLength+3*8)
	copy(bs[:AddressLength], principal[:])
	p := AddressLength
	binary.BigEndian.PutUint64(bs[p:p+8], amount)
	p += 8
	binary.BigEndian.PutUint64(bs[p:p+8], gasLimit)
	p += 8
	binary.BigEndian.PutUint64(bs[p:p+8], gasFee)
	e.FromBytes(bs)
	return e
}
