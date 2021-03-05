package process

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// EconomicsDataPath represents the path where an observer exposes his economics data
const EconomicsDataPath = "/network/economics"

const thresholdCountConsecutiveFails = 10

// GetEconomicsDataMetrics will return the economic metrics from cache
func (nsp *NodeStatusProcessor) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	return nsp.economicMetricsCacher.Load()
}

func (nsp *NodeStatusProcessor) getEconomicsDataMetricsFromApi() (*data.GenericAPIResponse, error) {
	metaObservers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	return nsp.getEconomicsDataMetrics(metaObservers)
}

func (nsp *NodeStatusProcessor) getEconomicsDataMetrics(observers []*data.NodeData) (*data.GenericAPIResponse, error) {
	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, EconomicsDataPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("economics data request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("economics data request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil
	}

	return nil, ErrSendingRequest
}

// StartCacheUpdate will update the economic metrics cache at a given time
func (nsp *NodeStatusProcessor) StartCacheUpdate() {
	go func() {
		countConsecutiveFails := 0
		for {
			economicMetrics, err := nsp.getEconomicsDataMetricsFromApi()
			if err != nil {
				countConsecutiveFails++
				log.Warn("economic metrics: get from API", "error", err.Error())
			}

			if countConsecutiveFails >= thresholdCountConsecutiveFails {
				nsp.economicMetricsCacher.Store(nil)
			}

			if economicMetrics != nil {
				countConsecutiveFails = 0
				nsp.economicMetricsCacher.Store(economicMetrics)
			}

			time.Sleep(nsp.cacheValidityDuration)
		}
	}()
}
