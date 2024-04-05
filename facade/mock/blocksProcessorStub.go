package mock

import (
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// BlocksProcessorStub -
type BlocksProcessorStub struct {
	GetBlocksByRoundCalled func(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error)
}

// GetBlocksByRound -
func (bps *BlocksProcessorStub) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	if bps.GetBlocksByRoundCalled != nil {
		return bps.GetBlocksByRoundCalled(round, options)
	}
	return nil, nil
}
