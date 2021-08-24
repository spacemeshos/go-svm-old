package svm

import (
	"go-svm/codec"
	"go-svm/common"
)

const (
	KVValueSize = 32
	StateSize   = 32
	AddressSize = common.AddressSize
)

type (
	DeployTemplateReceipt = common.ReceiptDeployTemplate
	SpawnAppReceipt       = common.ReceiptSpawnApp
	ExecAppReceipt        = common.ReceiptExecApp

	Address = common.Address
)

var (
	BytesToAddress = common.BytesToAddress
)

func DeployTemplate(runtime Runtime, appTemplate []byte, author Address, gasMetering bool, gasLimit uint64) (*DeployTemplateReceipt, error) {
	rawReceipt, err := cSvmDeployTemplate(runtime, appTemplate, author, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	return codec.DecodeReceiptDeployTemplate(rawReceipt)
}

func SpawnApp(runtime Runtime, spawnAppData []byte, creator Address, gasMetering bool, gasLimit uint64) (*SpawnAppReceipt, error) {
	rawReceipt, err := cSvmSpawnApp(runtime, spawnAppData, creator, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	return codec.DecodeReceiptSpawnApp(rawReceipt)
}

func ExecApp(runtime Runtime, tx, appState []byte, gasMetering bool, gasLimit uint64) (*ExecAppReceipt, error) {
	rawReceipt, err := cSvmExecApp(runtime, tx, appState, gasMetering, gasLimit)
	if err != nil {
		return nil, err
	}

	return codec.DecodeReceiptExecApp(rawReceipt)
}

func ValidateTemplate(runtime Runtime, appTemplate []byte) error {
	return cSvmValidateTemplate(runtime, appTemplate)
}

func ValidateApp(runtime Runtime, app []byte) error {
	return cSvmValidateApp(runtime, app)
}

func ValidateAppTx(runtime Runtime, appTx []byte) (Address, error) {
	return cSvmValidateTx(runtime, appTx)
}
