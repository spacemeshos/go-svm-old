package svm

const addressLength = 20

// Address is the address abstraction
type Address [addressLength]byte

// StringAddress creates address from string
func StringAddress(str string) Address {
	a := Address{}
	copy(a[:], str)
	return a
}
