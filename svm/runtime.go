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

// ValidateCall is a wrapped to svm_validate_call endpoint
func (rt *Runtime) ValidateCall(msg *Message) error {
	err := Error{}
	if res := C.svm_validate_call(rt.svmRuntime, msg.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return err.ToError("failed to validate call message: ")
	}
	return nil
}

// Call is a wrapper to svm_call endpoint
func (rt *Runtime) Call(envelope *Envelope, msg *Message, ctx *Context) (*CallReceipt, error) {
	var err Error
	rcpt := &ByteArray{}
	defer rcpt.Destroy()
	/*if res := C.svm_validate_call(rt.svmRuntime, msg.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("failed to validate call message: ")
	}*/
	if res := C.svm_call(&rcpt.byteArray, rt.svmRuntime, envelope.byteArray, msg.byteArray, ctx.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("failed to validate call message: ")
	}
	cr := &CallReceipt{}
	if err := cr.Decode(rcpt.Bytes()); err != nil {
		return nil, err
	}
	return cr, nil
}

// ValidateSpawn is a wrapper to svm_validate_spawn endpoint
func (rt *Runtime) ValidateSpawn(msg *Message) error {
	err := Error{}
	if res := C.svm_validate_call(rt.svmRuntime, msg.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return err.ToError("failed to validate spawn contract: ")
	}
	return nil
}

// Spawn is a wrapper to svm_spawn endpoint
func (rt *Runtime) Spawn(envelope *Envelope, msg *Message, ctx *Context) (*SpawnReceipt, error) {
	var err Error
	rcpt := &ByteArray{}
	defer rcpt.Destroy()
	/*if res := C.svm_validate_spawn(rt.svmRuntime, msg.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("failed to validate spawn message: ")
	}*/
	if res := C.svm_spawn(&rcpt.byteArray, rt.svmRuntime, envelope.byteArray, msg.byteArray, ctx.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("failed spawn contract: ")
	}
	sr := &SpawnReceipt{}
	if err := sr.Decode(rcpt.Bytes()); err != nil {
		return nil, err
	}
	return sr, nil
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

// Deploy is a wrapper to svm_deploy endpoint
func (rt *Runtime) Deploy(envelope *Envelope, msg *Message, ctx *Context) (*DeployReceipt, error) {
	var err Error
	rcpt := &ByteArray{}
	defer rcpt.Destroy()
	/*if res := C.svm_validate_deploy(rt.svmRuntime, msg.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("failed to validate deploy message: ")
	}*/
	if res := C.svm_deploy(&rcpt.byteArray, rt.svmRuntime, envelope.byteArray, msg.byteArray, ctx.byteArray, err.ptr()); res != C.SVM_SUCCESS {
		defer err.Destroy()
		return nil, err.ToError("failed deploy contract: ")
	}
	dr := &DeployReceipt{}
	if err := dr.Decode(rcpt.Bytes()); err != nil {
		return nil, err
	}
	return dr, nil
}
