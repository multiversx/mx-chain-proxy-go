package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetBlockByShardIDAndNonceCalled func(shardID uint32, nonce uint64) (data.ApiBlock, error)
	GetBlockByHashCalled            func(shardID uint32, hash string, withTxs bool) (*data.GenericAPIResponse, error)
	GetBlockByNonceCalled           func(shardID uint32, nonce uint64, withTxs bool) (*data.GenericAPIResponse, error)
}

func (bps *BlockProcessorStub) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.GenericAPIResponse, error) {
	return bps.GetBlockByHashCalled(shardID, hash, withTxs)
}

func (bps *BlockProcessorStub) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.GenericAPIResponse, error) {
	return bps.GetBlockByNonceCalled(shardID, nonce, withTxs)
}

// GetAtlasBlockByShardIDAndNonce -
func (bps *BlockProcessorStub) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error) {
	return bps.GetBlockByShardIDAndNonceCalled(shardID, nonce)
}
