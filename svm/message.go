package svm

/*
#include "svm.h"
*/
import "C"
import (
	"encoding/binary"
	"errors"
)

var errMessageTooLong = errors.New("message too long")

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

// CallMessage is an message abstraction for Runtime.Call endpoint
type CallMessage struct {
	Version    int
	Target     Address
	FuncName   string
	VerifyData []byte
	CallData   []byte
}

// Encode converts call message abstraction into svm_byte_array wrapped by Message
func (cm CallMessage) Encode() (*Message, error) {
	p := 0
	fn := []byte(cm.FuncName)
	bs := make([]byte, 2+addressLength+1+len(fn)+1+len(cm.VerifyData)+1+len(cm.CallData))
	encode := func(b []byte) error {
		ln := len(b)
		if ln > 255 {
			return errMessageTooLong
		}
		bs[p] = uint8(ln)
		copy(bs[p+1:p+1+ln], b)
		p += 1 + ln
		return nil
	}
	binary.BigEndian.PutUint16(bs[p:p+2], uint16(cm.Version))
	p += 2
	copy(bs[p:p+addressLength], cm.Target[:])
	p += addressLength
	if err := encode(fn); err != nil {
		return nil, err
	}
	if err := encode(cm.VerifyData); err != nil {
		return nil, err
	}
	if err := encode(cm.CallData); err != nil {
		return nil, err
	}
	return NewMessageFromBytes(bs), nil
}

// SpawnMessage is am message abstraction for Runtime.Spawn endpoint
type SpawnMessage struct {
	Version  int
	Template Address
	Name     string
	Ctor     string
	CallData []byte
}

// Encode converts spaen message abstraction into svm_byte_array wrapped by Message
func (cm SpawnMessage) Encode() (*Message, error) {
	p := 0
	fn := []byte(cm.Name)
	ct := []byte(cm.Ctor)
	bs := make([]byte, 2+addressLength+1+len(fn)+1+len(ct)+1+len(cm.CallData))
	encode := func(b []byte) error {
		ln := len(b)
		if ln > 255 {
			return errMessageTooLong
		}
		bs[p] = uint8(ln)
		copy(bs[p+1:p+1+ln], b)
		p += 1 + ln
		return nil
	}
	binary.BigEndian.PutUint16(bs[p:p+2], uint16(cm.Version))
	p += 2
	copy(bs[p:p+addressLength], cm.Template[:])
	p += addressLength
	if err := encode(fn); err != nil {
		return nil, err
	}
	if err := encode(ct); err != nil {
		return nil, err
	}
	if err := encode(cm.CallData); err != nil {
		return nil, err
	}
	return NewMessageFromBytes(bs), nil
}

// DeployMessage is an abstraction for Runtime.Deploy endpoint
type DeployMessage []section

// Encode converts deply message abstraction into svm_byte_array wrapped by Message
func (dm DeployMessage) Encode() (*Message, error) {
	var err error
	sectionsLength := 0
	for _, s := range dm {
		sectionsLength += s.length()
	}
	bs := make([]byte, 2, 2+previewSectionLength*len(dm)+sectionsLength)
	binary.BigEndian.PutUint16(bs, uint16(len(dm)))
	for _, s := range dm {
		bs, err = encodePreview(bs, s)
		if err != nil {
			return nil, err
		}
		bs, err = s.encode(bs)
		if err != nil {
			return nil, err
		}
	}
	return NewMessageFromBytes(bs), nil
}
