package process

import (
	"context"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// ValidatorStatisticsPath represents the path where an observer exposes his validator statistics data
const ValidatorStatisticsPath = "/validator/statistics"

// ValidatorStatisticsProcessor is able to process validator statistics data requests
type ValidatorStatisticsProcessor struct {
	proc                  Processor
	cacher                ValidatorStatisticsCacheHandler
	cacheValidityDuration time.Duration
	cancelFunc            func()
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
func (vsp *ValidatorStatisticsProcessor) GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error) {
	valStatsToReturn, err := vsp.cacher.LoadValStats()
	if err == nil {
		return &data.ValidatorStatisticsResponse{Statistics: valStatsToReturn}, nil
	}

	log.Info("validator statistics: cannot get from cache. Will fetch from API", "error", err.Error())

	return vsp.getValidatorStatisticsFromApi()
}

func (vsp *ValidatorStatisticsProcessor) getValidatorStatisticsFromApi() (*data.ValidatorStatisticsResponse, error) {
	observers, errFetchObs := vsp.proc.GetObservers(core.MetachainShardId, data.AvailabilityRecent)
	if errFetchObs != nil {
		return nil, errFetchObs
	}

	var valStatsResponse data.ValidatorStatisticsApiResponse
	var err error
	for _, observer := range observers {
		_, err = vsp.proc.CallGetRestEndPoint(observer.Address, ValidatorStatisticsPath, &valStatsResponse)
		if err == nil {
			log.Info("validator statistics fetched from API", "observer", observer.Address)
			return &valStatsResponse.Data, nil
		}
		log.Error("validator statistics", "observer", observer.Address, "error", "no response")
	}
	return nil, ErrValidatorStatisticsNotAvailable
}

// StartCacheUpdate will start the updating of the cache from the API at a given period
func (vsp *ValidatorStatisticsProcessor) StartCacheUpdate() {
	if vsp.cancelFunc != nil {
		log.Error("ValidatorStatisticsProcessor - cache update already started")
		return
	}

	var ctx context.Context
	ctx, vsp.cancelFunc = context.WithCancel(context.Background())

	go func(ctx context.Context) {
		timer := time.NewTimer(vsp.cacheValidityDuration)
		defer timer.Stop()

		vsp.handleCacheUpdate()

		for {
			timer.Reset(vsp.cacheValidityDuration)

			select {
			case <-timer.C:
				vsp.handleCacheUpdate()
			case <-ctx.Done():
				log.Debug("finishing ValidatorStatisticsProcessor cache update...")
				return
			}
		}
	}(ctx)
}

func (vsp *ValidatorStatisticsProcessor) handleCacheUpdate() {
	valStats, err := vsp.getValidatorStatisticsFromApi()
	if err != nil {
		log.Warn("validator statistics: get from API", "error", err.Error())
	}

	if valStats != nil {
		err = vsp.cacher.StoreValStats(valStats.Statistics)
		if err != nil {
			log.Warn("validator statistics: store in cache", "error", err.Error())
		}
	}
}

// Close will handle the closing of the cache update go routine
func (vsp *ValidatorStatisticsProcessor) Close() error {
	if vsp.cancelFunc != nil {
		vsp.cancelFunc()
	}

	return nil
}
