package svm

import "C"
import "errors"

type Error struct {
	ByteArray
}

func (e Error) ToError(prefix string) error {
	return errors.New(prefix + string(e.Bytes()))
}

func (e *Error) ptr() *byteArray {
	return &e.byteArray
}
