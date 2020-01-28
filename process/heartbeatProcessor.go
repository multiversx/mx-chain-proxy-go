package process

import (
	"time"

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

// NewHeartbeatProcessor creates a new instance of TransactionProcessor
func NewHeartbeatProcessor(
	proc Processor,
	cacher HeartbeatCacheHandler,
	cacheValidityDuration time.Duration,
) (*HeartbeatProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}
	if cacher == nil || cacher.IsInterfaceNil() {
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
	heartbeatsToReturn, err := hbp.cacher.Heartbeats()
	if err == nil {
		return heartbeatsToReturn, nil
	}

	log.Info("heartbeat: cannot get from cache. Will fetch from API", "error", err.Error())

	return hbp.getHeartbeatsFromApi()
}

func (hbp *HeartbeatProcessor) getHeartbeatsFromApi() (*data.HeartbeatResponse, error) {
	if hbp.proc.AreObserversBalanced() {
		observersRing := hbp.proc.GetAllObserversRing()
		numTries := 0
		for numTries < observersRing.Len() {
			hbtRes, err := hbp.callApiEndpointForHeartbeat(observersRing.Next())
			if err == nil {
				return hbtRes, nil
			}
			numTries++
		}

		return nil, ErrHeartbeatNotAvailable
	} else {
		observers, err := hbp.proc.GetAllObservers()
		if err != nil {
			return nil, err
		}
		for _, observer := range observers {
			hbtRes, err := hbp.callApiEndpointForHeartbeat(observer.Address)
			if err == nil {
				return hbtRes, nil
			}
		}

		return nil, ErrHeartbeatNotAvailable
	}
}

func (hbp *HeartbeatProcessor) callApiEndpointForHeartbeat(observerAddress string) (*data.HeartbeatResponse, error) {
	var heartbeatResponse data.HeartbeatResponse
	err := hbp.proc.CallGetRestEndPoint(observerAddress, HeartBeatPath, &heartbeatResponse)
	if err == nil {
		log.Info("heartbeat fetched from API", "observer", observerAddress)
		return &heartbeatResponse, nil
	}

	log.Error("heartbeat", "observer", observerAddress, "error", "no response")
	return nil, ErrNoResponseFromObserver
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
