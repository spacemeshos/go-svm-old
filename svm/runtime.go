package svm

import "C"
import (
	"fmt"
	"unsafe"
)

type Runtime struct {
	p unsafe.Pointer
}

func (r Runtime) Free() {
	cSvmRuntimeDestroy(r)
}

type RuntimeBuilder struct {
	imports    unsafe.Pointer
	memKV      unsafe.Pointer
	diskKVPath string
	host       unsafe.Pointer
}

func NewRuntimeBuilder() RuntimeBuilder {
	return RuntimeBuilder{}
}

func (rb RuntimeBuilder) WithImports(imports Imports) RuntimeBuilder {
	rb.imports = imports.p
	return rb
}

func (rb RuntimeBuilder) WithMemKVStore(kv MemKVStore) RuntimeBuilder {
	rb.memKV = kv.p
	return rb
}

func (rb RuntimeBuilder) WithDiskKV(path string) RuntimeBuilder {
	rb.diskKVPath = path
	return rb
}

func (rb RuntimeBuilder) WithHost(p unsafe.Pointer) RuntimeBuilder {
	rb.host = p
	return rb
}

func (rb RuntimeBuilder) Build() (Runtime, error) {
	var p unsafe.Pointer

	if err := cSvmMemoryRuntimeCreate(
		&p,
		rb.memKV,
		rb.host,
		rb.imports,
	); err != nil {
		return Runtime{}, fmt.Errorf("failed to create runtime: %v", err)
	}

	return Runtime{p}, nil
}

func InstanceContextHostGet(ctx unsafe.Pointer) unsafe.Pointer {
	return cSvmInstanceContextHostGet(ctx)
}
