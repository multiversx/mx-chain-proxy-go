package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetBlockByShardIDAndNonceCalled func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetBlockByHashCalled            func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled           func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
}

func (bps *BlockProcessorStub) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	return bps.GetBlockByHashCalled(shardID, hash, withTxs)
}

func (bps *BlockProcessorStub) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	return bps.GetBlockByNonceCalled(shardID, nonce, withTxs)
}

// GetAtlasBlockByShardIDAndNonce -
func (bps *BlockProcessorStub) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return bps.GetBlockByShardIDAndNonceCalled(shardID, nonce)
}

// GetHyperBlockByHash -
func (bp *BlockProcessorStub) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	panic("not implemented")
}

// GetHyperBlockByNonce -
func (bp *BlockProcessorStub) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	panic("not implemented")
}
