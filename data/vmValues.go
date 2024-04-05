package data

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
)

// VmValuesResponseData follows the format of the data field in an API response for a VM values query
type VmValuesResponseData struct {
	Data      *vm.VMOutputApi `json:"data"`
	BlockInfo BlockInfo       `json:"blockInfo"`
}

// ResponseVmValue defines a wrapper over string containing returned data in hex format
type ResponseVmValue struct {
	Data  VmValuesResponseData `json:"data"`
	Error string               `json:"error"`
	Code  string               `json:"code"`
}

// VmValueRequest defines the request struct for values available in a VM
type VmValueRequest struct {
	Address        string   `json:"scAddress"`
	FuncName       string   `json:"funcName"`
	CallerAddr     string   `json:"caller"`
	CallValue      string   `json:"value"`
	SameScState    bool     `json:"sameScState"`
	ShouldBeSynced bool     `json:"shouldBeSynced"`
	Args           []string `json:"args"`
}

// SCQuery represents a prepared query for executing a function of the smart contract
type SCQuery struct {
	ScAddress      string
	FuncName       string
	CallerAddr     string
	CallValue      string
	SameScState    bool `json:"sameScState"`
	ShouldBeSynced bool `json:"shouldBeSynced"`
	Arguments      [][]byte
	BlockNonce     core.OptionalUint64
	BlockHash      []byte
}
