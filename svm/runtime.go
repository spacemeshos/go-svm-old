package svm

/*
#include "svm.h"
#include "memory.h"
*/
import "C"
import "unsafe"

// Runtime is a wrapper for svm_runtime
type Runtime struct {
	svmRuntime unsafe.Pointer
}

// NewRuntime creates new SVM runtime
func NewRuntime() (*Runtime, error) {
	rt := &Runtime{}
	err := Error{}
	if res := C.svm_memory_runtime_create(&rt.svmRuntime, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("can't create new SVM runtime: ")
	}
	return rt, nil
}

// Destroy releases SVM runtime
func (rt *Runtime) Destroy() {
	if rt.svmRuntime != nil {
		C.svm_runtime_destroy(rt.svmRuntime)
		rt.svmRuntime = nil
	}
}

// Close is an idiomatic way to release runtime
func (rt *Runtime) Close() error {
	rt.Destroy()
	return nil
}

// ValidateDeploy is a wrapper to svm_validate_deploy endpoint
func (rt *Runtime) ValidateDeploy(msg *Message) error {
	err := Error{}
	if res := C.svm_validate_deploy(rt.svmRuntime, msg.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return err.ToError("failed to validate deploy contract: ")
	}
	return nil
}
