package svm

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Imports struct {
	p unsafe.Pointer
}

func (imports Imports) Free() {
	cSvmImportsDestroy(imports)
}

// ImportFunction represents a SVM runtime imported function.
type ImportFunction struct {
	// An implementation must be of type:
	// `func(ctx unsafe.Pointer, arguments ...interface{}) interface{}`.
	// It represents the real function implementation written in Go.
	implementation interface{}

	// The pointer to the cgo function implementation,
	// something like `C.foo`.
	cgoPointer unsafe.Pointer

	// The namespace of the imported function.
	namespace string

	// The function implementation signature as a WebAssembly signature.
	args ValueTypes

	// The function implementation signature as a WebAssembly signature.
	returns ValueTypes
}

type ImportsBuilder struct {
	// All imports.
	imports map[string]ImportFunction

	// Current namespace where to register the import.
	currentNamespace string
}

func NewImportsBuilder() ImportsBuilder {
	var imports = make(map[string]ImportFunction)
	var currentNamespace = "env"

	return ImportsBuilder{imports, currentNamespace}
}

// Namespace changes the current namespace of the next imported functions.
func (ib ImportsBuilder) Namespace(namespace string) ImportsBuilder {
	ib.currentNamespace = namespace
	return ib
}

func (ib ImportsBuilder) AppendFunction(name string, implementation interface{}, cgoPointer unsafe.Pointer) (ImportsBuilder, error) {
	args, returns, err := validateImport(name, implementation)
	if err != nil {
		return ImportsBuilder{}, err
	}

	namespace := ib.currentNamespace
	ib.imports[name] = ImportFunction{
		implementation,
		cgoPointer,
		namespace,
		args,
		returns,
	}

	return ib, nil
}

func (ib ImportsBuilder) Build() (Imports, error) {
	imports := Imports{}

	if res := cSvmImportsAlloc(&imports.p, uint(len(ib.imports))); res != cSvmSuccess {
		return Imports{}, fmt.Errorf("failed to allocate imports")
	}

	for importName, importFunction := range ib.imports {
		if err := cSvmImportFuncBuild(
			imports,
			importFunction.namespace,
			importName,
			importFunction.cgoPointer,
			importFunction.args,
			importFunction.returns,
		); err != nil {
			return Imports{}, fmt.Errorf("failed to build import `%v`: %v", importName, err)
		}
	}

	return imports, nil
}

func validateImport(name string, implementation interface{}) (args ValueTypes, returns ValueTypes, err error) {
	var importType = reflect.TypeOf(implementation)

	if importType.Kind() != reflect.Func {
		err = fmt.Errorf("imported function `%s` must be a function; given `%s`", name, importType.Kind())
		return
	}

	var inputArity = importType.NumIn()

	if inputArity < 1 {
		err = fmt.Errorf("imported function `%s` must at least have one argument (for the runtime context)", name)
		return
	}

	if importType.In(0).Kind() != reflect.UnsafePointer {
		err = fmt.Errorf("the runtime context of the `%s` imported function must be of kind `unsafe.Pointer`; given `%s`", name, importType.In(0).Kind())
		return
	}

	inputArity--

	var outputArity = importType.NumOut()
	args = make(ValueTypes, inputArity)
	returns = make(ValueTypes, outputArity)

	for i := 0; i < inputArity; i++ {
		var importInput = importType.In(i + 1)

		switch importInput.Kind() {
		case reflect.Int32:
			args[i] = TypeI32
		case reflect.Int64:
			args[i] = TypeI64
		default:
			err = fmt.Errorf("invalid input type for the `%s` imported function; given `%s`; only accept `int32` and `int64`", name, importInput.Kind())
			return
		}
	}

	if outputArity > 1 {
		err = fmt.Errorf("the `%s` imported function must have at most one output value", name)
		return
	} else if outputArity == 1 {
		switch importType.Out(0).Kind() {
		case reflect.Int32:
			returns[0] = TypeI32
		case reflect.Int64:
			returns[0] = TypeI64
		default:
			err = fmt.Errorf("invalid output type for the `%s` imported function; given `%s`; only accept `int32` and `int64`", name, importType.Out(0).Kind())
			return
		}
	}

	return
}
