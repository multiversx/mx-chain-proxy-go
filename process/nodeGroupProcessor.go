package process

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const (
	// heartbeatPath represents the path where an observer exposes his heartbeat status
	heartbeatPath = "/node/heartbeatstatus"
	// waitingEpochsLeftPath represents the path where an observer the number of epochs left in waiting state for a key
	waitingEpochsLeftPath = "/node/waiting-epochs-left/%s"
	systemAccountAddress  = "erd1lllllllllllllllllllllllllllllllllllllllllllllllllllsckry7t"
)

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
	ngp := &NodeGroupProcessor{
		proc:                  proc,
		cacher:                cacher,
		cacheValidityDuration: cacheValidityDuration,
	}

	return ngp, nil
}

// IsOldStorageForToken returns true if the token is stored in the old fashion
func (ngp *NodeGroupProcessor) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	observers, err := ngp.proc.GetAllObservers(data.AvailabilityRecent)
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
		respCode, err := ngp.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
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
			return false, WrapObserversError(apiResponse.Error)
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
func (ngp *NodeGroupProcessor) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	heartbeatsToReturn, err := ngp.cacher.LoadHeartbeats()
	if err == nil {
		return heartbeatsToReturn, nil
	}

	log.Info("heartbeat: cannot get from cache. Will fetch from API", "error", err.Error())

	return ngp.getHeartbeatsFromApi()
}

func (ngp *NodeGroupProcessor) getHeartbeatsFromApi() (*data.HeartbeatResponse, error) {
	shardIDs := ngp.proc.GetShardIDs()

	responseMap := make(map[string]data.PubKeyHeartbeat)
	for _, shard := range shardIDs {
		observers, err := ngp.proc.GetObservers(shard, data.AvailabilityRecent)
		if err != nil {
			log.Error("could not get observers", "shard", shard, "error", err.Error())
			continue
		}

		errorsCount := 0
		var response data.HeartbeatApiResponse
		for _, observer := range observers {
			_, err = ngp.proc.CallGetRestEndPoint(observer.Address, heartbeatPath, &response)
			heartbeats := response.Data.Heartbeats
			if err == nil && len(heartbeats) > 0 {
				ngp.addMessagesToMap(responseMap, heartbeats, shard)
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

	return ngp.mapToResponse(responseMap), nil
}

func (ngp *NodeGroupProcessor) addMessagesToMap(responseMap map[string]data.PubKeyHeartbeat, heartbeats []data.PubKeyHeartbeat, observerShard uint32) {
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

func (ngp *NodeGroupProcessor) mapToResponse(responseMap map[string]data.PubKeyHeartbeat) *data.HeartbeatResponse {
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
func (ngp *NodeGroupProcessor) StartCacheUpdate() {
	if ngp.cancelFunc != nil {
		log.Error("NodeGroupProcessor - cache update already started")
		return
	}

	var ctx context.Context
	ctx, ngp.cancelFunc = context.WithCancel(context.Background())

	go func(ctx context.Context) {
		timer := time.NewTimer(ngp.cacheValidityDuration)
		defer timer.Stop()

		ngp.handleHeartbeatCacheUpdate()

		for {
			timer.Reset(ngp.cacheValidityDuration)

			select {
			case <-timer.C:
				ngp.handleHeartbeatCacheUpdate()
			case <-ctx.Done():
				log.Debug("finishing NodeGroupProcessor cache update...")
				return
			}
		}
	}(ctx)
}

func (ngp *NodeGroupProcessor) handleHeartbeatCacheUpdate() {
	hbts, err := ngp.getHeartbeatsFromApi()
	if err != nil {
		log.Warn("heartbeat: get from API", "error", err.Error())
	}

	if hbts != nil {
		err = ngp.cacher.StoreHeartbeats(hbts)
		if err != nil {
			log.Warn("heartbeat: store in cache", "error", err.Error())
		}
	}
}

// GetWaitingEpochsLeftForPublicKey returns the number of epochs left for the public key until it becomes eligible
func (ngp *NodeGroupProcessor) GetWaitingEpochsLeftForPublicKey(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
	if len(publicKey) == 0 {
		return nil, ErrEmptyPubKey
	}

	observers, err := ngp.proc.GetAllObservers(data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	var lastErr error
	responseWaitingEpochsLeft := data.WaitingEpochsLeftApiResponse{}
	path := fmt.Sprintf(waitingEpochsLeftPath, publicKey)
	for _, observer := range observers {
		_, lastErr = ngp.proc.CallGetRestEndPoint(observer.Address, path, &responseWaitingEpochsLeft)
		if lastErr != nil {
			log.Error("waiting epochs left request", "observer", observer.Address, "public key", publicKey, "error", lastErr.Error())
			continue
		}

		log.Info("waiting epochs left request", "shard ID", observer.ShardId, "observer", observer.Address, "public key", publicKey)
		return &responseWaitingEpochsLeft, nil

	}

	return nil, WrapObserversError(responseWaitingEpochsLeft.Error)
}

// Close will handle the closing of the cache update go routine
func (ngp *NodeGroupProcessor) Close() error {
	if ngp.cancelFunc != nil {
		ngp.cancelFunc()
	}

	return nil
}
