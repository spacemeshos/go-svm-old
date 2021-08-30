package svm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// TxType is a transaction type abstraction
type TxType uint8

var strInvalidReceipt = "invalid receipt"
var errInvalidReceipt = errors.New(strInvalidReceipt)

// Receipt is a common receipt part
type Receipt struct {
	TxType
	Version int
	Success bool
}

const commonReceiptLength = 1 + 2 + 1

// Decode fills Receipt from bytes array
func (r *Receipt) Decode(bs []byte) error {
	if len(bs) < commonReceiptLength {
		return errInvalidReceipt
	}
	r.TxType = TxType(bs[0])
	r.Version = int(binary.BigEndian.Uint16(bs[1:]))
	r.Success = bs[3] != 0
	return nil
}

// DeployReceipt is a receipt returned by Runtime.Deply endpoint
type DeployReceipt struct {
	Receipt
	Address Address
	GasUsed uint64
}

// Decode fills receipt from bytes array
func (dr *DeployReceipt) Decode(bs []byte) error {
	if err := dr.Receipt.Decode(bs); err != nil {
		return err
	}
	if len(bs) < commonReceiptLength+addressLength+8 {
		return errInvalidReceipt
	}
	copy(dr.Address[:], bs[commonReceiptLength:commonReceiptLength+addressLength])
	dr.GasUsed = binary.BigEndian.Uint64(bs[commonReceiptLength+addressLength:])
	return nil
}

// ReturnData is a wrapper to data returnde from Runtime.Call/Spawn endpoints
type ReturnData []byte

// Decode constructs array of reflect.Value with bytes returned from Runtime endpoint
func (rd ReturnData) Decode() ([]reflect.Value, error) {
	return DecodeReturnData(rd)
}

// ReceiptResult is a receipt part returned from Runtime.Call/Spawn endpoints
type ReceiptResult struct {
	State
	Return  ReturnData
	GasUsed uint64
	Logs    [][]byte
}

// Decode fills receipt from bytes array
func (rr *ReceiptResult) Decode(bs []byte) (err error) {
	if len(bs) < 32 {
		return fmt.Errorf(strInvalidReceipt + ": no state here")
	}
	copy(rr.State[:], bs[:32])
	bs = bs[32:]
	returnLen := int(binary.BigEndian.Uint16(bs))
	if len(bs) < 2+returnLen {
		return fmt.Errorf(strInvalidReceipt + ": no return data here")
	}
	rr.Return = bs[2 : returnLen+2]
	if err != nil {
		return
	}
	bs = bs[2+returnLen:]
	if len(bs) < 8 {
		return fmt.Errorf(strInvalidReceipt + ": no gasUsed here")
	}
	rr.GasUsed = binary.BigEndian.Uint64(bs)
	bs = bs[8:]
	if len(bs) > 0 {
		logsCount := int(bs[0])
		rr.Logs = make([][]byte, logsCount)
		for i := range rr.Logs {
			// get logs
			rr.Logs[i] = make([]byte, 0)
		}
	}
	return
}

// SpawnReceipt is an abstraction for receipt returned from Runtime.Spawn endpoint
type SpawnReceipt struct {
	Receipt
	Address Address
	ReceiptResult
}

// Decode fills receipt from bytes array
func (sr *SpawnReceipt) Decode(bs []byte) error {
	if err := sr.Receipt.Decode(bs); err != nil {
		return err
	}
	if sr.Success {
		if len(bs) < commonReceiptLength+addressLength+32 {
			return fmt.Errorf(strInvalidReceipt + ": no initial state here")
		}
		copy(sr.Address[:], bs[commonReceiptLength:commonReceiptLength+addressLength])
		return sr.ReceiptResult.Decode(bs[commonReceiptLength+addressLength:])
	}
	return nil
}

// CallReceipt is an abstraction for receipt returned from runtime.Call endpoint
type CallReceipt struct {
	Receipt
	ReceiptResult
}

// Decode fills receipt from bytes array
func (cr *CallReceipt) Decode(bs []byte) error {
	if err := cr.Receipt.Decode(bs); err != nil {
		return err
	}
	if cr.Success {
		return cr.ReceiptResult.Decode(bs[commonReceiptLength:])
	}
	return nil
}
