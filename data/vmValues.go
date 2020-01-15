package data

import vmcommon "github.com/ElrondNetwork/elrond-vm-common"

// ResponseVmValue defines a wrapper over string containing returned data in hex format
type ResponseVmValue struct {
	Error string             `json:"error"`
	Data  *vmcommon.VMOutput `json:"data"`
}

// VmValueRequest defines the request struct for values available in a VM
type VmValueRequest struct {
	Address  string   `json:"scAddress"`
	FuncName string   `json:"funcName"`
	Args     []string `json:"args"`
}
