package svm

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR} -lsvm_runtime_c_api
// #include "./svm.h"
// #include <string.h>
//
import "C"

import (
	"fmt"
	"unsafe"
)

type cUchar = C.uchar
type cUint = C.uint
type cSvmByteArray = C.svm_byte_array
type cSvmResultT = C.svm_result_t

const cSvmSuccess = (C.svm_result_t)(C.SVM_SUCCESS)

func cSvmImportsAlloc(imports *unsafe.Pointer, count uint) cSvmResultT {
	return (cSvmResultT)(C.svm_imports_alloc(imports, C.uint(count)))
}

func cSvmImportFuncBuild(
	imports Imports,
	moduleName string,
	importName string,
	cgoFuncPointer unsafe.Pointer,
	params ValueTypes,
	returns ValueTypes,
) error {
	cImports := imports.p
	cModuleName := bytesCloneToSvmByteArray([]byte(moduleName))
	cImportName := bytesCloneToSvmByteArray([]byte(importName))
	cParams := bytesCloneToSvmByteArray(params.Encode())
	cReturns := bytesCloneToSvmByteArray(returns.Encode())
	cErr := cSvmByteArray{}

	defer func() {
		cModuleName.Free()
		cImportName.Free()
		cParams.Free()
		cReturns.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_import_func_build(
		cImports,
		cModuleName,
		cImportName,
		cgoFuncPointer,
		cParams,
		cReturns,
		&cErr,
	); res != cSvmSuccess {
		return cErr.svmError()
	}

	return nil
}

func cSvmMemoryRuntimeCreate(runtime *unsafe.Pointer, kv, host, imports unsafe.Pointer) error {
	err := cSvmByteArray{}
	defer err.SvmFree()

	if res := C.svm_memory_runtime_create(
		runtime,
		kv,
		host,
		imports,
		&err,
	); res != cSvmSuccess {
		return err.svmError()
	}

	return nil
}

func cSvmMemoryKVCreate(p *unsafe.Pointer) cSvmResultT {
	return (cSvmResultT)(C.svm_memory_kv_create(p))
}

func cSvmEncodeAppTemplate(version int, name string, code []byte, dataLayout DataLayout) ([]byte, error) {
	appTemplate := cSvmByteArray{}
	cName := bytesCloneToSvmByteArray([]byte(name))
	cCode := bytesCloneToSvmByteArray(code)
	cDataLayout := bytesCloneToSvmByteArray(dataLayout.Encode())
	err := cSvmByteArray{}

	defer func() {
		appTemplate.SvmFree()
		cName.Free()
		cCode.Free()
		cDataLayout.Free()
		err.SvmFree()
	}()

	if res := C.svm_encode_app_template(
		&appTemplate,
		C.uint(version),
		cName,
		cCode,
		cDataLayout,
		&err,
	); res != cSvmSuccess {
		return nil, err.svmError()
	}

	return svmByteArrayCloneToBytes(appTemplate), nil
}

func cSvmValidateTemplate(runtime Runtime, appTemplate []byte) error {
	cRuntime := runtime.p
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

func cSvmDeployTemplate(runtime Runtime, appTemplate []byte, author Address, hostCtx []byte, gasMetering bool, gasLimit uint64) ([]byte, error) {
	cReceipt := cSvmByteArray{}
	cRuntime := runtime.p
	cAppTemplate := bytesCloneToSvmByteArray(appTemplate)
	cAuthor := bytesCloneToSvmByteArray(author[:])
	cHostCtx := bytesCloneToSvmByteArray(hostCtx)
	cGasMetering := C.bool(gasMetering)
	cGasLimit := C.uint64_t(gasLimit)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.SvmFree()
		cAppTemplate.Free()
		cAuthor.Free()
		cHostCtx.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_deploy_template(
		&cReceipt,
		cRuntime,
		cAppTemplate,
		cAuthor,
		cHostCtx,
		cGasMetering,
		cGasLimit,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cReceipt), nil
}

func cSvmTemplateReceiptAddr(receipt []byte) (Address, error) {
	cTemplateAddr := cSvmByteArray{}
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cTemplateAddr.SvmFree()
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_template_receipt_addr(
		&cTemplateAddr,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return Address{}, cErr.svmError()
	}

	return svmByteArrayCloneToAddress(cTemplateAddr), nil
}

func cSvmTemplateReceiptGas(receipt []byte) (uint64, error) {
	var cGasUsed C.uint64_t
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_template_receipt_gas(
		&cGasUsed,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return 0, cErr.svmError()
	}

	return uint64(cGasUsed), nil
}

func cSvmSpawnApp(runtime Runtime, spawnApp []byte, creator Address, hostCtx []byte,
	gasMetering bool, gasLimit uint64) ([]byte, error) {
	cReceipt := cSvmByteArray{}
	cRuntime := runtime.p
	cSpawnApp := bytesCloneToSvmByteArray(spawnApp)
	cCreator := bytesCloneToSvmByteArray(creator[:])
	cHostCtx := bytesCloneToSvmByteArray(hostCtx)
	cGasMetering := C.bool(gasMetering)
	cGasLimit := C.uint64_t(gasLimit)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.SvmFree()
		cSpawnApp.Free()
		cCreator.Free()
		cHostCtx.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_spawn_app(
		&cReceipt,
		cRuntime,
		cSpawnApp,
		cCreator,
		cHostCtx,
		cGasMetering,
		cGasLimit,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cReceipt), nil
}

func cSvmAppReceiptState(receipt []byte) ([]byte, error) {
	cInitialState := cSvmByteArray{}
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cInitialState.SvmFree()
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_app_receipt_state(
		&cInitialState,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cInitialState), nil
}

func cSvmAppReceiptAddr(receipt []byte) (Address, error) {
	cAppAddr := cSvmByteArray{}
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cAppAddr.SvmFree()
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_app_receipt_addr(
		&cAppAddr,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return Address{}, cErr.svmError()
	}

	return svmByteArrayCloneToAddress(cAppAddr), nil
}

func cSvmAppReceiptGas(receipt []byte) (uint64, error) {
	var cGasUsed C.uint64_t
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_app_receipt_gas(
		&cGasUsed,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return 0, cErr.svmError()
	}

	return uint64(cGasUsed), nil
}

func cSvmEncodeSpawnApp(version int, templateAddr Address, ctorIndex uint16, ctorBuffer []byte, ctorArgs Values) ([]byte, error) {
	spawnApp := cSvmByteArray{}
	cVersion := C.uint(version)
	cTemplateAddr := bytesCloneToSvmByteArray(templateAddr[:])
	cCtorIndex := C.ushort(ctorIndex)
	cCtorBuffer := bytesCloneToSvmByteArray(ctorBuffer)
	cCtorArgs := bytesCloneToSvmByteArray(ctorArgs.Encode())
	cErr := cSvmByteArray{}

	defer func() {
		spawnApp.SvmFree()
		cTemplateAddr.Free()
		cCtorBuffer.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_encode_spawn_app(
		&spawnApp,
		cVersion,
		cTemplateAddr,
		cCtorIndex,
		cCtorBuffer,
		cCtorArgs,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(spawnApp), nil
}

func cSvmValidateApp(runtime Runtime, app []byte) error {
	cRuntime := runtime.p
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

func cSvmEncodeAppTx(
	version int,
	AppAddr Address,
	funcIndex uint16,
	funcBuffer []byte,
	funcArgs Values,
) ([]byte, error) {
	appTx := cSvmByteArray{}
	cVersion := C.uint(version)
	cTemplateAddr := bytesCloneToSvmByteArray(AppAddr[:])
	cFuncIndex := C.ushort(funcIndex)
	cFuncBuffer := bytesCloneToSvmByteArray(funcBuffer)
	cFuncArgs := bytesCloneToSvmByteArray(funcArgs.Encode())
	cErr := cSvmByteArray{}

	defer func() {
		appTx.SvmFree()
		cTemplateAddr.Free()
		cFuncBuffer.Free()
		cFuncArgs.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_encode_app_tx(
		&appTx,
		cVersion,
		cTemplateAddr,
		cFuncIndex,
		cFuncBuffer,
		cFuncArgs,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(appTx), nil
}

func cSvmValidateTx(runtime Runtime, appTx []byte) (Address, error) {
	cAppAddr := cSvmByteArray{}
	cRuntime := runtime.p
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

func cSvmExecApp(runtime Runtime, appTx []byte, appState []byte, hostCtx []byte, gasMetering bool,
	gasLimit uint64) ([]byte, error) {
	cReceipt := cSvmByteArray{}
	cRuntime := runtime.p
	cAppTx := bytesCloneToSvmByteArray(appTx)
	cAppState := bytesCloneToSvmByteArray(appState)
	cHostCtx := bytesCloneToSvmByteArray(hostCtx)
	cGasMetering := C.bool(gasMetering)
	cGasLimit := C.uint64_t(gasLimit)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.SvmFree()
		cAppTx.Free()
		cAppState.Free()
		cHostCtx.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_exec_app(
		&cReceipt,
		cRuntime,
		cAppTx,
		cAppState,
		cHostCtx,
		cGasMetering,
		cGasLimit,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cReceipt), nil
}

func cSvmExecReceiptState(receipt []byte) ([]byte, error) {
	cNewState := cSvmByteArray{}
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cNewState.SvmFree()
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_exec_receipt_state(
		&cNewState,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	return svmByteArrayCloneToBytes(cNewState), nil
}

func cSvmExecReceiptReturns(receipt []byte) (Values, error) {
	cReturns := cSvmByteArray{}
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cReturns.SvmFree()
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_exec_receipt_returns(
		&cReturns,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return nil, cErr.svmError()
	}

	var nativeReturns Values
	if err := (&nativeReturns).Decode(svmByteArrayCloneToBytes(cReturns)); err != nil {
		return nil, fmt.Errorf("failed to decode returns: %v", err)
	}

	return nativeReturns, nil
}

func cSvmExecReceiptGas(receipt []byte) (uint64, error) {
	var cGasUsed C.uint64_t
	cReceipt := bytesCloneToSvmByteArray(receipt)
	cErr := cSvmByteArray{}

	defer func() {
		cReceipt.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_exec_receipt_gas(
		&cGasUsed,
		cReceipt,
		&cErr,
	); res != cSvmSuccess {
		return 0, cErr.svmError()
	}

	return uint64(cGasUsed), nil
}

func cSvmInstanceContextHostGet(ctx unsafe.Pointer) unsafe.Pointer {
	return C.svm_instance_context_host_get(ctx)
}

func cSvmByteArrayDestroy(ba cSvmByteArray) {
	C.svm_byte_array_destroy(ba)
}

func cSvmRuntimeDestroy(runtime Runtime) {
	C.svm_runtime_destroy(runtime.p)
}

func cSvmImportsDestroy(imports Imports) {
	C.svm_imports_destroy(imports.p)
}

func cSvmMemKVDestroy(kv MemKVStore) {
	C.svm_memory_kv_destroy(kv.p)
}

func cFree(p unsafe.Pointer) {
	C.free(p)
}
