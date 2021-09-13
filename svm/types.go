package svm

import "unsafe"

type TxType uint8
type Amount uint64
type Address [20]byte
type TemplateAddr [20]byte
type TxId [32]byte
type State [32]byte
type Gas uint64
type Layer uint64
type Log []byte

// Runtime is a wrapper for svm_runtime
type Runtime struct {
	raw unsafe.Pointer
}

type Context struct {
	layer Layer
	txId  TxId
}

type AccountCounter struct {
	upper uint64
	lower uint64
}

const (
	DeployType TxType = 0
	SpawnType  TxType = 1
	CallType   TxType = 2
)

type Envelope struct {
	Type      TxType
	Principal Address
	Amount    Amount
	Nonce     AccountCounter
	GasLimit  uint64
	GasFee    uint64
}

type Account struct {
	Addr    Address
	Balance Amount
	Counter AccountCounter
}

type DeployReceipt struct {
	Success bool
	Error   RuntimeError
	Addr    TemplateAddr
	GasUsed Gas
	Logs    []Log
}

type SpawnReceipt struct {
	Success         bool
	Error           RuntimeError
	AccountAddr     Address
	InitState       State
	GasUsed         Gas
	Logs            []Log
	TouchedAccounts []Address
}

type CallReceipt struct {
	Success         bool
	Error           RuntimeError
	NewState        State
	ReturnData      []byte
	GasUsed         Gas
	Logs            []Log
	TouchedAccounts []Address
}
