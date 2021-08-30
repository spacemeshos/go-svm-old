package svm

import (
	"encoding/binary"
	"errors"
	"reflect"
)

// Call32u the uint32 abstraction for SVM calldata
type Call32u uint32

// Encode uint32 value as SVM calldata
func (c Call32u) Encode() []byte {
	bs := []byte{0b_0_111_0011, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(bs[1:], uint32(c))
	return bs
}

// DecodeReturnData as array of reflect.Value
func DecodeReturnData(bs []byte) (r []reflect.Value, err error) {
	if len(bs) == 0 {
		return []reflect.Value{}, nil
	}
	if bs[0]&0x0e == 0b0110 { // array
		arrLen := int(bs[0] >> 4)
		if bs[0]&1 == 1 {
			arrLen += 8
		}
		r = make([]reflect.Value, arrLen)
		bs = bs[1:]
		for i := range r {
			r[i], bs, err = decodeValue(bs)
			if err != nil {
				return
			}
		}
		return
	}
	r = []reflect.Value{{}}
	r[0], _, err = decodeValue(bs)
	return
}

func copyBeVarlen(bs []byte, ln int, ext []byte) {
	offs := len(ext) - ln
	for i := 0; i < ln; i++ {
		ext[offs+i] = bs[i+1] // skip layout marker
	}
}

func decodeValue(bs []byte) (reflect.Value, []byte, error) {
	switch bs[0] {
	case 0b_0_000_0000, 0b_0_001_0000: // False/True
		return reflect.ValueOf(bs[0] != 0), bs[1:], nil
	case 0b_0_010_0000, 0b_0_011_0000:
		return reflect.ValueOf(nil), bs[1:], nil
	case 0b_0_100_0000: // Address
		addr := Address{}
		copy(addr[:], bs[1:1+addressLength])
		return reflect.ValueOf(addr), bs[addressLength+1:], nil
	case 0b_0_000_0010: // signed int8
		return reflect.ValueOf(int8(bs[1])), bs[2:], nil
	case 0b_0_001_0010: // unsigned int8
		return reflect.ValueOf(bs[1]), bs[2:], nil
	case 0b_0_010_0010: // signed int16
		return reflect.ValueOf(int16(binary.BigEndian.Uint16(bs[1:]))), bs[3:], nil
	case 0b_0_011_0010: // unsigned int16
		return reflect.ValueOf(binary.BigEndian.Uint16(bs[1:])), bs[3:], nil
	case 0b_0_000_0011, 0b_0_001_0011, 0b_0_010_0011, 0b_0_011_0011: // signed int32
		b := []byte{0, 0, 0, 0}
		ln := int(bs[0]>>4) + 1
		copyBeVarlen(bs, ln, b)
		return reflect.ValueOf(int32(binary.BigEndian.Uint32(b))), bs[ln+1:], nil
	case 0b_0_100_0011, 0b_0_101_0011, 0b_0_110_0011, 0b_0_111_0011: // unsigned int32
		b := []byte{0, 0, 0, 0}
		ln := int((bs[0]&0b00110000)>>4) + 1
		copyBeVarlen(bs, ln, b)
		return reflect.ValueOf(binary.BigEndian.Uint32(b)), bs[ln+1:], nil
	case 0b_0_000_0100, 0b_0_001_0100, 0b_0_010_0100, 0b_0_011_0100, 0b_0_100_0100,
		0b_0_101_0100, 0b_0_110_0100, 0b_0_111_0100: // signed int64
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		ln := int(bs[0]>>4) + 1
		copyBeVarlen(bs, ln, b)
		return reflect.ValueOf(int32(binary.BigEndian.Uint64(b))), bs[ln+1:], nil
	case 0b_0_000_0101, 0b_0_001_0101, 0b_0_010_0101, 0b_0_011_0101, 0b_0_100_0101,
		0b_0_101_0101, 0b_0_110_0101, 0b_0_111_0101: // unsigned int64
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		ln := int(bs[0]>>4) + 1
		copyBeVarlen(bs, ln, b)
		return reflect.ValueOf(binary.BigEndian.Uint64(b)), bs[ln+1:], nil
	case 0b_0_000_0001, 0b_0_001_0001, 0b_0_010_0001, 0b_0_011_0001, 0b_0_100_0001,
		0b_0_101_0001, 0b_0_110_0001, 0b_0_111_0001: // Amount
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		ln := int(bs[0]>>4) + 1
		copyBeVarlen(bs, ln, b)
		return reflect.ValueOf(binary.BigEndian.Uint64(b)), bs[ln+1:], nil
	default:
		return reflect.Value{}, bs, errors.New("invalid value")
	}
}
