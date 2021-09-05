package svm

// #include "svm.h"
import "C"

type Resource struct {
	raw *C.svm_resource_t
}

func Destroy(r *Resource) {
	C.svm_resource_destroy(r.raw)
}

func Name(r *Resource) string {
	raw := C.svm_resource_type_name_resolve(r.raw.type_id)
	bytes := ByteArrayToSlice(raw)

	return string(bytes)
}
