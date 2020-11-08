package svm

func EncodeAppTemplate(version int, name string, code []byte, dataLayout DataLayout) ([]byte, error) {
	return cSvmEncodeAppTemplate(version, name, code, dataLayout)
}

func EncodeSpawnApp(version int, templateAddr Address, ctorIndex uint16, ctorBuffer []byte, ctorArgs Values) ([]byte, error) {
	return cSvmEncodeSpawnApp(version, templateAddr, ctorIndex, ctorBuffer, ctorArgs)
}

func EncodeAppTx(version int, appAddr Address, funcIndex uint16, funcBuffer []byte, funcArgs Values) ([]byte, error) {
	return cSvmEncodeAppTx(version, appAddr, funcIndex, funcBuffer, funcArgs)
}
