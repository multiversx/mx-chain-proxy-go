package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetBlockByShardIDAndNonceCalled func(shardID uint32, nonce uint64) (data.ApiBlock, error)
}

// GetBlockByShardIDAndNonce -
func (bps *BlockProcessorStub) GetBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error) {
	return bps.GetBlockByShardIDAndNonceCalled(shardID, nonce)
}
