package process

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const (
	blocksByRoundPath = "/blocks/by-round"
)

type BlocksProcessor struct {
	proc Processor
}

func NewBlocksProcessor(proc Processor) (*BlocksProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}

	return &BlocksProcessor{
		proc: nil,
	}, nil
}

func (bp *BlocksProcessor) GetBlocksByRound(round uint64, withTxs bool) (*data.BlocksApiResponse, error) {
	shardIDs := bp.proc.GetShardIDs()
	ret := &data.BlocksApiResponse{
		Data: data.BlocksApiResponsePayload{
			Blocks: make([]*data.Block, 0, len(shardIDs)),
		},
	}

	path := fmt.Sprintf("%s/%d", blocksByRoundPath, round)
	if withTxs {
		path += withTxsParamTrue
	}

	for shardID := range shardIDs {
		observers, err := bp.proc.GetObservers(uint32(shardID))
		if err != nil {
			return nil, err
		}

		for _, observer := range observers {
			block, err := bp.getBlockFromObserver(observer, path)
			if err != nil {
				log.Error("block request", "shard id", observer.ShardId, "observer", observer.Address, "error", err.Error())
				continue
			}

			log.Info("block request", "shard id", observer.ShardId, "observer", observer.Address, "round", round)
			ret.Data.Blocks = append(ret.Data.Blocks, block)
			break
		}
	}

	return ret, nil
}

func (bp *BlocksProcessor) getBlockFromObserver(observer *data.NodeData, path string) (*data.Block, error) {
	response := data.BlockApiResponse{}

	_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Block, nil
}
