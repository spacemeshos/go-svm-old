package svm

import "go-svm/common"

func svmByteArrayCloneToAddress(ba cSvmByteArray) Address {
	b := svmByteArrayCloneToBytes(ba)
	return common.BytesToAddress(b)
}
