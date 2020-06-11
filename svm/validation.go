package svm

func ValidateAppTx(runtime Runtime, appTx []byte) (Address, error) {
	return cSvmValidateTx(runtime, appTx)
}
