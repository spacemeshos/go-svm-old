package svm

import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ValueType represents the `Value` type.
type ValueType uint8

const (
	// TypeI32 represents the SVM `i32` type.
	TypeI32 ValueType = 0

	// TypeI64 represents the SVM `i64` type.
	TypeI64 ValueType = 1
)

type ValueTypes []ValueType

// Encode encodes ValueTypes according to the following format:
//
// +---------------------------------+
// |  type #1  |  . . . |  type #N   |
// |  (1 byte) |        |  (1 byte)  |
// +-----------+--------+------------+
//
// `type` can be either 0 (TypeI32) or 1 (TypeI64).
// Note: the number of `type` values equals the number of bytes (one byte per-type).
func (v ValueTypes) Encode() []byte {
	b := make([]byte, len(v))
	for i, vt := range v {
		b[i] = byte(vt)
	}
	return b
}

// Value represents a SVM value of a particular type.
type Value struct {
	// The SVM value (as bits).
	value uint64

	// The SVM value type.
	ty ValueType
}

// I32 constructs a SVM value of type `i32`.
func I32(value int32) Value {
	return Value{
		value: uint64(value),
		ty:    TypeI32,
	}
}

// I64 constructs a SVM value of type `i64`.
func I64(value int64) Value {
	return Value{
		value: uint64(value),
		ty:    TypeI64,
	}
}

// GetType gets the type of the SVM value.
func (v Value) Type() ValueType {
	return v.ty
}

// ToI32 reads the SVM value bits as an `int32`.
// The SVM value type is ignored.
func (v Value) ToI32() int32 {
	return int32(v.value)
}

// ToI64 reads the SVM value bits as an `int64`.
// The SVM value type is ignored.
func (v Value) ToI64() int64 {
	return int64(v.value)
}

// String helps Value to implement the Stringer interface.
func (v Value) String() string {
	switch v.ty {
	case TypeI32:
		return fmt.Sprintf("i32 %d", v.ToI32())
	case TypeI64:
		return fmt.Sprintf("i64 %d", v.ToI64())
	default:
		return ""
	}
}

// Encode encodes Value according to the following format:
//
// +--------------------------------------+
// | type (1 byte) | value (4 or 8 bytes) |
// +---------------+----------------------+
//
// `type` can be either 0 (TypeI32) or 1 (TypeI64).
// `value` byte order is Big-Endian.
func (v Value) Encode() []byte {
	switch v.ty {
	case TypeI32:
		b := make([]byte, 1+4)
		b[0] = byte(v.ty)
		binary.BigEndian.PutUint32(b[1:], uint32(v.ToI32()))
		return b
	case TypeI64:
		b := make([]byte, 1+8)
		b[0] = byte(v.ty)
		binary.BigEndian.PutUint64(b[1:], uint64(v.ToI64()))
		return b
	default:
		return nil
	}
}

type Values []Value

// Encode encodes Values according to the following format:
//
/// +------------------------------------------------------+
/// | #values  | value #1       |  . . .  | value #N       |
/// | (1 byte) | (5 or 9 bytes) |         | (5 or 9 bytes) |
/// +----------+----------------+---------+----------------+
//
// `value` encoding is defined separately.
func (values Values) Encode() []byte {
	buf := &bytes.Buffer{}

	numValues := byte(len(values))
	buf.Write([]byte{numValues})

	for _, v := range values {
		buf.Write(v.Encode())
	}

	return buf.Bytes()
}

// Decode decodes []byte slice according to the encoding format
// defined in the `Encode` method.
// If completed successfully, the result is assigned to the
// method pointer receiver value, hence the previous value is overridden.
// This method is intended to be called on a zero-value instance.
func (values *Values) Decode(data []byte) error {
	if len(data) == 0 {
		return errors.New("invalid input: empty data")
	}

	buf := bytes.NewBuffer(data)

	numValues, err := buf.ReadByte()
	if err != nil {
		return err
	}

	decodeValues := make(Values, numValues)

	for i := range decodeValues {
		ty, err := buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("failed to decode value #%v: bytes are missing", i)
			}
			return err
		}

		v := &decodeValues[i]

		switch ValueType(ty) {
		case TypeI32:
			next := buf.Next(4)
			if len(next) < 4 {
				return fmt.Errorf("failed to decode value #%v: "+
					"bytes are missing; expected: 4, given: %v", i, len(next))
			}
			v.ty = TypeI32
			v.value = uint64(binary.BigEndian.Uint32(next))
		case TypeI64:
			next := buf.Next(8)
			if len(next) < 8 {
				return fmt.Errorf("failed to decode value #%v: "+
					"bytes are missing; expected: 8, given: %v", i, len(next))
			}
			v.ty = TypeI64
			v.value = binary.BigEndian.Uint64(next)
		default:
			return fmt.Errorf("invalid type; expected: %v or %v, given: %v",
				TypeI32, TypeI64, ty)
		}
	}

	if buf.Len() > 0 {
		return fmt.Errorf("too many bytes; num expected: %v, num given: %v",
			len(data)-buf.Len(), len(data))
	}

	// Once completed successfully, override the method pointer receiver value.
	*values = decodeValues

	return nil
}
