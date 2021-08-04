package svm

import "C"
import (
	"fmt"
	"unsafe"
)

type Runtime struct {
	// _inner is a pointer to an SVM-managed heap allocation.
	_inner unsafe.Pointer
}

func (r Runtime) Free() {
	cSvmRuntimeDestroy(r)
}

type RuntimeBuilder struct {
	imports unsafe.Pointer
	kv      unsafe.Pointer
	host    unsafe.Pointer
}

func NewRuntimeBuilder() RuntimeBuilder {
	return RuntimeBuilder{}
}

func (rb RuntimeBuilder) WithImports(imports *Imports) RuntimeBuilder {
	rb.imports = imports._inner
	return rb
}

func (rb RuntimeBuilder) WithStateKV_Mem(kv *StateKV_Mem) RuntimeBuilder {
	rb.kv = kv._inner
	return rb
}

func (rb RuntimeBuilder) WithStateKV_FFI(kv *StateKV_FFI) RuntimeBuilder {
	rb.kv = kv._inner
	return rb
}

func (rb RuntimeBuilder) Build() (Runtime, error) {
	var p unsafe.Pointer

	if err := cSvmMemoryRuntimeCreate(
		&p,
		rb.kv,
		rb.imports,
	); err != nil {
		return Runtime{}, fmt.Errorf("failed to create runtime: %v", err)
	}

	return Runtime{p}, nil
}
