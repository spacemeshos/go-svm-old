package svm

import "C"
import "errors"

// Error is an svm error wrapper
type Error struct {
	ByteArray
}

// ToError converts SVM error into Golang error
func (e Error) ToError(prefix string) error {
	return errors.New(prefix + string(e.Bytes()))
}

func (e *Error) ptr() *byteArray {
	return &e.byteArray
}
