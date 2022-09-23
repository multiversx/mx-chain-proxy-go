package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// BlockProcessorStub -
type BlockProcessorStub struct {
	GetBlockByShardIDAndNonceCalled        func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetBlockByHashCalled                   func(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled                  func(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetHyperBlockByHashCalled              func(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled             func(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetInternalBlockByHashCalled           func(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonceCalled          func(shardID uint32, round uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalMiniBlockByHashCalled       func(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error)
	GetInternalStartOfEpochMetaBlockCalled func(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
}

func (bps *BlockProcessorStub) GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return bps.GetBlockByHashCalled(shardID, hash, options)
}

func (bps *BlockProcessorStub) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return bps.GetBlockByNonceCalled(shardID, nonce, options)
}

// GetAtlasBlockByShardIDAndNonce -
func (bps *BlockProcessorStub) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return bps.GetBlockByShardIDAndNonceCalled(shardID, nonce)
}

// GetHyperBlockByHash -
func (bps *BlockProcessorStub) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	if bps.GetHyperBlockByHashCalled != nil {
		return bps.GetHyperBlockByHashCalled(hash, options)
	}

	panic("not implemented: GetHyperBlockByHash")
}

// GetHyperBlockByNonce -
func (bps *BlockProcessorStub) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	if bps.GetHyperBlockByNonceCalled != nil {
		return bps.GetHyperBlockByNonceCalled(nonce, options)
	}

	panic("not implemented: GetHyperBlockByNonce")
}

// GetInternalBlockByHash -
func (bps *BlockProcessorStub) GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalBlockByHashCalled(shardID, hash, format)
}

// GetInternalBlockByNonce -
func (bps *BlockProcessorStub) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalBlockByNonceCalled(shardID, nonce, format)
}

// GetInternalMiniBlockByHash -
func (bps *BlockProcessorStub) GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
	return bps.GetInternalMiniBlockByHashCalled(shardID, hash, epoch, format)
}

// GetInternalStartOfEpochMetaBlock -
func (bps *BlockProcessorStub) GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return bps.GetInternalStartOfEpochMetaBlockCalled(epoch, format)
}
