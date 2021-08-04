package svm

import (
	"fmt"
	"unsafe"
)

type StateKV_Mem struct {
	// _inner is a pointer to an SVM-managed heap allocation.
	_inner unsafe.Pointer
}

func NewStateKV_Mem() (StateKV_Mem, error) {
	var p unsafe.Pointer
	if res := cSvmMemoryKVCreate(&p); res != cSvmSuccess {
		return StateKV_Mem{}, fmt.Errorf("failed to create memory state KV store")
	}

	return StateKV_Mem{p}, nil
}

func (kv StateKV_Mem) Free() {
	cSvmStateKVDestroy(kv._inner)
}
