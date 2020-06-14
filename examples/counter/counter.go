package main

import (
	"fmt"
	"go-svm/svm"
	"io/ioutil"
	"unsafe"
)

// Declare `inc` and `get` C function signatures in the cgo preamble.
// Doing so is required for SVM to be able to invoke their Go implementation.

// #include <stdlib.h>
//
// extern void inc(void *ctx, int value);
// extern int get(void *ctx);
import "C"

// Define `inc` and `get` Go implementation.
// The first argument is the runtime context, and must be included by all import functions.
// Notice the `//export` comment which is the way cgo uses to map Go code to C code.

//export inc
func inc(ctx unsafe.Pointer, value int32) {
	// SVM import function can access a closure variable,
	closure.value += value

	// or the runtime's provided host object pointer.
	host := (*counter)(svm.InstanceContextHostGet(ctx))
	host.value += value
}

//export get
func get(ctx unsafe.Pointer) int32 {
	host := (*counter)(svm.InstanceContextHostGet(ctx))

	// Host and closure variable should be synced.
	if host.value != closure.value {
		panic("Mayday")
	}

	return host.value
}

type counter struct {
	value int32
}

var closure counter

var host counter

func main() {
	// 1) Initialize runtime.
	ib := svm.NewImportsBuilder()
	ib, err := ib.AppendFunction("inc", inc, C.inc)
	noError(err)
	ib, err = ib.AppendFunction("get", get, C.get)
	noError(err)
	imports, err := ib.Build()
	noError(err)
	defer imports.Free()

	kv, err := svm.NewMemoryKVStore()
	noError(err)

	rawKV, err := svm.NewMemoryRawKVStore()
	noError(err)

	host := unsafe.Pointer(&host)

	runtime, err := svm.NewRuntimeBuilder().
		WithImports(imports).
		WithMemoryKV(kv).
		WithMemoryRawKV(rawKV).
		WithHost(host).
		Build()
	noError(err)
	defer runtime.Free()
	fmt.Printf("1) Runtime: %v\n\n", runtime)

	version := 0
	hostCtx := svm.NewHostCtx().Encode()
	gasMetering := false
	gasLimit := uint64(0)

	// 2) Deploy Template.
	// TODO: add on-the-fly wat2wasm translation
	code, err := ioutil.ReadFile("counter_template.wasm")
	noError(err)
	name := "name"
	pageCount := 1
	author := svm.Address{}
	deployTemplateResult, err := deployTemplate(
		runtime,
		code,
		version,
		name,
		pageCount,
		author,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	fmt.Printf("2) %v\n", deployTemplateResult)

	// 3) Spawn App.
	creator := svm.Address{}
	spawnAppResult, err := spawnApp(
		runtime,
		version,
		deployTemplateResult.TemplateAddr,
		uint16(0),
		[]byte(nil),
		svm.Values{svm.I32(5)},
		creator,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	fmt.Printf("3) %v\n", spawnAppResult)

	// 4) Exec App
	// 4.0) Storage value increment.
	execAppResult, err := execApp(
		runtime,
		version,
		spawnAppResult.AppAddr,
		uint16(0),
		[]byte(nil),
		svm.Values{svm.I32(5)},
		spawnAppResult.InitialState,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	fmt.Printf("4.0) %v\n", execAppResult)

	// 4.1) Storage value get.
	execAppResult, err = execApp(
		runtime,
		version,
		spawnAppResult.AppAddr,
		uint16(1),
		[]byte(nil),
		svm.Values(nil),
		execAppResult.NewState,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	fmt.Printf("4.1) %v\n", execAppResult)

	// 4.2) Host import function value increment.
	execAppResult, err = execApp(
		runtime,
		version,
		spawnAppResult.AppAddr,
		uint16(2),
		[]byte(nil),
		svm.Values{svm.I32(25)},
		execAppResult.NewState,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	fmt.Printf("4.2) %v\n", execAppResult)

	// 4.3) Host import function value get.
	execAppResult, err = execApp(
		runtime,
		version,
		spawnAppResult.AppAddr,
		uint16(3),
		[]byte(nil),
		svm.Values(nil),
		execAppResult.NewState,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	fmt.Printf("4.3) %v\n", execAppResult)
}

func deployTemplate(
	runtime svm.Runtime,
	code []byte,
	version int,
	name string,
	pageCount int,
	author svm.Address,
	hostCtx []byte,
	gasMetering bool,
	gasLimit uint64,
) (*svm.DeployTemplateResult, error) {
	appTemplate, err := svm.EncodeAppTemplate(
		version,
		name,
		pageCount,
		code,
	)
	if err != nil {
		return nil, err
	}

	res, err := svm.DeployTemplate(
		runtime,
		appTemplate,
		author,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func spawnApp(
	runtime svm.Runtime,
	version int,
	templateAddr svm.Address,
	ctorIndex uint16,
	ctorBuffer []byte,
	ctorArgs svm.Values,
	creator svm.Address,
	hostCtx []byte,
	gasMetering bool,
	gasLimit uint64,
) (*svm.SpawnAppResult, error) {
	spawnApp, err := svm.EncodeSpawnApp(
		version,
		templateAddr,
		ctorIndex,
		ctorBuffer,
		ctorArgs,
	)
	if err != nil {
		return nil, err
	}

	res, err := svm.SpawnApp(
		runtime,
		spawnApp,
		creator,
		hostCtx,
		gasMetering,
		gasLimit,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func execApp(
	runtime svm.Runtime,
	version int,
	appAddr svm.Address,
	funcIndex uint16,
	funcBuffer []byte,
	funcArgs svm.Values,
	appState []byte,
	hostCtx []byte,
	gasMetering bool,
	gasLimit uint64,
) (*svm.ExecAppResult, error) {
	appTx, err := svm.EncodeAppTx(
		version,
		appAddr,
		funcIndex,
		funcBuffer,
		funcArgs,
	)
	if err != nil {
		return nil, err
	}

	if _, err = svm.ValidateAppTx(runtime, appTx); err != nil {
		return nil, err
	}

	execAppResult, err := svm.ExecApp(runtime, appTx, appState, hostCtx, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	return execAppResult, nil
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}
