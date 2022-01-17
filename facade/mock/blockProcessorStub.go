package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetBlockByShardIDAndNonceCalled  func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetBlockByHashCalled             func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled            func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetHyperBlockByHashCalled        func(hash string) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled       func(nonce uint64) (*data.HyperblockApiResponse, error)
	GetInternalBlockByHashCalled     func(shardID uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonceCalled    func(shardID uint32, round uint64, format common.OutportFormat) (*data.InternalBlockApiResponse, error)
	GetInternalMiniBlockByHashCalled func(shardID uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error)
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
func (bps *BlockProcessorStub) GetInternalBlockByHash(shardID uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalBlockByHash(shardID, hash, format)
}

// GetInternalBlockByNonce -
func (bps *BlockProcessorStub) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalBlockByNonce(shardID, nonce, format)
}

// GetInternalMiniBlockByHash -
func (bps *BlockProcessorStub) GetInternalMiniBlockByHash(shardID uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalMiniBlockByHash(shardID, hash, format)
}
