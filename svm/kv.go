package svm

import (
	"fmt"
	"unsafe"
)

func NewMemoryKVStore() (unsafe.Pointer, error) {
	var p unsafe.Pointer
	if res := cSvmMemoryKVCreate(&p); res != cSvmSuccess {
		return nil, fmt.Errorf("failed to create memory kv-store")
	}
	return p, nil
}

func NewMemoryRawKVStore() (unsafe.Pointer, error) {
	var p unsafe.Pointer
	if res := cSvmMemoryRawKVCreate(&p); res != cSvmSuccess {
		return nil, fmt.Errorf("failed to create memory kv-store")
	}
	return p, nil
}
