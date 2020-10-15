package process

import (
	"fmt"
	"math"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const (
	blockByHashPath  = "/block/by-hash"
	blockByNoncePath = "/block/by-nonce"
	withTxsParamTrue = "?withTxs=true"
)

type blockProcessor struct {
	proc                                      Processor
	dbReader                                  ExternalStorageConnector
	getLatestFullySynchronizedHyperblockNonce func() (uint64, error)
}

// NewBlockProcessor will create a new block processor
func NewBlockProcessor(
	dbReader ExternalStorageConnector,
	proc Processor,
	getLatestFullySynchronizedHyperblockNonce func() (uint64, error),
) (*blockProcessor, error) {
	if check.IfNil(dbReader) {
		return nil, ErrNilDatabaseConnector
	}
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	// if function is nil will return always MaxUint64 -> result of function will be ignored
	if getLatestFullySynchronizedHyperblockNonce == nil {
		getLatestFullySynchronizedHyperblockNonce = func() (uint64, error) {
			return uint64(math.MaxUint64), nil
		}
	}

	return &blockProcessor{
		dbReader: dbReader,
		proc:     proc,
		getLatestFullySynchronizedHyperblockNonce: getLatestFullySynchronizedHyperblockNonce,
	}, nil
}

// GetAtlasBlockByShardIDAndNonce return the block byte shardID and nonce
func (bp *blockProcessor) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return bp.dbReader.GetAtlasBlockByShardIDAndNonce(shardID, nonce)
}

// GetBlockByHash will return the block based on its hash
func (bp *blockProcessor) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", blockByHashPath, hash)
	if withTxs {
		path += withTxsParamTrue
	}

	for _, observer := range observers {
		var response data.BlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

// GetBlockByNonce will return the block based on the nonce
func (bp *blockProcessor) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%d", blockByNoncePath, nonce)
	if withTxs {
		path += withTxsParamTrue
	}

	for _, observer := range observers {
		var response data.BlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "nonce", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

func (bp *blockProcessor) getObserversOrFullHistoryNodes(shardID uint32) ([]*data.NodeData, error) {
	fullHistoryNodes, err := bp.proc.GetFullHistoryNodes(shardID)
	if err == nil {
		return fullHistoryNodes, nil
	}

	return bp.proc.GetObservers(shardID)
}

// GetHyperBlockByHash returns the hyperblock by hash
func (bp *blockProcessor) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	builder := &hyperblockBuilder{}

	metaBlockResponse, err := bp.GetBlockByHash(core.MetachainShardId, hash, true)
	if err != nil {
		return nil, err
	}

	if highestBlockNonceThatCanBeReturned, err := bp.getLatestFullySynchronizedHyperblockNonce(); err == nil {
		if metaBlockResponse.Data.Block.Nonce > highestBlockNonceThatCanBeReturned {
			return nil, fmt.Errorf("%w with hash %s: has nonce (%d) greater than highest block nonce that can be returned(%d)",
				ErrCannotGetHyperblock,
				hash,
				metaBlockResponse.Data.Block.Nonce,
				highestBlockNonceThatCanBeReturned,
			)
		}
	}

	metaBlock := metaBlockResponse.Data.Block
	builder.addMetaBlock(&metaBlock)

	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, true)
		if err != nil {
			return nil, err
		}

		builder.addShardBlock(&shardBlockResponse.Data.Block)
	}

	hyperblock := builder.build()
	return data.NewHyperblockApiResponse(hyperblock), nil
}

// GetHyperBlockByNonce returns the hyperblock by nonce
func (bp *blockProcessor) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	if highestBlockNonceThatCanBeReturned, err := bp.getLatestFullySynchronizedHyperblockNonce(); err == nil {
		if nonce > highestBlockNonceThatCanBeReturned {
			return nil, fmt.Errorf("%w: has nonce (%d) greater than highest block nonce that can be returned(%d)",
				ErrCannotGetHyperblock,
				nonce,
				highestBlockNonceThatCanBeReturned,
			)
		}
	}

	builder := &hyperblockBuilder{}

	metaBlockResponse, err := bp.GetBlockByNonce(core.MetachainShardId, nonce, true)
	if err != nil {
		return nil, err
	}

	metaBlock := metaBlockResponse.Data.Block
	builder.addMetaBlock(&metaBlock)

	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, true)
		if err != nil {
			return nil, err
		}

		builder.addShardBlock(&shardBlockResponse.Data.Block)
	}

	hyperblock := builder.build()
	return data.NewHyperblockApiResponse(hyperblock), nil
}
