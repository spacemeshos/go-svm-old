package svm

import (
	"fmt"
	"unsafe"
)

type MemKVStore struct {
	p unsafe.Pointer
}

func NewMemKVStore() (MemKVStore, error) {
	var p unsafe.Pointer
	if res := cSvmMemoryKVCreate(&p); res != cSvmSuccess {
		return MemKVStore{}, fmt.Errorf("failed to create memory kv-store")
	}

	return MemKVStore{p}, nil
}

func (kv MemKVStore) Free() {
	cSvmMemKVDestroy(kv)
}
