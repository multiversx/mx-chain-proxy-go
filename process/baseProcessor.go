package process

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/sharding"
	proxyData "github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/observer"
)

var log = logger.GetOrCreate("process")
var mutHttpClient sync.RWMutex

// BaseProcessor represents an implementation of CoreProcessor that helps
// processing requests
type BaseProcessor struct {
	mutState                 sync.RWMutex
	shardCoordinator         sharding.Coordinator
	observersProvider        observer.NodesProviderHandler
	fullHistoryNodesProvider observer.NodesProviderHandler
	pubKeyConverter          core.PubkeyConverter
	shardIDs                 []uint32

	httpClient *http.Client
}

// NewBaseProcessor creates a new instance of BaseProcessor struct
func NewBaseProcessor(
	requestTimeoutSec int,
	shardCoord sharding.Coordinator,
	observersProvider observer.NodesProviderHandler,
	fullHistoryNodesProvider observer.NodesProviderHandler,
	pubKeyConverter core.PubkeyConverter,
) (*BaseProcessor, error) {
	if check.IfNil(shardCoord) {
		return nil, ErrNilShardCoordinator
	}
	if requestTimeoutSec <= 0 {
		return nil, ErrInvalidRequestTimeout
	}
	if check.IfNil(observersProvider) {
		return nil, fmt.Errorf("%w for observers", ErrNilNodesProvider)
	}
	if check.IfNil(fullHistoryNodesProvider) {
		return nil, fmt.Errorf("%w for full history nodes", ErrNilNodesProvider)
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	httpClient := http.DefaultClient
	mutHttpClient.Lock()
	httpClient.Timeout = time.Duration(requestTimeoutSec) * time.Second
	mutHttpClient.Unlock()

	return &BaseProcessor{
		shardCoordinator:         shardCoord,
		observersProvider:        observersProvider,
		fullHistoryNodesProvider: fullHistoryNodesProvider,
		httpClient:               httpClient,
		pubKeyConverter:          pubKeyConverter,
		shardIDs:                 computeShardIDs(shardCoord),
	}, nil
}

// GetShardIDs will return the shard IDs slice
func (bp *BaseProcessor) GetShardIDs() []uint32 {
	return bp.shardIDs
}

// ReloadObservers will call the nodes reloading from the observers provider
func (bp *BaseProcessor) ReloadObservers() proxyData.NodesReloadResponse {
	return bp.observersProvider.ReloadNodes(proxyData.Observer)
}

// ReloadFullHistoryObservers will call the nodes reloading from the full history observers provider
func (bp *BaseProcessor) ReloadFullHistoryObservers() proxyData.NodesReloadResponse {
	return bp.fullHistoryNodesProvider.ReloadNodes(proxyData.FullHistoryNode)
}

// GetObservers returns the registered observers on a shard
func (bp *BaseProcessor) GetObservers(shardID uint32) ([]*proxyData.NodeData, error) {
	return bp.observersProvider.GetNodesByShardId(shardID)
}

// GetAllObservers will return all the observers, regardless of shard ID
func (bp *BaseProcessor) GetAllObservers() ([]*proxyData.NodeData, error) {
	return bp.observersProvider.GetAllNodes()
}

// GetObserversOnePerShard will return a slice containing an observer for each shard
func (bp *BaseProcessor) GetObserversOnePerShard() ([]*proxyData.NodeData, error) {
	return bp.getNodesOnePerShard(bp.observersProvider.GetNodesByShardId)
}

// GetFullHistoryNodes returns the registered full history nodes on a shard
func (bp *BaseProcessor) GetFullHistoryNodes(shardID uint32) ([]*proxyData.NodeData, error) {
	return bp.fullHistoryNodesProvider.GetNodesByShardId(shardID)
}

// GetAllFullHistoryNodes will return all the full history nodes, regardless of shard ID
func (bp *BaseProcessor) GetAllFullHistoryNodes() ([]*proxyData.NodeData, error) {
	return bp.fullHistoryNodesProvider.GetAllNodes()
}

// GetFullHistoryNodesOnePerShard will return a slice containing a full history node for each shard
func (bp *BaseProcessor) GetFullHistoryNodesOnePerShard() ([]*proxyData.NodeData, error) {
	return bp.getNodesOnePerShard(bp.fullHistoryNodesProvider.GetNodesByShardId)
}

func (bp *BaseProcessor) getNodesOnePerShard(
	observersInShardGetter func(shardID uint32) ([]*proxyData.NodeData, error),
) ([]*proxyData.NodeData, error) {
	numShards := bp.shardCoordinator.NumberOfShards()
	sliceToReturn := make([]*proxyData.NodeData, 0)

	for shardID := uint32(0); shardID < numShards; shardID++ {
		observersInShard, err := observersInShardGetter(shardID)
		if err != nil || len(observersInShard) < 1 {
			continue
		}

		sliceToReturn = append(sliceToReturn, observersInShard[0])
	}

	observersInShardMeta, err := observersInShardGetter(core.MetachainShardId)
	if err == nil && len(observersInShardMeta) > 0 {
		sliceToReturn = append(sliceToReturn, observersInShardMeta[0])
	}

	if len(sliceToReturn) == 0 {
		return nil, ErrNoObserverAvailable
	}

	return sliceToReturn, nil
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
) (int, error) {

	req, err := http.NewRequest("GET", address+path, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Requesting data from nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return http.StatusRequestTimeout, err
		}

		return http.StatusNotFound, err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		if errNotCritical != nil {
			log.Warn("base process GET: close body", "error", errNotCritical.Error())
		}
	}()

	err = json.NewDecoder(resp.Body).Decode(value)
	if err != nil {
		return 0, err
	}

	responseStatusCode := resp.StatusCode
	if responseStatusCode == http.StatusOK { // everything ok, return status ok and the expected response
		return responseStatusCode, nil
	}

	// status response not ok, return the error
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseStatusCode, err
	}

	return responseStatusCode, errors.New(string(responseBytes))
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

		return http.StatusNotFound, err
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

	genericApiResponse := proxyData.GenericAPIResponse{}
	err = json.Unmarshal(responseBytes, &genericApiResponse)
	if err != nil {
		return responseStatusCode, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return responseStatusCode, errors.New(genericApiResponse.Error)
}

func isTimeoutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}

// GetShardCoordinator returns the shard coordinator
func (bp *BaseProcessor) GetShardCoordinator() sharding.Coordinator {
	return bp.shardCoordinator
}

// GetPubKeyConverter returns the public key converter
func (bp *BaseProcessor) GetPubKeyConverter() core.PubkeyConverter {
	return bp.pubKeyConverter
}

// GetObserverProvider returns the observers provider
func (bp *BaseProcessor) GetObserverProvider() observer.NodesProviderHandler {
	return bp.observersProvider
}

// GetFullHistoryNodesProvider returns the full history nodes provider object
func (bp *BaseProcessor) GetFullHistoryNodesProvider() observer.NodesProviderHandler {
	return bp.fullHistoryNodesProvider
}

func computeShardIDs(shardCoordinator sharding.Coordinator) []uint32 {
	shardIDs := make([]uint32, 0)
	for i := uint32(0); i < shardCoordinator.NumberOfShards(); i++ {
		shardIDs = append(shardIDs, i)
	}

	shardIDs = append(shardIDs, core.MetachainShardId)

	return shardIDs
}

// IsInterfaceNil returns true if there is no value under the interface
func (bp *BaseProcessor) IsInterfaceNil() bool {
	return bp == nil
}
