package svm

// #include "svm.h"
import "C"
import "unsafe"

type ResourcesIter struct {
	raw unsafe.Pointer
}

func TotalLiveResources() uint32 {
	return (uint32)(C.svm_total_live_resources())
}

func NewResourcesIter() *ResourcesIter {
	raw := C.svm_resource_iter_new()
	iter := ResourcesIter{raw}

	return &iter
}

func Next(iter *ResourcesIter) *Resource {
	raw := C.svm_resource_iter_next(iter.raw)

	if raw != nil {
		return &Resource{raw}
	}
	return nil
}

func ResourcesIterDestroy(iter *ResourcesIter) {
	C.svm_resource_iter_destroy(iter.raw)
}
