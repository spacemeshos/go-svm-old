package svm

func EncodeAppTemplate(version int, name string, pageCount int, code []byte) ([]byte, error) {
	return cSvmEncodeAppTemplate(version, name, pageCount, code)
}

func EncodeSpawnApp(version int, templateAddr Address, ctorIndex uint16, ctorBuffer []byte, ctorArgs Values) ([]byte, error) {
	return cSvmEncodeSpawnApp(version, templateAddr, ctorIndex, ctorBuffer, ctorArgs)
}

func EncodeAppTx(version int, appAddr Address, funcIndex uint16, funcBuffer []byte, funcArgs Values) ([]byte, error) {
	return cSvmEncodeAppTx(version, appAddr, funcIndex, funcBuffer, funcArgs)
}
