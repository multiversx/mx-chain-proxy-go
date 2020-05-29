package process

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// HeartBeatPath represents the path where an observer exposes his heartbeat status
const HeartBeatPath = "/node/heartbeatstatus"

// HeartbeatProcessor is able to process transaction requests
type HeartbeatProcessor struct {
	proc                  Processor
	cacher                HeartbeatCacheHandler
	cacheValidityDuration time.Duration
}

// NewHeartbeatProcessor creates a new instance of HeartbeatProcessor
func NewHeartbeatProcessor(
	proc Processor,
	cacher HeartbeatCacheHandler,
	cacheValidityDuration time.Duration,
) (*HeartbeatProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(cacher) {
		return nil, ErrNilHeartbeatCacher
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}
	hbp := &HeartbeatProcessor{
		proc:                  proc,
		cacher:                cacher,
		cacheValidityDuration: cacheValidityDuration,
	}

	return hbp, nil
}

// GetHeartbeatData will simply forward the heartbeat status from an observer
func (hbp *HeartbeatProcessor) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	heartbeatsToReturn, err := hbp.cacher.LoadHeartbeats()
	if err == nil {
		return heartbeatsToReturn, nil
	}

	log.Info("heartbeat: cannot get from cache. Will fetch from API", "error", err.Error())

	return hbp.getHeartbeatsFromApi()
}

func (hbp *HeartbeatProcessor) getHeartbeatsFromApi() (*data.HeartbeatResponse, error) {
	observers := hbp.proc.GetAllObservers()

	var response data.HeartbeatApiResponse
	var err error
	for _, observer := range observers {
		err = hbp.proc.CallGetRestEndPoint(observer.Address, HeartBeatPath, &response)
		if err == nil {
			log.Info("heartbeat fetched from API", "observer", observer.Address)
			return &response.Data, nil
		}
		log.Error("heartbeat", "observer", observer.Address, "error", "no response")
	}
	return nil, ErrHeartbeatNotAvailable
}

// StartCacheUpdate will start the updating of the cache from the API at a given period
func (hbp *HeartbeatProcessor) StartCacheUpdate() {
	go func() {
		for {
			hbts, err := hbp.getHeartbeatsFromApi()
			if err != nil {
				log.Warn("heartbeat: get from API", "error", err.Error())
			}

			if hbts != nil {
				err = hbp.cacher.StoreHeartbeats(hbts)
				if err != nil {
					log.Warn("heartbeat: store in cache", "error", err.Error())
				}
			}

			time.Sleep(hbp.cacheValidityDuration)
		}
	}()
}
