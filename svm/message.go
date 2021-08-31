package svm

/*
#include "svm.h"
*/
import "C"

// Message is an SVM message wrapper
type Message struct {
	ByteArray
}

// NewMessage creates new message with specified data length
func NewMessage(length int) *Message {
	m := &Message{}
	m.byteArray = C.svm_message_alloc(C.uint(length))
	return m
}

// NewMessageFromBytes creates new SVM message from specified bytes
func NewMessageFromBytes(bs []byte) *Message {
	m := NewMessage(len(bs))
	m.FromBytes(bs)
	return m
}
