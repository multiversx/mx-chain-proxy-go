package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type BlocksProcessorStub struct {
	GetBlocksByRoundCalled func(round uint64, withTxs bool) (*data.BlocksApiResponse, error)
}

func (bps *BlocksProcessorStub) GetBlocksByRound(round uint64, withTxs bool) (*data.BlocksApiResponse, error) {
	if bps.GetBlocksByRoundCalled != nil {
		return bps.GetBlocksByRoundCalled(round, withTxs)
	}
	return nil, nil
}
