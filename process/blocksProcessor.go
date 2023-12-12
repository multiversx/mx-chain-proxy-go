package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const (
	blockByRoundPath = "/block/by-round"
)

// BlocksProcessor handles blocks retrieving from all shards
type BlocksProcessor struct {
	proc Processor
}

// NewBlocksProcessor creates a new block processor
func NewBlocksProcessor(proc Processor) (*BlocksProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}

	return &BlocksProcessor{
		proc: proc,
	}, nil
}

// GetBlocksByRound return all blocks(from all shards) by a specific round. For each shard, a block is requested
// (from only one observer) and added in a slice of blocks => should have max blocks = no of shards.
// If there are more observers in a shard which can be queried for a block by round, we get the block from
// the first one which responds (no sanity checks are performed)
func (bp *BlocksProcessor) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	shardIDs := bp.proc.GetShardIDs()
	ret := &data.BlocksApiResponse{
		Data: data.BlocksApiResponsePayload{
			Blocks: make([]*api.Block, 0, len(shardIDs)),
		},
	}

	path := common.BuildUrlWithBlockQueryOptions(fmt.Sprintf("%s/%d", blockByRoundPath, round), options)

	for _, shardID := range shardIDs {
		observers, err := bp.proc.GetObservers(shardID, data.AvailabilityAll)
		if err != nil {
			return nil, err
		}

		for _, observer := range observers {
			block, err := bp.getBlockFromObserver(observer, path)
			if err != nil {
				log.Error("block request failed", "shard id", observer.ShardId, "observer", observer.Address, "error", err.Error())
				continue
			}

			log.Info("block requested successfully", "shard id", observer.ShardId, "observer", observer.Address, "round", round)
			ret.Data.Blocks = append(ret.Data.Blocks, block)
			break
		}
	}

	return ret, nil
}

func (bp *BlocksProcessor) getBlockFromObserver(observer *data.NodeData, path string) (*api.Block, error) {
	var response data.BlockApiResponse

	_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Block, nil
}
