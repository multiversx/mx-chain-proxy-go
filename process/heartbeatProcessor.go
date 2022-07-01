package process

import (
	"context"
	"sort"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// HeartBeatPath represents the path where an observer exposes his heartbeat status
const HeartBeatPath = "/node/heartbeatstatus"

// HeartbeatProcessor is able to process transaction requests
type HeartbeatProcessor struct {
	proc                  Processor
	cacher                HeartbeatCacheHandler
	cacheValidityDuration time.Duration
	cancelFunc            func()
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
	shardIDs := hbp.proc.GetShardIDs()

	responseMap := make(map[string]data.PubKeyHeartbeat)
	for _, shard := range shardIDs {
		observers, err := hbp.proc.GetObservers(shard)
		if err != nil {
			log.Error("could not get observers", "shard", shard, "error", err.Error())
			continue
		}

		var response data.HeartbeatApiResponse
		for _, observer := range observers {
			_, err = hbp.proc.CallGetRestEndPoint(observer.Address, HeartBeatPath, &response)
			if err == nil {
				hbp.addMessagesToMap(responseMap, response.Data.Heartbeats, shard)
				break
			}

			log.Error("heartbeat", "observer", observer.Address, "shard", shard, "error", err.Error())
		}
	}

	if len(responseMap) == 0 {
		return nil, ErrHeartbeatNotAvailable
	}

	return hbp.mapToResponse(responseMap), nil
}

func (hbp *HeartbeatProcessor) addMessagesToMap(responseMap map[string]data.PubKeyHeartbeat, heartbeats []data.PubKeyHeartbeat, observerShard uint32) {
	for _, heartbeatMessage := range heartbeats {
		// TODO: fix these merges when the heartbeat v2 will be active. Within this implementation, if a shard won't
		// respond, then the final heartbeat message won't include data from that shard. Analyze if returning error
		// in case of an unresponsive shard is ok, or a better solution is to be found

		//isMessageFromCurrentShard := heartbeatMessage.ReceivedShardID == observerShard
		//if !isMessageFromCurrentShard {
		//	continue
		//}

		_, found := responseMap[heartbeatMessage.PublicKey]
		if !found {
			responseMap[heartbeatMessage.PublicKey] = heartbeatMessage
		}
	}
}

func (hbp *HeartbeatProcessor) mapToResponse(responseMap map[string]data.PubKeyHeartbeat) *data.HeartbeatResponse {
	heartbeats := make([]data.PubKeyHeartbeat, 0)
	for _, heartbeatMessage := range responseMap {
		heartbeats = append(heartbeats, heartbeatMessage)
	}

	sort.Slice(heartbeats, func(i, j int) bool {
		return heartbeats[i].PublicKey < heartbeats[j].PublicKey
	})

	return &data.HeartbeatResponse{
		Heartbeats: heartbeats,
	}
}

// StartCacheUpdate will start the updating of the cache from the API at a given period
func (hbp *HeartbeatProcessor) StartCacheUpdate() {
	if hbp.cancelFunc != nil {
		log.Error("HeartbeatProcessor - cache update already started")
		return
	}

	var ctx context.Context
	ctx, hbp.cancelFunc = context.WithCancel(context.Background())

	go func(ctx context.Context) {
		timer := time.NewTimer(hbp.cacheValidityDuration)
		defer timer.Stop()

		hbp.handleHeartbeatCacheUpdate()

		for {
			timer.Reset(hbp.cacheValidityDuration)

			select {
			case <-timer.C:
				hbp.handleHeartbeatCacheUpdate()
			case <-ctx.Done():
				log.Debug("finishing HeartbeatProcessor cache update...")
				return
			}
		}
	}(ctx)
}

func (hbp *HeartbeatProcessor) handleHeartbeatCacheUpdate() {
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
}

// Close will handle the closing of the cache update go routine
func (hbp *HeartbeatProcessor) Close() error {
	if hbp.cancelFunc != nil {
		hbp.cancelFunc()
	}

	return nil
}
