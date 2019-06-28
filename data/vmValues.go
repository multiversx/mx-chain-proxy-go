package data

// ResponseVmValue defines a wrapper over string containing returned data in hex format
type ResponseVmValue struct {
	HexData string `json:"data"`
}

// VmValueRequest defines the request struct for values available in a VM
type VmValueRequest struct {
	Address  string   `json:"scAddress"`
	FuncName string   `json:"funcName"`
	Args     []string `json:"args"`
}
