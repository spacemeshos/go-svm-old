package svm

// #include "./svm.h"
//
// extern svm_byte_array* svm_trampoline(svm_env_t *env, svm_byte_array *args, svm_byte_array *results);
//
import "C"
import (
	"fmt"
	"unsafe"
)

// Imports is used to register and coordinate the invocation of functions
// written in Go, directly from the SVM-managed WebAssembly modules.
type Imports struct {
	// _inner is a pointer to an SVM-managed heap allocation.
	_inner unsafe.Pointer

	// envs holds host import functions environment objects.
	// `svm_trampoline` will get the respective environment object raw pointer directly from SVM.
	// tracking it here is needed merely so that it won't get GC-ed.
	envs []*functionEnvironment
}

func (imports Imports) Free() {
	cSvmImportsDestroy(imports)
}

// ImportFunction represents an SVM-runtime imported function.
type ImportFunction struct {
	// f represents the actual function implementation, written in Go.
	f hostFunction

	// params is the WebAssembly signature of the function implementation params.
	params ValueTypes

	// returns is the WebAssembly signature of the function implementation returns.
	returns ValueTypes

	// namespace is the imported function WebAssembly namespace.
	namespace string
}

type ImportsBuilder struct {
	imports          map[string]ImportFunction
	currentNamespace string
}

func NewImportsBuilder() ImportsBuilder {
	var imports = make(map[string]ImportFunction)
	var currentNamespace = "host"

	return ImportsBuilder{imports, currentNamespace}
}

// Namespace changes the current namespace of the next imported functions.
func (ib ImportsBuilder) Namespace(namespace string) ImportsBuilder {
	ib.currentNamespace = namespace
	return ib
}

func (ib ImportsBuilder) RegisterFunction(name string, params ValueTypes, returns ValueTypes, f hostFunction) ImportsBuilder {
	ib.imports[name] = ImportFunction{
		f,
		params,
		returns,
		ib.currentNamespace,
	}

	return ib
}

func (ib ImportsBuilder) Build() (*Imports, error) {
	imports := Imports{}
	imports.envs = make([]*functionEnvironment, 0)

	if res := cSvmImportsAlloc(&imports._inner, uint(len(ib.imports))); res != cSvmSuccess {
		return nil, fmt.Errorf("failed to allocate imports")
	}

	for imprtName, imprt := range ib.imports {
		// hostEnv is used to define the minimal context of the import function,
		// to be used by `svm_trampoline` for its invocation.
		hostEnv := functionEnvironment{
			hostFunctionStoreIndex: hostFunctionStore.add(imprt.f),
		}
		imports.envs = append(imports.envs, &hostEnv)

		if err := cSvmImportFuncNew(
			imports,
			imprt.namespace,
			imprtName,
			unsafe.Pointer(&hostEnv),
			imprt.params,
			imprt.returns,
		); err != nil {
			return nil, fmt.Errorf("failed to build import function `%v`: %v", imprtName, err)
		}
	}

	return &imports, nil
}

func cSvmImportFuncNew(
	imports Imports,
	namespace string,
	name string,
	hostEnv unsafe.Pointer,
	params ValueTypes,
	returns ValueTypes,
) error {
	cImports := imports._inner
	cNamespace := bytesCloneToSvmByteArray([]byte(namespace))
	cImportName := bytesCloneToSvmByteArray([]byte(name))
	cParams := bytesCloneToSvmByteArray(params.Encode())
	cReturns := bytesCloneToSvmByteArray(returns.Encode())
	cErr := cSvmByteArray{}

	defer func() {
		cNamespace.Free()
		cImportName.Free()
		cParams.Free()
		cReturns.Free()
		cErr.SvmFree()
	}()

	if res := C.svm_import_func_new(
		cImports,
		cNamespace,
		cImportName,
		(C.svm_func_callback_t)(C.svm_trampoline),
		hostEnv,
		cParams,
		cReturns,
		&cErr,
	); res != cSvmSuccess {
		return cErr.svmError()
	}

	return nil
}
