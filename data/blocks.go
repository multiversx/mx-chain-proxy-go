package data

import "github.com/multiversx/mx-chain-core-go/data/api"

// BlocksApiResponse is a response holding(possibly) multiple block
type BlocksApiResponse struct {
	Data  BlocksApiResponsePayload `json:"data"`
	Error string                   `json:"error"`
	Code  ReturnCode               `json:"code"`
}

// BlocksApiResponsePayload wraps a block
type BlocksApiResponsePayload struct {
	Blocks []*api.Block `json:"blocks"`
}
