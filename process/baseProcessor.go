package process

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/core/logger"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin/json"
)

var log = logger.DefaultLogger()

// HttpNoStatusCode is returned when a function which makes a Rest API call fails before it gets to make the
// actual call
const HttpNoStatusCode = 0

// BaseProcessor represents an implementation of CoreProcessor that helps
// processing requests
type BaseProcessor struct {
	addressConverter state.AddressConverter
	lastConfig       *config.Config
	mutState         sync.RWMutex
	shardCoordinator sharding.Coordinator
	observers        map[uint32][]*data.Observer

	httpClient *http.Client
}

// NewBaseProcessor creates a new instance of BaseProcessor struct
func NewBaseProcessor(addressConverter state.AddressConverter, requestTimeoutSec int, shardCoord sharding.Coordinator) (*BaseProcessor, error) {
	if addressConverter == nil {
		return nil, ErrNilAddressConverter
	}
	if shardCoord == nil {
		return nil, ErrNilShardCoordinator
	}
	if requestTimeoutSec <= 0 {
		return nil, ErrInvalidRequestTimeout
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = time.Duration(requestTimeoutSec) * time.Second

	return &BaseProcessor{
		observers:        make(map[uint32][]*data.Observer),
		shardCoordinator: shardCoord,
		httpClient:       httpClient,
		addressConverter: addressConverter,
	}, nil
}

// ApplyConfig applies a config on a base processor
func (bp *BaseProcessor) ApplyConfig(cfg *config.Config) error {
	if cfg == nil {
		return ErrNilConfig
	}
	if len(cfg.Observers) == 0 {
		return ErrEmptyObserversList
	}

	newObservers := make(map[uint32][]*data.Observer)
	for _, observer := range cfg.Observers {
		shardId := observer.ShardId
		newObservers[shardId] = append(newObservers[shardId], observer)
	}

	bp.mutState.Lock()
	bp.observers = newObservers
	bp.mutState.Unlock()

	return nil
}

// GetObservers returns the registered observers on a shard
func (bp *BaseProcessor) GetObservers(shardId uint32) ([]*data.Observer, error) {
	bp.mutState.RLock()
	defer bp.mutState.RUnlock()

	observers := bp.observers[shardId]
	if len(observers) == 0 {
		return nil, ErrMissingObserver
	}

	return observers, nil
}

// GetAllObservers will return all the observers, regardless of shard ID
func (bp *BaseProcessor) GetAllObservers() ([]*data.Observer, error) {
	bp.mutState.RLock()
	defer bp.mutState.RUnlock()

	var observers []*data.Observer
	for _, observersByShard := range bp.observers {
		observers = append(observers, observersByShard...)
	}

	if len(observers) == 0 {
		return nil, ErrNoObserverConnected
	}

	return observers, nil
}

// ComputeShardId computes the shard id in which the account resides
func (bp *BaseProcessor) ComputeShardId(addressBuff []byte) (uint32, error) {
	bp.mutState.RLock()
	defer bp.mutState.RUnlock()

	address, err := bp.addressConverter.CreateAddressFromPublicKeyBytes(addressBuff)
	if err != nil {
		return 0, err
	}

	return bp.shardCoordinator.ComputeId(address), nil
}

// CallGetRestEndPoint calls an external end point (sends a request on a node)
func (bp *BaseProcessor) CallGetRestEndPoint(
	address string,
	path string,
	value interface{},
) error {

	req, err := http.NewRequest("GET", address+path, nil)
	if err != nil {
		return err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Requesting data from nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		log.LogIfError(errNotCritical)
	}()

	return json.NewDecoder(resp.Body).Decode(value)
}

// CallPostRestEndPoint calls an external end point (sends a request on a node)
func (bp *BaseProcessor) CallPostRestEndPoint(
	address string,
	path string,
	data interface{},
	response interface{},
) (int, error) {

	buff, err := json.Marshal(data)
	if err != nil {
		return HttpNoStatusCode, err
	}

	req, err := http.NewRequest("POST", address+path, bytes.NewReader(buff))
	if err != nil {
		return HttpNoStatusCode, err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Posting to nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		return HttpNoStatusCode, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		log.LogIfError(errNotCritical)
	}()

	responseStatusCode := resp.StatusCode
	if responseStatusCode == http.StatusOK { // everything ok, return status ok and the expected response
		return responseStatusCode, json.NewDecoder(resp.Body).Decode(response)
	}

	// status response not ok, return the error
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseStatusCode, err
	}

	return responseStatusCode, errors.New(string(responseBytes))
}
