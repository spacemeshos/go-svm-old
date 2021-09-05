package svm

// #include "svm.h"
import "C"

// Message is an SVM message wrapper
type Message struct {
	ByteArray
}

// NewMessage allocated a new `Message` of 
func NewMessage(length int) *Message {
	msg := &Message{}
	msg.byteArray = C.svm_message_alloc(C.uint(length))
	return msg
}

// NewMessageFromBytes creates new SVM message from specified bytes
func NewMessageFromBytes(bytes []byte) *Message {
	msg := NewMessage(len(bytes))
	msg.FromBytes(bytes)
	return msg
}
