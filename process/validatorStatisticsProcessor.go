package process

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ValidatorStatisticsPath represents the path where an observer exposes his validator statistics data
const ValidatorStatisticsPath = "/validator/statistics"

// validatorStatisticsProcessor is able to process validator statistics data requests
type validatorStatisticsProcessor struct {
	proc                  Processor
	cacher                ValidatorStatisticsCacheHandler
	cacheValidityDuration time.Duration
}

// NewValidatorStatisticsProcessor creates a new instance of validatorStatisticsProcessor
func NewValidatorStatisticsProcessor(
	proc Processor,
	cacher ValidatorStatisticsCacheHandler,
	cacheValidityDuration time.Duration,
) (*validatorStatisticsProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(cacher) {
		return nil, ErrNilValidatorStatisticsCacher
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}
	hbp := &validatorStatisticsProcessor{
		proc:                  proc,
		cacher:                cacher,
		cacheValidityDuration: cacheValidityDuration,
	}

	return hbp, nil
}

// GetValidatorStatistics will simply forward the validator statistics data from an observer
func (hbp *validatorStatisticsProcessor) GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error) {
	valStatsToReturn, err := hbp.cacher.LoadValStats()
	if err == nil {
		return &data.ValidatorStatisticsResponse{Statistics: valStatsToReturn}, nil
	}

	log.Info("validator statistics: cannot get from cache. Will fetch from API", "error", err.Error())

	return hbp.getValidatorStatisticsFromApi()
}

func (hbp *validatorStatisticsProcessor) getValidatorStatisticsFromApi() (*data.ValidatorStatisticsResponse, error) {
	observers := hbp.proc.GetAllObservers()

	var valStatsResponse data.ValidatorStatisticsResponse
	var err error
	for _, observer := range observers {
		err = hbp.proc.CallGetRestEndPoint(observer.Address, ValidatorStatisticsPath, &valStatsResponse)
		if err == nil {
			log.Info("validator statistics fetched from API", "observer", observer.Address)
			return &valStatsResponse, nil
		}
		log.Error("validator statistics", "observer", observer.Address, "error", "no response")
	}
	return nil, ErrValidatorStatisticsNotAvailable
}

// StartCacheUpdate will start the updating of the cache from the API at a given period
func (hbp *validatorStatisticsProcessor) StartCacheUpdate() {
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
