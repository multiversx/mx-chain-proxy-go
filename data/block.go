package data

import (
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
	"github.com/multiversx/mx-chain-core-go/data/api"
)

// BlockApiResponse is a response holding a block
type BlockApiResponse struct {
	Data  BlockApiResponsePayload `json:"data"`
	Error string                  `json:"error"`
	Code  ReturnCode              `json:"code"`
}

// BlockApiResponsePayload wraps a block
type BlockApiResponsePayload struct {
	Block api.Block `json:"block"`
}

// HyperblockApiResponse is a response holding a hyperblock
type HyperblockApiResponse struct {
	Data  HyperblockApiResponsePayload `json:"data"`
	Error string                       `json:"error"`
	Code  ReturnCode                   `json:"code"`
}

// NewHyperblockApiResponse creates a HyperblockApiResponse
func NewHyperblockApiResponse(hyperblock api.Hyperblock) *HyperblockApiResponse {
	return &HyperblockApiResponse{
		Data: HyperblockApiResponsePayload{
			Hyperblock: hyperblock,
		},
		Code: ReturnCodeSuccess,
	}
}

// HyperblockApiResponsePayload wraps a hyperblock
type HyperblockApiResponsePayload struct {
	Hyperblock api.Hyperblock `json:"hyperblock"`
}

// InternalBlockApiResponse is a response holding an internal block
type InternalBlockApiResponse struct {
	Data  InternalBlockApiResponsePayload `json:"data"`
	Error string                          `json:"error"`
	Code  ReturnCode                      `json:"code"`
}

// InternalBlockApiResponsePayload wraps a internal generic block
type InternalBlockApiResponsePayload struct {
	Block interface{} `json:"block"`
}

// ValidatorsInfoApiResponse is a response holding validators info
type ValidatorsInfoApiResponse struct {
	Data  InternalStartOfEpochValidators `json:"data"`
	Error string                         `json:"error"`
	Code  ReturnCode                     `json:"code"`
}

// InternalBlockApiResponsePayload wraps a internal generic validators info
type InternalStartOfEpochValidators struct {
	ValidatorsInfo interface{} `json:"validators"`
}

// InternalMiniBlockApiResponse is a response holding an internal miniblock
type InternalMiniBlockApiResponse struct {
	Data  InternalMiniBlockApiResponsePayload `json:"data"`
	Error string                              `json:"error"`
	Code  ReturnCode                          `json:"code"`
}

// InternalMiniBlockApiResponsePayload wraps an internal miniblock
type InternalMiniBlockApiResponsePayload struct {
	MiniBlock interface{} `json:"miniblock"`
}

// AlteredAccountsApiResponse is a response holding a altered accounts
type AlteredAccountsApiResponse struct {
	Data  AlteredAccountsPayload `json:"data"`
	Error string                 `json:"error"`
	Code  ReturnCode             `json:"code"`
}

// AlteredAccountsPayload wraps altered accounts payload
type AlteredAccountsPayload struct {
	Accounts []*alteredAccount.AlteredAccount `json:"accounts"`
}
