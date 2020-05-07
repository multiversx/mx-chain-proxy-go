package process

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/observer"
	"github.com/gin-gonic/gin/json"
)

var log = logger.GetOrCreate("process")
var mutHttpClient sync.RWMutex

// BaseProcessor represents an implementation of CoreProcessor that helps
// processing requests
type BaseProcessor struct {
	lastConfig        *config.Config
	mutState          sync.RWMutex
	shardCoordinator  sharding.Coordinator
	observersProvider observer.ObserversProviderHandler
	pubKeyConverter   state.PubkeyConverter

	httpClient *http.Client
}

// NewBaseProcessor creates a new instance of BaseProcessor struct
func NewBaseProcessor(
	requestTimeoutSec int,
	shardCoord sharding.Coordinator,
	observersProvider observer.ObserversProviderHandler,
	pubKeyConverter state.PubkeyConverter,
) (*BaseProcessor, error) {
	if check.IfNil(shardCoord) {
		return nil, ErrNilShardCoordinator
	}
	if requestTimeoutSec <= 0 {
		return nil, ErrInvalidRequestTimeout
	}
	if check.IfNil(observersProvider) {
		return nil, ErrNilObserversProvider
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	httpClient := http.DefaultClient
	mutHttpClient.Lock()
	httpClient.Timeout = time.Duration(requestTimeoutSec) * time.Second
	mutHttpClient.Unlock()

	return &BaseProcessor{
		shardCoordinator:  shardCoord,
		observersProvider: observersProvider,
		httpClient:        httpClient,
		pubKeyConverter:   pubKeyConverter,
	}, nil
}

// GetObservers returns the registered observers on a shard
func (bp *BaseProcessor) GetObservers(shardId uint32) ([]*data.Observer, error) {
	return bp.observersProvider.GetObserversByShardId(shardId)
}

// GetAllObservers will return all the observers, regardless of shard ID
func (bp *BaseProcessor) GetAllObservers() []*data.Observer {
	return bp.observersProvider.GetAllObservers()
}

// ComputeShardId computes the shard id in which the account resides
func (bp *BaseProcessor) ComputeShardId(addressBuff []byte) (uint32, error) {
	bp.mutState.RLock()
	defer bp.mutState.RUnlock()

	return bp.shardCoordinator.ComputeId(addressBuff), nil
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
		if errNotCritical != nil {
			log.Warn("base process GET: close body", "error", errNotCritical.Error())
		}
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
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", address+path, bytes.NewReader(buff))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Posting to nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return http.StatusRequestTimeout, err
		}

		return http.StatusBadRequest, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		if errNotCritical != nil {
			log.Warn("base process POST: close body", "error", errNotCritical.Error())
		}
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

func isTimeoutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (bp *BaseProcessor) IsInterfaceNil() bool {
	return bp == nil
}
