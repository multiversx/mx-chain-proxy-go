package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const (
	blockByHashPath  = "/block/by-hash"
	blockByNoncePath = "/block/by-nonce"

	internalMetaBlockByHashPath  = "/internal/%s/metablock/by-hash"
	internalShardBlockByHashPath = "/internal/%s/shardblock/by-hash"

	internalMetaBlockByNoncePath  = "/internal/%s/metablock/by-nonce"
	internalShardBlockByNoncePath = "/internal/%s/shardblock/by-nonce"

	internalMiniBlockByHashPath = "/internal/%s/miniblock/by-hash/%s/epoch/%d"

	internalStartOfEpochMetaBlockPath      = "/internal/%s/startofepoch/metablock/by-epoch/%d"
	internalStartOfEpochValidatorsInfoPath = "/internal/json/startofepoch/validators/by-epoch/%d"

	alteredAccountByBlockNonce = "/block/altered-accounts/by-nonce"
	alteredAccountByBlockHash  = "/block/altered-accounts/by-hash"
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

	response := data.BlockApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}

// GetBlockByNonce will return the block based on the nonce
func (bp *BlockProcessor) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := common.BuildUrlWithBlockQueryOptions(fmt.Sprintf("%s/%d", blockByNoncePath, nonce), options)

	response := data.BlockApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "nonce", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}

func (bp *BlockProcessor) getObserversOrFullHistoryNodes(shardID uint32) ([]*data.NodeData, error) {
	fullHistoryNodes, err := bp.proc.GetFullHistoryNodes(shardID, data.AvailabilityAll)
	if err == nil {
		return fullHistoryNodes, nil
	}

	return bp.proc.GetObservers(shardID, data.AvailabilityAll)
}

// GetHyperBlockByHash returns the hyperblock by hash
func (bp *BlockProcessor) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	builder := &hyperblockBuilder{}

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

	err = bp.addShardBlocks(metaBlock, builder, options, blockQueryOptions)
	if err != nil {
		return nil, err
	}

	hyperblock := builder.build(options.NotarizedAtSource)
	return data.NewHyperblockApiResponse(hyperblock), nil
}

func (bp *BlockProcessor) addShardBlocks(
	metaBlock api.Block,
	builder *hyperblockBuilder,
	options common.HyperblockQueryOptions,
	blockQueryOptions common.BlockQueryOptions,
) error {
	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, blockQueryOptions)
		if err != nil {
			return err
		}

		alteredAccounts, err := bp.getAlteredAccountsIfNeeded(options, notarizedBlock)
		if err != nil {
			return err
		}

		builder.addShardBlock(&shardBlockWithAlteredAccounts{
			shardBlock:      &shardBlockResponse.Data.Block,
			alteredAccounts: alteredAccounts,
		})
	}

	return nil
}

func (bp *BlockProcessor) getAlteredAccountsIfNeeded(options common.HyperblockQueryOptions, notarizedBlock *api.NotarizedBlock) ([]*alteredAccount.AlteredAccount, error) {
	ret := make([]*alteredAccount.AlteredAccount, 0)
	if !options.WithAlteredAccounts {
		return ret, nil
	}

	alteredAccountsApiResponse, err := bp.GetAlteredAccountsByHash(notarizedBlock.Shard, notarizedBlock.Hash, options.AlteredAccountsOptions)
	if err != nil {
		return nil, err
	}

	return alteredAccountsApiResponse.Data.Accounts, nil
}

// GetHyperBlockByNonce returns the hyperblock by nonce
func (bp *BlockProcessor) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	builder := &hyperblockBuilder{}

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

	err = bp.addShardBlocks(metaBlock, builder, options, blockQueryOptions)
	if err != nil {
		return nil, err
	}

	hyperblock := builder.build(options.NotarizedAtSource)
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

	response := data.InternalBlockApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
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

	response := data.InternalBlockApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal block request", "shard id", observer.ShardId, "round", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
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

	response := data.InternalMiniBlockApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("miniblock request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("miniblock request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
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

	response := data.InternalBlockApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal block request", "shard id", observer.ShardId, "epoch", epoch, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}

// GetInternalStartOfEpochValidatorsInfo will return the internal start of epoch validators info based on epoch
func (bp *BlockProcessor) GetInternalStartOfEpochValidatorsInfo(epoch uint32) (*data.ValidatorsInfoApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(internalStartOfEpochValidatorsInfoPath, epoch)

	response := data.ValidatorsInfoApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("internal validators info request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("internal validators info request", "shard id", observer.ShardId, "epoch", epoch, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}

// GetAlteredAccountsByNonce will return altered accounts by block nonce
func (bp *BlockProcessor) GetAlteredAccountsByNonce(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	observers, err := bp.proc.GetObservers(shardID, data.AvailabilityAll)
	if err != nil {
		return nil, err
	}
	path := common.BuildUrlWithAlteredAccountsQueryOptions(fmt.Sprintf("%s/%d", alteredAccountByBlockNonce, nonce), options)

	response := data.AlteredAccountsApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("altered accounts request by nonce", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("altered accounts request by nonce", "shard id", observer.ShardId, "nonce", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}

// GetAlteredAccountsByHash will return altered accounts by block hash
func (bp *BlockProcessor) GetAlteredAccountsByHash(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	observers, err := bp.proc.GetObservers(shardID, data.AvailabilityAll)
	if err != nil {
		return nil, err
	}
	path := common.BuildUrlWithAlteredAccountsQueryOptions(fmt.Sprintf("%s/%s", alteredAccountByBlockHash, hash), options)

	response := data.AlteredAccountsApiResponse{}
	for _, observer := range observers {

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("altered accounts request by hash", "observer", observer.Address, "hash", hash, "error", err.Error())
			continue
		}

		log.Info("altered accounts request by hash", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}
