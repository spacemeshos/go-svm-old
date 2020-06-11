package svm

import "fmt"

type DeployTemplateResult struct {
	Receipt      []byte
	TemplateAddr Address
	GasUsed      uint64
}

func (r DeployTemplateResult) String() string {
	return fmt.Sprintf(
		"DeployTemplate result:\n"+
			"  Receipt: %x\n"+
			"  TemplateAddr: %x\n"+
			"  GasUsed: %v\n",
		r.Receipt, r.TemplateAddr, r.GasUsed)
}

type SpawnAppResult struct {
	Receipt      []byte
	InitialState []byte
	AppAddr      Address
	GasUsed      uint64
}

func (r SpawnAppResult) String() string {
	return fmt.Sprintf(
		"SpawnApp result:\n"+
			"  Receipt: %x\n"+
			"  InitialState: %x\n"+
			"  AppAddr: %x\n"+
			"  GasUsed: %v\n",
		r.Receipt, r.InitialState, r.AppAddr, r.GasUsed)
}

type ExecAppResult struct {
	Receipt  []byte
	NewState []byte
	Returns  Values
	GasUsed  uint64
}

func (r ExecAppResult) String() string {
	return fmt.Sprintf(
		"ExecApp result:\n"+
			"  Receipt: %x\n"+
			"  NewState: %x\n"+
			"  Returns: %v\n"+
			"  GasUsed: %v\n",
		r.Receipt, r.NewState, r.Returns, r.GasUsed)
}

func DeployTemplate(runtime Runtime, appTemplate []byte, author Address, hostCtx []byte, gasMetering bool, gasLimit uint64) (*DeployTemplateResult, error) {
	receipt, err := cSvmDeployTemplate(runtime, appTemplate, author, hostCtx, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	addr, err := cSvmTemplateReceiptAddr(receipt)
	if err != nil {
		return nil, err
	}

	gasUsed, err := cSvmTemplateReceiptGas(receipt)
	if err != nil {
		return nil, err
	}

	return &DeployTemplateResult{
		Receipt:      receipt,
		TemplateAddr: addr,
		GasUsed:      gasUsed,
	}, nil
}

func SpawnApp(runtime Runtime, spawnAppData []byte, creator Address, hostCtx []byte,
	gasMetering bool, gasLimit uint64) (*SpawnAppResult, error) {

	receipt, err := cSvmSpawnApp(runtime, spawnAppData, creator, hostCtx, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	initialState, err := cSvmAppReceiptState(receipt)
	if err != nil {
		return nil, err
	}

	addr, err := cSvmAppReceiptAddr(receipt)
	if err != nil {
		return nil, err
	}

	gasUsed, err := cSvmAppReceiptGas(receipt)
	if err != nil {
		return nil, err
	}

	return &SpawnAppResult{
		Receipt:      receipt,
		InitialState: initialState,
		AppAddr:      addr,
		GasUsed:      gasUsed,
	}, nil
}

func ExecApp(runtime Runtime, appTx, appState, hostCtx []byte, gasMetering bool,
	gasLimit uint64) (*ExecAppResult, error) {

	receipt, err := cSvmExecApp(runtime, appTx, appState, hostCtx, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	newState, err := cSvmExecReceiptState(receipt)
	if err != nil {
		return nil, err
	}

	returns, err := cSvmExecReceiptReturns(receipt)
	if err != nil {
		return nil, err
	}

	gasUsed, err := cSvmExecReceiptGas(receipt)
	if err != nil {
		return nil, err
	}

	return &ExecAppResult{
		Receipt:  receipt,
		NewState: newState,
		Returns:  returns,
		GasUsed:  gasUsed,
	}, nil
}
