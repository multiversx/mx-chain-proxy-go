package data

import "encoding/json"

// RequestBodyWeb3 defines the data structure for web3 request body
type RequestBodyWeb3 struct {
	JsonRpc  string          `json:"jsonrpc"`
	Id       int             `json:"id"`
	FuncName string          `json:"method"`
	Params   json.RawMessage `json:"params"`
}

// ResponseWeb3 defines the data structure for web3 response
type ResponseWeb3 struct {
	JsonRpc string
	Id      int
	Result  interface{}
}
