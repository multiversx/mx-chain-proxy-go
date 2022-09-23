package process

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const (
	blockByHashPath  = "/block/by-hash"
	blockByNoncePath = "/block/by-nonce"

	internalMetaBlockByHashPath  = "/internal/%s/metablock/by-hash"
	internalShardBlockByHashPath = "/internal/%s/shardblock/by-hash"

	internalMetaBlockByNoncePath  = "/internal/%s/metablock/by-nonce"
	internalShardBlockByNoncePath = "/internal/%s/shardblock/by-nonce"

	internalMiniBlockByHashPath = "/internal/%s/miniblock/by-hash/%s/epoch/%d"

	internalStartOfEpochMetaBlockPath = "/internal/%s/startofepoch/metablock/by-epoch/%d"
)

const (
	jsonPathStr = "json"
	rawPathStr  = "raw"
)

// BlockProcessor handles blocks retrieving
type BlockProcessor struct {
	proc     Processor
	dbReader ExternalStorageConnector
}

// NewBlockProcessor will create a new block processor
func NewBlockProcessor(dbReader ExternalStorageConnector, proc Processor) (*BlockProcessor, error) {
	if check.IfNil(dbReader) {
		return nil, ErrNilDatabaseConnector
	}
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}

	return &BlockProcessor{
		dbReader: dbReader,
		proc:     proc,
	}, nil
}

// GetAtlasBlockByShardIDAndNonce return the block byte shardID and nonce
func (bp *BlockProcessor) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return bp.dbReader.GetAtlasBlockByShardIDAndNonce(shardID, nonce)
}

// GetBlockByHash will return the block based on its hash
func (bp *BlockProcessor) GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := common.BuildUrlWithBlockQueryOptions(fmt.Sprintf("%s/%s", blockByHashPath, hash), options)

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
func (bp *BlockProcessor) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := common.BuildUrlWithBlockQueryOptions(fmt.Sprintf("%s/%d", blockByNoncePath, nonce), options)

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

func (bp *BlockProcessor) getObserversOrFullHistoryNodes(shardID uint32) ([]*data.NodeData, error) {
	fullHistoryNodes, err := bp.proc.GetFullHistoryNodes(shardID)
	if err == nil {
		return fullHistoryNodes, nil
	}

	return bp.proc.GetObservers(shardID)
}

// GetHyperBlockByHash returns the hyperblock by hash
func (bp *BlockProcessor) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	builder := &HyperblockBuilder{}

	blockQueryOptions := common.BlockQueryOptions{
		WithTransactions: true,
		WithLogs:         options.WithLogs,
	}

	metaBlockResponse, err := bp.GetBlockByHash(core.MetachainShardId, hash, blockQueryOptions)
	if err != nil {
		return nil, err
	}

	metaBlock := metaBlockResponse.Data.Block
	builder.addMetaBlock(&metaBlock)

	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, blockQueryOptions)
		if err != nil {
			return nil, err
		}

		builder.addShardBlock(&shardBlockResponse.Data.Block)
	}

	hyperblock := builder.build()
	return data.NewHyperblockApiResponse(hyperblock), nil
}

// GetHyperBlockByNonce returns the hyperblock by nonce
func (bp *BlockProcessor) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	builder := &HyperblockBuilder{}

	blockQueryOptions := common.BlockQueryOptions{
		WithTransactions: true,
		WithLogs:         options.WithLogs,
	}

	metaBlockResponse, err := bp.GetBlockByNonce(core.MetachainShardId, nonce, blockQueryOptions)
	if err != nil {
		return nil, err
	}

	metaBlock := metaBlockResponse.Data.Block
	builder.addMetaBlock(&metaBlock)

	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, blockQueryOptions)
		if err != nil {
			return nil, err
		}

		builder.addShardBlock(&shardBlockResponse.Data.Block)
	}

	hyperblock := builder.build()
	return data.NewHyperblockApiResponse(hyperblock), nil
}

// GetInternalBlockByHash will return the internal block based on its hash
func (bp *BlockProcessor) GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path, err := getInternalBlockByHashPath(shardID, format, hash)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var response data.InternalBlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

func getInternalBlockByHashPath(shardID uint32, format common.OutputFormat, hash string) (string, error) {
	var path string

	outputStr, err := getOutputFormat(format)
	if err != nil {
		return "", err
	}

	if shardID == core.MetachainShardId {
		path = fmt.Sprintf(internalMetaBlockByHashPath, outputStr)
	} else {
		path = fmt.Sprintf(internalShardBlockByHashPath, outputStr)
	}

	return fmt.Sprintf("%s/%s", path, hash), nil
}

// GetInternalBlockByNonce will return the internal block based on its nonce
func (bp *BlockProcessor) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path, err := getInternalBlockByNoncePath(shardID, format, nonce)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var response data.InternalBlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal block request", "shard id", observer.ShardId, "round", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

func getInternalBlockByNoncePath(shardID uint32, format common.OutputFormat, nonce uint64) (string, error) {
	var path string

	outputStr, err := getOutputFormat(format)
	if err != nil {
		return "", err
	}

	if shardID == core.MetachainShardId {
		path = fmt.Sprintf(internalMetaBlockByNoncePath, outputStr)
	} else {
		path = fmt.Sprintf(internalShardBlockByNoncePath, outputStr)
	}

	return fmt.Sprintf("%s/%d", path, nonce), nil
}

// GetInternalMiniBlockByHash will return the miniblock based on its hash
func (bp *BlockProcessor) GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	outputStr, err := getOutputFormat(format)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf(internalMiniBlockByHashPath, outputStr, hash, epoch)

	for _, observer := range observers {
		var response data.InternalMiniBlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("miniblock request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("miniblock request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

func getOutputFormat(format common.OutputFormat) (string, error) {
	var outputStr string

	switch format {
	case common.Internal:
		outputStr = jsonPathStr
	case common.Proto:
		outputStr = rawPathStr
	default:
		return "", ErrInvalidOutputFormat
	}

	return outputStr, nil
}

// GetInternalStartOfEpochMetaBlock will return the internal start of epoch meta block based on epoch
func (bp *BlockProcessor) GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	outputStr, err := getOutputFormat(format)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(internalStartOfEpochMetaBlockPath, outputStr, epoch)

	for _, observer := range observers {
		var response data.InternalBlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal block request", "shard id", observer.ShardId, "epoch", epoch, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}
