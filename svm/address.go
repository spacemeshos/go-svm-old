package svm

const AddressLen = 20

type Address [AddressLen]byte

func bytesToAddress(b []byte) Address {
	var addr Address
	if len(b) <= AddressLen {
		copy(addr[:], b)
	} else {
		copy(addr[:], b[:AddressLen])
	}

	return addr
}

func svmByteArrayCloneToAddress(ba cSvmByteArray) Address {
	b := svmByteArrayCloneToBytes(ba)
	return bytesToAddress(b)
}
