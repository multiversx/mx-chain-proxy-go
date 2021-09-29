package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// BlocksProcessorStub -
type BlocksProcessorStub struct {
	GetBlocksByRoundCalled func(round uint64, withTxs bool) (*data.BlocksApiResponse, error)
}

// GetBlocksByRound -
func (bps *BlocksProcessorStub) GetBlocksByRound(round uint64, withTxs bool) (*data.BlocksApiResponse, error) {
	if bps.GetBlocksByRoundCalled != nil {
		return bps.GetBlocksByRoundCalled(round, withTxs)
	}
	return nil, nil
}
