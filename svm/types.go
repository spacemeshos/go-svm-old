package svm

import "unsafe"

// Declaring types aliases used throughout the project.
type TxType uint8
type Amount uint64
type Address [20]byte
type TemplateAddr [20]byte
type TxId [32]byte
type State [32]byte
type Gas uint64
type Layer uint64
type Log []byte

// `Runtime` wraps the raw-Runtime returned by SVM C-API
type Runtime struct {
	raw unsafe.Pointer
}

// Holds the currently executed `Node Context`.
// Addionally, contains data implied/computed from the `input` transaction.
type Context struct {
	Layer Layer
	TxId  TxId
}

// Encapsulates an `Account Counter`. (Since `Golang` has no `unit128` primitive out-of-the-box).
//
// Used for implementing the `Nonce Scheme` implemented within the `Template` associated with the `Account`
type AccountCounter struct {
	Upper uint64
	Lower uint64
}

// A `Transaction Type` enum
const (
	DeployType TxType = 0
	SpawnType  TxType = 1
	CallType   TxType = 2
)

// Holds the `Envelope` of a transaction.
//
// In other words, holds fields which are part of any transaction regardless of its type (i.e `Deploy/Spawn/Call`).
type Envelope struct {
	Type      TxType
	Principal Address
	Amount    Amount
	Nonce     AccountCounter
	GasLimit  uint64
	GasFee    uint64
}

// Holds an `Account` basic information.
type Account struct {
	Addr    Address
	Balance Amount
	Counter AccountCounter
}

// Holds the data returned after executing a `Deploy` transaction.
type DeployReceipt struct {
	Success      bool
	Error        RuntimeError
	TemplateAddr TemplateAddr
	GasUsed      Gas
	Logs         []Log
}

// Holds the data returned after executing a `Spawn` transaction.
type SpawnReceipt struct {
	Success         bool
	Error           RuntimeError
	AccountAddr     Address
	InitState       State
	GasUsed         Gas
	Logs            []Log
	TouchedAccounts []Address
}

// Holds the data returned after executing a `Call` transaction.
type CallReceipt struct {
	Success         bool
	Error           RuntimeError
	NewState        State
	ReturnData      []byte
	GasUsed         Gas
	Logs            []Log
	TouchedAccounts []Address
}
