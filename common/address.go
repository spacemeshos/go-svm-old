package common

import (
	"encoding/hex"
)

const AddressSize = 20

type Address [AddressSize]byte

func (addr Address) String() string {
	return hex.EncodeToString(addr[:])
}

func BytesToAddress(b []byte) Address {
	var addr Address
	if len(b) <= AddressSize {
		copy(addr[:], b)
	} else {
		copy(addr[:], b[:AddressSize])
	}

	return addr
}
