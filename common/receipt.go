package common

type ReceiptDeployTemplate struct {
	Success      bool
	Version      int
	TemplateAddr Address
	GasUsed      uint64
}

type ReceiptSpawnApp struct {
	Success    bool
	Version    int
	AppAddr    Address
	State      []byte
	Returndata []byte
	Logs       []string
	GasUsed    uint64
}

type ReceiptExecApp struct {
	Success    bool
	Version    int
	NewState   []byte
	Returndata []byte
	Logs       []string
	GasUsed    uint64
}
