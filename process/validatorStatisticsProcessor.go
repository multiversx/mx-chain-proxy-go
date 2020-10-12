package process

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ValidatorStatisticsPath represents the path where an observer exposes his validator statistics data
const ValidatorStatisticsPath = "/validator/statistics"

// ValidatorStatisticsProcessor is able to process validator statistics data requests
type ValidatorStatisticsProcessor struct {
	proc                  Processor
	cacher                ValidatorStatisticsCacheHandler
	cacheValidityDuration time.Duration
}

// NewValidatorStatisticsProcessor creates a new instance of ValidatorStatisticsProcessor
func NewValidatorStatisticsProcessor(
	proc Processor,
	cacher ValidatorStatisticsCacheHandler,
	cacheValidityDuration time.Duration,
) (*ValidatorStatisticsProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(cacher) {
		return nil, ErrNilValidatorStatisticsCacher
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}
	hbp := &ValidatorStatisticsProcessor{
		proc:                  proc,
		cacher:                cacher,
		cacheValidityDuration: cacheValidityDuration,
	}

	return hbp, nil
}

// GetValidatorStatistics will simply forward the validator statistics data from an observer
func (hbp *ValidatorStatisticsProcessor) GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error) {
	valStatsToReturn, err := hbp.cacher.LoadValStats()
	if err == nil {
		return &data.ValidatorStatisticsResponse{Statistics: valStatsToReturn}, nil
	}

	log.Info("validator statistics: cannot get from cache. Will fetch from API", "error", err.Error())

	return hbp.getValidatorStatisticsFromApi()
}

func (hbp *ValidatorStatisticsProcessor) getValidatorStatisticsFromApi() (*data.ValidatorStatisticsResponse, error) {
	observers, errFetchObs := hbp.proc.GetObservers(core.MetachainShardId)
	if errFetchObs != nil {
		return nil, errFetchObs
	}

	var valStatsResponse data.ValidatorStatisticsApiResponse
	var err error
	for _, observer := range observers {
		_, err = hbp.proc.CallGetRestEndPoint(observer.Address, ValidatorStatisticsPath, &valStatsResponse)
		if err == nil {
			log.Info("validator statistics fetched from API", "observer", observer.Address)
			return &valStatsResponse.Data, nil
		}
		log.Error("validator statistics", "observer", observer.Address, "error", "no response")
	}
	return nil, ErrValidatorStatisticsNotAvailable
}

// StartCacheUpdate will start the updating of the cache from the API at a given period
func (hbp *ValidatorStatisticsProcessor) StartCacheUpdate() {
	go func() {
		for {
			valStats, err := hbp.getValidatorStatisticsFromApi()
			if err != nil {
				log.Warn("validator statistics: get from API", "error", err.Error())
			}

			if valStats != nil {
				err = hbp.cacher.StoreValStats(valStats.Statistics)
				if err != nil {
					log.Warn("validator statistics: store in cache", "error", err.Error())
				}
			}

			time.Sleep(hbp.cacheValidityDuration)
		}
	}()
}
