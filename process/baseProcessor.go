package process

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
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

const (
	nodeSyncedNonceDifferenceThreshold = 10
	stepDelayForCheckingNodesSyncState = 1 * time.Minute
	timeoutDurationForNodeStatus       = 2 * time.Second
)

// BaseProcessor represents an implementation of CoreProcessor that helps
// processing requests
type BaseProcessor struct {
	mutState                       sync.RWMutex
	shardCoordinator               sharding.Coordinator
	observersProvider              observer.NodesProviderHandler
	fullHistoryNodesProvider       observer.NodesProviderHandler
	pubKeyConverter                core.PubkeyConverter
	shardIDs                       []uint32
	delayForCheckingNodesSyncState time.Duration
	cancelFunc                     func()

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

	bp := &BaseProcessor{
		shardCoordinator:               shardCoord,
		observersProvider:              observersProvider,
		fullHistoryNodesProvider:       fullHistoryNodesProvider,
		httpClient:                     httpClient,
		pubKeyConverter:                pubKeyConverter,
		shardIDs:                       computeShardIDs(shardCoord),
		delayForCheckingNodesSyncState: stepDelayForCheckingNodesSyncState,
	}

	return bp, nil
}

// StartNodesSyncStateChecks will simply start the goroutine that handles the nodes sync state
func (bp *BaseProcessor) StartNodesSyncStateChecks() {
	if bp.cancelFunc != nil {
		log.Error("BaseProcessor - cache update already started")
		return
	}

	var ctx context.Context
	ctx, bp.cancelFunc = context.WithCancel(context.Background())

	go bp.handleOutOfSyncNodes(ctx)
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

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal(responseBodyBytes, value)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	responseStatusCode := resp.StatusCode
	if responseStatusCode == http.StatusOK { // everything ok, return status ok and the expected response
		return responseStatusCode, nil
	}

	// status response not ok, return the error
	return responseStatusCode, errors.New(string(responseBodyBytes))
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

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	responseStatusCode := resp.StatusCode
	if responseStatusCode == http.StatusOK { // everything ok, return status ok and the expected response
		return responseStatusCode, json.Unmarshal(responseBodyBytes, response)
	}

	// status response not ok, return the error
	genericApiResponse := proxyData.GenericAPIResponse{}
	err = json.Unmarshal(responseBodyBytes, &genericApiResponse)
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

func (bp *BaseProcessor) handleOutOfSyncNodes(ctx context.Context) {
	timer := time.NewTimer(bp.delayForCheckingNodesSyncState)
	defer timer.Stop()

	bp.updateNodesWithSync()
	for {
		timer.Reset(bp.delayForCheckingNodesSyncState)

		select {
		case <-timer.C:
			bp.updateNodesWithSync()
		case <-ctx.Done():
			log.Debug("finishing BaseProcessor nodes state update...")
			return
		}
	}
}

func (bp *BaseProcessor) updateNodesWithSync() {
	observers := bp.observersProvider.GetAllNodesWithSyncState()
	observersWithSyncStatus := bp.getNodesWithSyncStatus(observers)
	bp.observersProvider.UpdateNodesBasedOnSyncState(observersWithSyncStatus)

	fullHistoryNodes := bp.fullHistoryNodesProvider.GetAllNodesWithSyncState()
	fullHistoryNodesWithSyncStatus := bp.getNodesWithSyncStatus(fullHistoryNodes)
	bp.fullHistoryNodesProvider.UpdateNodesBasedOnSyncState(fullHistoryNodesWithSyncStatus)
}

func (bp *BaseProcessor) getNodesWithSyncStatus(nodes []*proxyData.NodeData) []*proxyData.NodeData {
	nodesToReturn := make([]*proxyData.NodeData, 0)
	for _, node := range nodes {
		outOfSync, err := bp.isNodeOutOfSync(node)
		if err != nil {
			log.Warn("cannot get node status. will mark as inactive", "address", node.Address, "error", err)
			outOfSync = true
		}

		node.IsSynced = !outOfSync
		nodesToReturn = append(nodesToReturn, node)
	}

	return nodesToReturn
}

func (bp *BaseProcessor) isNodeOutOfSync(node *proxyData.NodeData) (bool, error) {
	var nodeStatusResponse proxyData.NodeStatusAPIResponse

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDurationForNodeStatus)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, node.Address+"/node/status", nil)
	if err != nil {
		return false, err
	}

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer func() {
		if resp != nil && resp.Body != nil {
			log.LogIfError(resp.Body.Close())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("observer %s responded with code %d", node.Address, resp.StatusCode)
	}

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(responseBodyBytes, &nodeStatusResponse)
	if err != nil {
		return false, err
	}

	nonce := nodeStatusResponse.Data.Metrics.Nonce
	probableHighestNonce := nodeStatusResponse.Data.Metrics.ProbableHighestNonce
	isReadyForVMQueries := isNodeReadyForVMQueries(nodeStatusResponse.Data.Metrics.AreVmQueriesReady)

	probableHighestNonceLessThanOrEqualToNonce := probableHighestNonce <= nonce
	nonceDifferenceBeyondThreshold := probableHighestNonce-nonce > nodeSyncedNonceDifferenceThreshold
	isNodeOutOfSync := !probableHighestNonceLessThanOrEqualToNonce && nonceDifferenceBeyondThreshold

	log.Info("node status",
		"address", node.Address,
		"shard", node.ShardId,
		"nonce", nonce,
		"probable highest nonce", probableHighestNonce,
		"is synced", !isNodeOutOfSync,
		"is ready for VM Queries", isReadyForVMQueries)

	if !isReadyForVMQueries {
		isNodeOutOfSync = true
	}

	return isNodeOutOfSync, nil
}

func isNodeReadyForVMQueries(metricValue string) bool {
	if strconv.FormatBool(true) == metricValue {
		return true
	}

	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (bp *BaseProcessor) IsInterfaceNil() bool {
	return bp == nil
}

// Close will handle the closing of the cache update go routine
func (bp *BaseProcessor) Close() error {
	if bp.cancelFunc != nil {
		bp.cancelFunc()
	}

	return nil
}
