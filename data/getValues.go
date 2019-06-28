package data

// ResponseGetValues defines a wrapper over string containing returned data in hex format
type ResponseGetValues struct {
	HexData string `json:"data"`
}

// GetValuesRequest defines the request struct for getValues
type GetValuesRequest struct {
	Address  string   `json:"scAddress"`
	FuncName string   `json:"funcName"`
	Args     []string `json:"args"`
}
