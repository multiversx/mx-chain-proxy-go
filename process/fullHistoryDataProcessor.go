package process

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// NetworkStatusPath represents the path where an observer exposes his network metrics
const blockByNoncePath = "/block/by-nonce"

// NetworkConfigPath represents the path where an observer exposes his network metrics
const blockByHashPath = "/block/by-hash"

// NodeStatusProcessor handles the action needed for fetching data related to status metrics from nodes
type fullHistoryDataProcessor struct {
	proc Processor
}

// NewNodeStatusProcessor creates a new instance of NodeStatusProcessor
func NewFullHistoryDataProcessor(processor Processor) (*fullHistoryDataProcessor, error) {
	if check.IfNil(processor) {
		return nil, ErrNilCoreProcessor
	}

	return &fullHistoryDataProcessor{
		proc: processor,
	}, nil
}

// GetNetworkStatusMetrics will simply forward the network status metrics from an observer in the given shard
func (nsp *fullHistoryDataProcessor) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", blockByHashPath, hash)
	if withTxs {
		path += "?withTxs=true"
	}

	for _, observer := range observers {
		var response *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("full history node - block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("full history node - block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return response, nil

	}

	return nil, ErrSendingRequest
}

func (nsp *fullHistoryDataProcessor) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%d", blockByNoncePath, nonce)
	if withTxs {
		path += "?withTxs=true"
	}

	for _, observer := range observers {
		var response *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("full history node - block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("full history node - block request", "shard id", observer.ShardId, "nonce", nonce, "observer", observer.Address)
		return response, nil

	}

	return nil, ErrSendingRequest
}
