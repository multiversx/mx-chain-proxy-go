package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetBlockByShardIDAndNonceCalled func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetBlockByHashCalled            func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled           func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetHyperBlockByHashCalled       func(hash string) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled      func(nonce uint64) (*data.HyperblockApiResponse, error)
	GetInternalBlockByHashCalled    func(shardID uint32, hash string) (*data.InternalBlockApiResponse, error)
	GetRawBlockByHashCalled         func(shardID uint32, hash string) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonceCalled   func(shardID uint32, round uint64) (*data.InternalBlockApiResponse, error)
	GetRawBlockByNonceCalled        func(shardID uint32, round uint64) (*data.InternalBlockApiResponse, error)
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
func (bps *BlockProcessorStub) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	if bps.GetHyperBlockByHashCalled != nil {
		return bps.GetHyperBlockByHashCalled(hash)
	}

	panic("not implemented: GetHyperBlockByHash")
}

// GetHyperBlockByNonce -
func (bps *BlockProcessorStub) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	if bps.GetHyperBlockByNonceCalled != nil {
		return bps.GetHyperBlockByNonceCalled(nonce)
	}

	panic("not implemented: GetHyperBlockByNonce")
}

// GetInternalBlockByHash -
func (bps *BlockProcessorStub) GetInternalBlockByHash(shardID uint32, hash string) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalBlockByHash(shardID, hash)
}

// GetInternalBlockByNonce -
func (bps *BlockProcessorStub) GetInternalBlockByNonce(shardID uint32, nonce uint64) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalBlockByNonce(shardID, nonce)
}

// GetRawBlockByHash -
func (bps *BlockProcessorStub) GetRawBlockByHash(shardID uint32, hash string) (*data.InternalBlockApiResponse, error) {
	return bps.GetRawBlockByHash(shardID, hash)
}

// GetRawBlockByNonce -
func (bps *BlockProcessorStub) GetRawBlockByNonce(shardID uint32, nonce uint64) (*data.InternalBlockApiResponse, error) {
	return bps.GetRawBlockByNonce(shardID, nonce)
}
