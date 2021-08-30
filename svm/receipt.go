package svm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type TxType uint8

var strInvalidReceipt = "invalid receipt"
var errInvalidReceipt = errors.New(strInvalidReceipt)

type Receipt struct {
	TxType
	Version int
	Success bool
}

const CommonReceiptLength = 1 + 2 + 1

func (r *Receipt) Decode(bs []byte) error {
	if len(bs) < CommonReceiptLength {
		return errInvalidReceipt
	}
	r.TxType = TxType(bs[0])
	r.Version = int(binary.BigEndian.Uint16(bs[1:]))
	r.Success = bs[3] != 0
	return nil
}

type DeployReceipt struct {
	Receipt
	Address Address
	GasUsed uint64
}

func (dr *DeployReceipt) Decode(bs []byte) error {
	if err := dr.Receipt.Decode(bs); err != nil {
		return err
	}
	if len(bs) < CommonReceiptLength+AddressLength+8 {
		return errInvalidReceipt
	}
	copy(dr.Address[:], bs[CommonReceiptLength:CommonReceiptLength+AddressLength])
	dr.GasUsed = binary.BigEndian.Uint64(bs[CommonReceiptLength+AddressLength:])
	return nil
}

type ReturnData []reflect.Value

func (rd *ReturnData) Decode(bs []byte) {
	for len(bs) > 0 {
		switch bs[0] {
		case 0b00010110:
		}
	}
}

type ReceiptResult struct {
	State
	Return  ReturnData
	GasUsed uint64
	Logs    [][]byte
}

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
	rr.Return, err = DecodeReturnData(bs[2 : returnLen+2])
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

type SpawnReceipt struct {
	Receipt
	Address Address
	ReceiptResult
}

func (sr *SpawnReceipt) Decode(bs []byte) error {
	if err := sr.Receipt.Decode(bs); err != nil {
		return err
	}
	if sr.Success {
		if len(bs) < CommonReceiptLength+AddressLength+32 {
			return fmt.Errorf(strInvalidReceipt + ": no initial state here")
		}
		copy(sr.Address[:], bs[CommonReceiptLength:CommonReceiptLength+AddressLength])
		return sr.ReceiptResult.Decode(bs[CommonReceiptLength+AddressLength:])
	}
	return nil
}

type CallReceipt struct {
	Receipt
	ReceiptResult
}

func (cr *CallReceipt) Decode(bs []byte) error {
	if err := cr.Receipt.Decode(bs); err != nil {
		return err
	}
	if cr.Success {
		return cr.ReceiptResult.Decode(bs[CommonReceiptLength:])
	}
	return nil
}
