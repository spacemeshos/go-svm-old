package svm

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lsvm_runtime_c_api
// #include "./svm.h"
// #include <string.h>
//
import "C"

import (
	"unsafe"
)

type (
	cUchar      = C.uchar
	cUint       = C.uint
	cSvmResultT = C.svm_result_t
)

const (
	cSvmSuccess = (C.svm_result_t)(C.SVM_SUCCESS)
)

func cSvmImportsAlloc(p *unsafe.Pointer, count uint) cSvmResultT {
	return (cSvmResultT)(C.svm_imports_alloc(p, C.uint(count)))
}

func cSvmMemoryRuntimeCreate(runtime *unsafe.Pointer, kv, imports unsafe.Pointer) error {
	err := cSvmByteArray{}
	defer err.SvmFree()

	if res := C.svm_memory_runtime_create(
		runtime,
		kv,
		imports,
		&err,
	); res != cSvmSuccess {
		return err.svmError()
	}

	return nil
}

func cSvmMemoryKVCreate(p *unsafe.Pointer) cSvmResultT {
	return (cSvmResultT)(C.svm_memory_state_kv_create(p))
}

func cSvmValidateTemplate(runtime Runtime, appTemplate []byte) error {
	cRuntime := runtime._inner
	cAppTemplate := bytesCloneToSvmByteArray(appTemplate)
	cErr := cSvmByteArray{}

	defer func() {
		cAppTemplate.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_validate_template(
		cRuntime,
		cAppTemplate,
		&cErr,
	); res != cSvmSuccess {
		return cErr.svmError()
	}

	return nil
}

func cSvmDeployTemplate(runtime Runtime, appTemplate []byte, author Address, gasMetering bool, gasLimit uint64) ([]byte, error) {
	cReceipt := cSvmByteArray{}
	cRuntime := runtime._inner
	cAppTemplate := bytesCloneToSvmByteArray(appTemplate)
	cAuthor := bytesCloneToSvmByteArray(author[:])
	cGasMetering := C.bool(gasMetering)
	cGasLimit := C.uint64_t(gasLimit)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.SvmFree()
		cAppTemplate.Free()
		cAuthor.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_deploy_template(
		&cReceipt,
		cRuntime,
		cAppTemplate,
		cAuthor,
		cGasMetering,
		cGasLimit,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cReceipt), nil
}

func cSvmSpawnApp(runtime Runtime, spawnApp []byte, creator Address, gasMetering bool, gasLimit uint64) ([]byte, error) {
	cReceipt := cSvmByteArray{}
	cRuntime := runtime._inner
	cSpawnApp := bytesCloneToSvmByteArray(spawnApp)
	cCreator := bytesCloneToSvmByteArray(creator[:])
	cGasMetering := C.bool(gasMetering)
	cGasLimit := C.uint64_t(gasLimit)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.SvmFree()
		cSpawnApp.Free()
		cCreator.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_spawn_app(
		&cReceipt,
		cRuntime,
		cSpawnApp,
		cCreator,
		cGasMetering,
		cGasLimit,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cReceipt), nil
}

func cSvmValidateApp(runtime Runtime, app []byte) error {
	cRuntime := runtime._inner
	cApp := bytesCloneToSvmByteArray(app)
	cErr := cSvmByteArray{}

	defer func() {
		cApp.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_validate_app(
		cRuntime,
		cApp,
		&cErr,
	); res != cSvmSuccess {
		return cErr.svmError()
	}

	return nil
}

func cSvmValidateTx(runtime Runtime, appTx []byte) (Address, error) {
	cAppAddr := cSvmByteArray{}
	cRuntime := runtime._inner
	cAppTx := bytesCloneToSvmByteArray(appTx)
	cErr := cSvmByteArray{}

	defer func() {
		cAppAddr.SvmFree()
		cAppTx.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_validate_tx(
		&cAppAddr,
		cRuntime,
		cAppTx,
		&cErr,
	); res != cSvmSuccess {
		return Address{}, cErr.svmError()
	}

	return svmByteArrayCloneToAddress(cAppAddr), nil
}

func cSvmExecApp(runtime Runtime, tx []byte, appState []byte, gasMetering bool,
	gasLimit uint64) ([]byte, error) {
	cReceipt := cSvmByteArray{}
	cRuntime := runtime._inner
	cTx := bytesCloneToSvmByteArray(tx)
	cAppState := bytesCloneToSvmByteArray(appState)
	cGasMetering := C.bool(gasMetering)
	cGasLimit := C.uint64_t(gasLimit)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.SvmFree()
		cTx.Free()
		cAppState.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_exec_app(
		&cReceipt,
		cRuntime,
		cTx,
		cAppState,
		cGasMetering,
		cGasLimit,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cReceipt), nil
}

func cSvmByteArrayDestroy(ba cSvmByteArray) {
	C.svm_byte_array_destroy(ba)
}

func cSvmRuntimeDestroy(runtime Runtime) {
	C.svm_runtime_destroy(runtime._inner)
}

func cSvmImportsDestroy(imports Imports) {
	C.svm_imports_destroy(imports._inner)
}

func cSvmStateKVDestroy(_inner unsafe.Pointer) {
	C.svm_state_kv_destroy(_inner)
}

func cFree(p unsafe.Pointer) {
	C.free(p)
}
