package svm

// #include "./svm.h"
//
import "C"
import (
	"fmt"
	"runtime"
	"sync"
)

// `svm_trampoline` is used to facilitate the invocation host import functions,
// written in Go, directly from a WebAssembly module instance.
//
// Schematically, the workflow looks like this:
//
//     +------------------------------+
//     |                              |
//     |  +------------------------+  |
//     |  |                        |  |
//     |  |   WebAssembly module   |  |
//     |  |  host func invocation  |  |
//     |  |                        |  |
//     |  +----+--------------+----+  |
//     |       |              ^       |
//     |       v              |       |
//     |  +----+--------------+----+  |
//     |  |                        |  |
//     |  |     Wasmer handler     |  |
//     |  |                        |  |
//     |  +----+--------------+----+  |
//     |       |              ^       |
//     |       v              |       |
//     |  +----+--------------+----+  |
//     |  |                        |  |
//     |  |      SVM handler       |  |
//     |  |                        |  |
//     |  +----+--------------+----+  |
//     |       |              ^       |
//     |       v              |       |
//     |  +----+--------------+----+  |
//     |  |                        |  |
//     |  |    `svm_trampoline`    |  |
//     |  |						   |  |
//     |  |	     (cgo handler)     |  |
//     |  |                        |  |
//     |  +----+--------------+----+  |
//     |       |              ^       |
//     |       v              |       |
//     |  +----+--------------+----+  |
//     |  |                        |  |
//     |  |    host Go function    |  |
//     |  |                        |  |
//     |  +------------------------+  |
//     |                              |
//     +------------------------------+

//export svm_trampoline
func svm_trampoline(env *C.svm_env_t, args *C.svm_byte_array, results *C.svm_byte_array) *C.svm_byte_array {
	// Fetch the target function.
	hostEnv := (*functionEnvironment)(env.host_env)
	f := hostFunctionStore.get(hostEnv.hostFunctionStoreIndex)
	if f == nil {
		panic(fmt.Sprintf("go-svm: host import function not found; index: %v", hostEnv.hostFunctionStoreIndex))
	}

	// Decode args.
	goArgs := Values{}
	if err := goArgs.Decode(svmByteArrayCloneToBytes(*args)); err != nil {
		panic(fmt.Sprintf("go-svm: %v", err))
	}

	// Invoke.
	goResults, err := f(goArgs)
	if err != nil {
		err := []byte(err.Error())
		cErr := bytesAliasToSvmByteArray(err)

		// Re-allocate the error on SVM side, so that it
		// would be able do de-allocate it once processing is done.
		cSvmErr := cSvmWasmErrorCreate(cErr)
		runtime.KeepAlive(err)

		return cSvmErr
	}

	// Encode results.
	rawResults := Values(goResults).Encode()
	*results = bytesCloneToSvmByteArray(rawResults)

	return nil
}

func cSvmWasmErrorCreate(err cSvmByteArray) *cSvmByteArray {
	return (*cSvmByteArray)(C.svm_wasm_error_create(err))
}

type functionEnvironment struct {
	hostFunctionStoreIndex uint
}

type hostFunction func([]Value) ([]Value, error)

// hostFunctionStore is a static container for the registered import functions, written in Go,
// to be invoked from the unsafe, cgo-exported `svm_trampoline`.
var hostFunctionStore = hostFunctions{
	functions: make(map[uint]hostFunction),
}

type hostFunctions struct {
	sync.RWMutex
	functions map[uint]hostFunction
}

func (hf *hostFunctions) get(index uint) hostFunction {
	return hf.functions[index]
}

func (hf *hostFunctions) add(function hostFunction) uint {
	hf.Lock()
	defer hf.Unlock()

	index := uint(len(hf.functions))
	hf.functions[index] = function

	return index
}
