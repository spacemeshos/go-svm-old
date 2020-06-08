package svm

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes(t *testing.T) {
	require.Equal(t, 0, int(TypeI32))
	require.Equal(t, 1, int(TypeI64))
}

func TestValueTypes_Encode(t *testing.T) {
	req := require.New(t)

	data := ValueTypes{TypeI32, TypeI64, TypeI32, TypeI64}.Encode()
	req.Equal(4, len(data))
	req.Equal(4, cap(data))
	req.Equal(uint8(TypeI32), data[0])
	req.Equal(uint8(TypeI64), data[1])
	req.Equal(uint8(TypeI32), data[2])
	req.Equal(uint8(TypeI64), data[3])
}

func TestValueTypes_Encode_Empty(t *testing.T) {
	req := require.New(t)

	data := ValueTypes{}.Encode()
	req.Equal(0, len(data))
	req.Equal(0, cap(data))

	data = ValueTypes(nil).Encode()
	req.Equal(0, len(data))
	req.Equal(0, cap(data))
}

func TestValues_Decode_Empty(t *testing.T) {
	req := require.New(t)
	v := Values{}

	err := v.Decode(nil)
	req.EqualError(err, "invalid input: empty data")

	err = v.Decode([]byte{})
	req.EqualError(err, "invalid input: empty data")
}

func TestValues_Encode_Decode_Empty(t *testing.T) {
	req := require.New(t)
	v := Values{}

	err := v.Decode(Values(nil).Encode())
	req.NoError(err)
	req.Equal(Values{}, v)

	err = v.Decode(Values{}.Encode())
	req.NoError(err)
	req.Equal(Values{}, v)
}

func TestValues_Encode_Decode(t *testing.T) {
	req := require.New(t)
	v := Values{}

	vBase := Values{I32(10), I64(20)}
	err := v.Decode(vBase.Encode())
	req.NoError(err)
	req.Equal(vBase, v)
}

func TestValues_Decode_Errors(t *testing.T) {
	req := require.New(t)
	v := Values{}

	v = Values{}
	vBase := Values{I32(10), I64(20)}
	vBaseData := vBase.Encode()
	err := v.Decode(append(vBaseData, byte(0)))
	req.EqualError(err, "too many bytes; num expected: 15, num given: 16")
	req.Equal(Values{}, v)

	vBase = Values{I32(10), I64(20)}
	vBaseData = vBase.Encode()
	err = v.Decode(vBaseData[:len(vBaseData)-1])
	req.EqualError(err, "failed to decode value #1: bytes are missing; expected: 8, given: 7")
	req.Equal(Values{}, v)
	err = v.Decode(vBaseData[:len(vBaseData)-9])
	req.EqualError(err, "failed to decode value #1: bytes are missing")
	req.Equal(Values{}, v)
	err = v.Decode(vBaseData[:len(vBaseData)-10])
	req.EqualError(err, "failed to decode value #0: bytes are missing; expected: 4, given: 3")
	req.Equal(Values{}, v)
	err = v.Decode(vBaseData[:len(vBaseData)-14])
	req.EqualError(err, "failed to decode value #0: bytes are missing")
	req.Equal(Values{}, v)
}
