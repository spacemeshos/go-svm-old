package svm

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lsvm_runtime_c_api
// #include "./svm.h"
// #include <string.h>
//
import "C"
import (
	"unsafe"
)

//type cBool C.bool
//type cChar C.char
//type cInt C.int
//type cUint C.uint
//type cUlonglong C.ulonglong
type cSvmByteArray = C.svm_byte_array

var cSuccess = (C.svm_result_t)(C.SVM_SUCCESS)
var cFailure = (C.svm_result_t)(C.SVM_FAILURE)

const SvmAddressLen = 20

func cFree(p unsafe.Pointer) {
	C.free(p)
}

func cSvmByteArrayDestroy(ba cSvmByteArray) {
	C.svm_byte_array_destroy(ba)
}

func cStrLen(p unsafe.Pointer) int {
	return int(C.strlen((*C.char)(p)))
}
