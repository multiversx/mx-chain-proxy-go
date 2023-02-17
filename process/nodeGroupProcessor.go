package process

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// HeartBeatPath represents the path where an observer exposes his heartbeat status
const HeartBeatPath = "/node/heartbeatstatus"

const systemAccountAddress = "erd1lllllllllllllllllllllllllllllllllllllllllllllllllllsckry7t"

// NodeGroupProcessor is able to process transaction requests
type NodeGroupProcessor struct {
	proc                  Processor
	cacher                HeartbeatCacheHandler
	cacheValidityDuration time.Duration
	cancelFunc            func()
}

// NewNodeGroupProcessor creates a new instance of NodeGroupProcessor
func NewNodeGroupProcessor(
	proc Processor,
	cacher HeartbeatCacheHandler,
	cacheValidityDuration time.Duration,
) (*NodeGroupProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(cacher) {
		return nil, ErrNilHeartbeatCacher
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}
	hbp := &NodeGroupProcessor{
		proc:                  proc,
		cacher:                cacher,
		cacheValidityDuration: cacheValidityDuration,
	}

	return hbp, nil
}

// IsOldStorageForToken returns true if the token is stored in the old fashion
func (hbp *NodeGroupProcessor) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	observers, err := hbp.proc.GetAllObservers()
	if err != nil {
		return false, err
	}

	tokenStorageKey := computeTokenStorageKey(tokenID, nonce)

	for _, observer := range observers {
		if observer.ShardId == core.MetachainShardId {
			continue
		}

		apiResponse := data.AccountKeyValueResponse{}
		apiPath := addressPath + systemAccountAddress + "/key/" + tokenStorageKey
		respCode, err := hbp.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account value for key request",
				"address", systemAccountAddress,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return false, errors.New(apiResponse.Error)
			}

			log.Info("load token from system account", "token", tokenID, "nonce", nonce, "shard ID", observer.ShardId, "value length", len(apiResponse.Data.Value))
			if len(apiResponse.Data.Value) > 0 {
				return false, nil
			}
		} else {
			return false, ErrSendingRequest
		}
	}

	return true, nil
}

func computeTokenStorageKey(tokenID string, nonce uint64) string {
	key := []byte(core.ProtectedKeyPrefix)
	key = append(key, core.ESDTKeyIdentifier...)
	key = append(key, []byte(tokenID)...)

	if nonce > 0 {
		nonceBI := big.NewInt(0).SetUint64(nonce)
		key = append(key, nonceBI.Bytes()...)
	}

	return hex.EncodeToString(key)
}

// GetHeartbeatData will simply forward the heartbeat status from an observer
func (hbp *NodeGroupProcessor) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	heartbeatsToReturn, err := hbp.cacher.LoadHeartbeats()
	if err == nil {
		return heartbeatsToReturn, nil
	}

	log.Info("heartbeat: cannot get from cache. Will fetch from API", "error", err.Error())

	return hbp.getHeartbeatsFromApi()
}

func (hbp *NodeGroupProcessor) getHeartbeatsFromApi() (*data.HeartbeatResponse, error) {
	shardIDs := hbp.proc.GetShardIDs()

	responseMap := make(map[string]data.PubKeyHeartbeat)
	for _, shard := range shardIDs {
		observers, err := hbp.proc.GetObservers(shard)
		if err != nil {
			log.Error("could not get observers", "shard", shard, "error", err.Error())
			continue
		}

		errorsCount := 0
		var response data.HeartbeatApiResponse
		for _, observer := range observers {
			_, err = hbp.proc.CallGetRestEndPoint(observer.Address, HeartBeatPath, &response)
			heartbeats := response.Data.Heartbeats
			if err == nil && len(heartbeats) > 0 {
				hbp.addMessagesToMap(responseMap, heartbeats, shard)
				break
			}

			errorsCount++
			errorMsg := "no heartbeat messages"
			if err != nil {
				errorMsg = err.Error()
			}
			log.Error("heartbeat", "observer", observer.Address, "shard", shard, "error", errorMsg)
		}

		// If no observer responded from a specific shard, log and return error
		if errorsCount == len(observers) {
			log.Error("heartbeat", "error", ErrHeartbeatNotAvailable.Error(), "shard", shard)
			return nil, ErrHeartbeatNotAvailable
		}
	}

	if len(responseMap) == 0 {
		return nil, ErrHeartbeatNotAvailable
	}

	return hbp.mapToResponse(responseMap), nil
}

func (hbp *NodeGroupProcessor) addMessagesToMap(responseMap map[string]data.PubKeyHeartbeat, heartbeats []data.PubKeyHeartbeat, observerShard uint32) {
	for _, heartbeatMessage := range heartbeats {
		isMessageFromCurrentShard := heartbeatMessage.ComputedShardID == observerShard
		isMessageFromShardAfterShuffleOut := heartbeatMessage.ReceivedShardID == observerShard
		belongToCurrentShard := isMessageFromCurrentShard || isMessageFromShardAfterShuffleOut
		if !belongToCurrentShard {
			continue
		}

		oldMessage, found := responseMap[heartbeatMessage.PublicKey]
		if !found {
			responseMap[heartbeatMessage.PublicKey] = heartbeatMessage
			continue // needed because the above get will return a default struct which has IsActive set to false
		}

		if !oldMessage.IsActive && heartbeatMessage.IsActive {
			responseMap[heartbeatMessage.PublicKey] = heartbeatMessage
		}
	}
}

func (hbp *NodeGroupProcessor) mapToResponse(responseMap map[string]data.PubKeyHeartbeat) *data.HeartbeatResponse {
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
func (hbp *NodeGroupProcessor) StartCacheUpdate() {
	if hbp.cancelFunc != nil {
		log.Error("NodeGroupProcessor - cache update already started")
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
				log.Debug("finishing NodeGroupProcessor cache update...")
				return
			}
		}
	}(ctx)
}

func (hbp *NodeGroupProcessor) handleHeartbeatCacheUpdate() {
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
func (hbp *NodeGroupProcessor) Close() error {
	if hbp.cancelFunc != nil {
		hbp.cancelFunc()
	}

	return nil
}
