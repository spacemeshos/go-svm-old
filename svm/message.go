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

type Message struct {
	ByteArray
}

func NewMessage(length int) *Message {
	m := &Message{}
	m.byteArray = C.svm_message_alloc(C.uint(length))
	return m
}

func NewMessageFromBytes(bs []byte) *Message {
	m := NewMessage(len(bs))
	m.FromBytes(bs)
	return m
}

type CallMessage struct {
	Version    int
	Target     Address
	FuncName   string
	VerifyData []byte
	CallData   []byte
}

func (cm CallMessage) Encode() (*Message, error) {
	p := 0
	fn := []byte(cm.FuncName)
	bs := make([]byte, 2+AddressLength+1+len(fn)+1+len(cm.VerifyData)+1+len(cm.CallData))
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
	copy(bs[p:p+AddressLength], cm.Target[:])
	p += AddressLength
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

type SpawnMessage struct {
	Version  int
	Template Address
	Name     string
	Ctor     string
	CallData []byte
}

func (cm SpawnMessage) Encode() (*Message, error) {
	p := 0
	fn := []byte(cm.Name)
	ct := []byte(cm.Ctor)
	bs := make([]byte, 2+AddressLength+1+len(fn)+1+len(ct)+1+len(cm.CallData))
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
	copy(bs[p:p+AddressLength], cm.Template[:])
	p += AddressLength
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

type DeployMessage []Section

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
