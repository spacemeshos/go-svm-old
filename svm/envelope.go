package svm

/*
#include "svm.h"
#include "memory.h"
*/
import "C"
import "encoding/binary"

// Envelope is an SVM envelope wrapper
type Envelope struct {
	ByteArray
}

// DefaultEnvelope returns default envelope structore with specified principal address
func DefaultEnvelope(principal Address) *Envelope {
	return NewEnvelope(principal, 0, 0, 0)
}

// NewEnvelope creates envelope for whole specified paremeters
func NewEnvelope(principal Address, amount, gasLimit, gasFee uint64) *Envelope {
	e := &Envelope{}
	e.byteArray = C.svm_envelope_alloc()
	bs := make([]byte, addressLength+3*8)
	copy(bs[:addressLength], principal[:])
	p := addressLength
	binary.BigEndian.PutUint64(bs[p:p+8], amount)
	p += 8
	binary.BigEndian.PutUint64(bs[p:p+8], gasLimit)
	p += 8
	binary.BigEndian.PutUint64(bs[p:p+8], gasFee)
	e.FromBytes(bs)
	return e
}
