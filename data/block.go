package data

import (
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// AtlasBlock is a block, as required by BlockAtlas
// Will be removed when using the "hyperblock" route in BlockAtlas as well.
type AtlasBlock struct {
	Nonce        uint64                `form:"nonce" json:"nonce"`
	Hash         string                `form:"hash" json:"hash"`
	Transactions []DatabaseTransaction `form:"transactions" json:"transactions"`
}

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
func NewHyperblockApiResponse(hyperblock Hyperblock) *HyperblockApiResponse {
	return &HyperblockApiResponse{
		Data: HyperblockApiResponsePayload{
			Hyperblock: hyperblock,
		},
		Code: ReturnCodeSuccess,
	}
}

// HyperblockApiResponsePayload wraps a hyperblock
type HyperblockApiResponsePayload struct {
	Hyperblock Hyperblock `json:"hyperblock"`
}

// Hyperblock contains all fully executed (both in source and in destination shards) transactions notarized in a given metablock
type Hyperblock struct {
	Hash                   string                              `json:"hash"`
	PrevBlockHash          string                              `json:"prevBlockHash"`
	StateRootHash          string                              `json:"stateRootHash"`
	Nonce                  uint64                              `json:"nonce"`
	Round                  uint64                              `json:"round"`
	Epoch                  uint32                              `json:"epoch"`
	NumTxs                 uint32                              `json:"numTxs"`
	AccumulatedFees        string                              `json:"accumulatedFees,omitempty"`
	DeveloperFees          string                              `json:"developerFees,omitempty"`
	AccumulatedFeesInEpoch string                              `json:"accumulatedFeesInEpoch,omitempty"`
	DeveloperFeesInEpoch   string                              `json:"developerFeesInEpoch,omitempty"`
	Timestamp              time.Duration                       `json:"timestamp,omitempty"`
	EpochStartInfo         *api.EpochStartInfo                 `json:"epochStartInfo,omitempty"`
	EpochStartShardsData   []*api.EpochStartShardData          `json:"epochStartShardsData,omitempty"`
	ShardBlocks            []*api.NotarizedBlock               `json:"shardBlocks"`
	Transactions           []*transaction.ApiTransactionResult `json:"transactions"`
	Status                 string                              `json:"status,omitempty"`
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
	Accounts []*outport.AlteredAccount `json:"accounts"`
}
