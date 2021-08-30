package svm

const AddressLength = 20

type Address [AddressLength]byte

func StringAddress(str string) Address {
	a := Address{}
	copy(a[:], str)
	return a
}
