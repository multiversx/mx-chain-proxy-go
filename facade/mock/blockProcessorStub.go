package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetHighestBlockNonceCalled func() (uint64, error)
	GetBlockByNonceCalled      func(nonce uint64) (data.ApiBlock, error)
}

// GetHighestBlockNonce -
func (bps *BlockProcessorStub) GetHighestBlockNonce() (uint64, error) {
	return bps.GetHighestBlockNonceCalled()
}

// GetBlockByNonce -
func (bps *BlockProcessorStub) GetBlockByNonce(nonce uint64) (data.ApiBlock, error) {
	return bps.GetBlockByNonceCalled(nonce)
}
