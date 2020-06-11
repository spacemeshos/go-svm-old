package svm

import "fmt"

// svmError is error type which represent an error originated in the SVM runtime.
type svmError struct {
	s string
}

// newSvmError creates a new svmError instance from []byte slice.
func newSvmError(b []byte) error {
	return &svmError{s: string(b)}
}

// Error helps svmError to implement the error interface.
func (e *svmError) Error() string {
	return fmt.Sprintf("svm error: %v", e.s)
}
