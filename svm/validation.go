package svm

func ValidateTemplate(runtime Runtime, appTemplate []byte) error {
	return cSvmValidateTemplate(runtime, appTemplate)
}

func ValidateApp(runtime Runtime, app []byte) error {
	return cSvmValidateApp(runtime, app)
}

func ValidateAppTx(runtime Runtime, appTx []byte) (Address, error) {
	return cSvmValidateTx(runtime, appTx)
}
