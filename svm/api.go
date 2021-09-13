package svm

var initialized = false

func Init(inMemory bool, path string) {
	if initialized {
		panic("`Init` can be called only once")
	}

	panic("TODO")
}

func AssertInitialized() {
	if !initialized {
		panic("Forgot to call `Init`")
	}
}

func NewRuntime() (*Runtime, error) {
	panic("TODO")
}

// Releases SVM runtime
func (rt *Runtime) Destroy() {
	panic("TODO")
}

func (rt *Runtime) Close() error {
	panic("TODO")
}

func ValidateDeploy(rt *Runtime, msg []byte) (bool, ValidateError) {
	panic("TODO")
}

func Deploy(rt *Runtime, env Envelope, msg []byte, ctx Context) DeployReceipt {
	panic("TODO")
}

func ValidateSpawn(rt *Runtime, msg []byte) (bool, ValidateError) {
	panic("TODO")
}

func Spawn(rt *Runtime, env Envelope, msg []byte, ctx Context) SpawnReceipt {
	panic("TODO")
}

func ValidateCall(rt *Runtime, msg []byte) (bool, ValidateError) {
	panic("TODO")
}

func Call(rt *Runtime, env Envelope, msg []byte, ctx Context) CallReceipt {
	panic("TODO")
}

func (rt *Runtime) Rewind(layer Layer) {
	panic("TODO")
}

func (rt *Runtime) Commit() {
	panic("TODO")
}

func (rt *Runtime) GetAccount(addr Address) Account {
	panic("TODO")
}

func (rt *Runtime) IncreaseBalance(amount Amount) {
	panic("TODO")
}
