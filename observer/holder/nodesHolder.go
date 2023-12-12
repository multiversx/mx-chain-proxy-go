package holder

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/availabilityCommon"
)

type cacheType string

const (
	syncedNodesCache            cacheType = "syncedNodes"
	outOfSyncNodesCache         cacheType = "outOfSyncNodes"
	syncedFallbackNodesCache    cacheType = "syncedFallbackNodes"
	outOfSyncFallbackNodesCache cacheType = "outOfSyncFallbackNodes"
)

var (
	log               = logger.GetOrCreate("observer/holder")
	errEmptyNodesList = errors.New("empty nodes list")
)

type nodesHolder struct {
	mut                  sync.RWMutex
	allNodes             map[uint32][]*data.NodeData
	cache                map[string][]*data.NodeData
	availability         data.ObserverDataAvailabilityType
	availabilityProvider availabilityCommon.AvailabilityProvider
}

// NewNodesHolder will return a new instance of a nodesHolder
func NewNodesHolder(syncedNodes []*data.NodeData, fallbackNodes []*data.NodeData, availability data.ObserverDataAvailabilityType) (*nodesHolder, error) {
	if len(syncedNodes) == 0 && len(fallbackNodes) == 0 && availability != data.AvailabilityRecent {
		return nil, errEmptyNodesList
	}
	return &nodesHolder{
		allNodes:             computeInitialNodeList(syncedNodes, fallbackNodes),
		cache:                make(map[string][]*data.NodeData),
		availability:         availability,
		availabilityProvider: availabilityCommon.AvailabilityProvider{},
	}, nil
}

// UpdateNodes will update the internal maps based on the provided nodes
func (nh *nodesHolder) UpdateNodes(nodesWithSyncStatus []*data.NodeData) {
	if len(nodesWithSyncStatus) == 0 {
		return
	}

	nh.mut.Lock()
	defer nh.mut.Unlock()

	nh.allNodes = make(map[uint32][]*data.NodeData)
	nh.cache = make(map[string][]*data.NodeData)
	for _, node := range nodesWithSyncStatus {
		if !nh.availabilityProvider.IsNodeValid(node, nh.availability) {
			continue
		}
		nh.allNodes[node.ShardId] = append(nh.allNodes[node.ShardId], node)
	}

	nh.printNodesInShardsUnprotected()
}

// GetSyncedNodes returns all the synced nodes
func (nh *nodesHolder) GetSyncedNodes(shardID uint32) []*data.NodeData {
	return nh.getObservers(syncedNodesCache, shardID)
}

// GetSyncedFallbackNodes returns all the synced fallback nodes
func (nh *nodesHolder) GetSyncedFallbackNodes(shardID uint32) []*data.NodeData {
	return nh.getObservers(syncedFallbackNodesCache, shardID)
}

// GetOutOfSyncNodes returns all the out of sync nodes
func (nh *nodesHolder) GetOutOfSyncNodes(shardID uint32) []*data.NodeData {
	return nh.getObservers(outOfSyncNodesCache, shardID)
}

// GetOutOfSyncFallbackNodes returns all the out of sync fallback nodes
func (nh *nodesHolder) GetOutOfSyncFallbackNodes(shardID uint32) []*data.NodeData {
	return nh.getObservers(outOfSyncFallbackNodesCache, shardID)
}

// Count computes and returns the total number of nodes
func (nh *nodesHolder) Count() int {
	counter := 0
	nh.mut.RLock()
	defer nh.mut.RUnlock()

	for _, nodes := range nh.allNodes {
		counter += len(nodes)
	}

	return counter
}

func (nh *nodesHolder) getObservers(cache cacheType, shardID uint32) []*data.NodeData {
	cacheKey := getCacheKey(cache, shardID)
	nh.mut.RLock()
	cachedValues, exists := nh.cache[cacheKey]
	nh.mut.RUnlock()

	if exists {
		return cachedValues
	}

	// nodes not cached, compute the list and update the cache
	recomputedList := make([]*data.NodeData, 0)
	nh.mut.Lock()
	defer nh.mut.Unlock()

	cachedValues, exists = nh.cache[cacheKey]
	if exists {
		return cachedValues
	}
	for _, node := range nh.allNodes[shardID] {
		if areCompatibleParameters(cache, node) {
			recomputedList = append(recomputedList, node)
		}
	}
	nh.cache[cacheKey] = recomputedList

	return recomputedList
}

func areCompatibleParameters(cache cacheType, node *data.NodeData) bool {
	isSynced, isFallback := node.IsSynced, node.IsFallback
	if cache == syncedFallbackNodesCache && isSynced && isFallback {
		return true
	}
	if cache == outOfSyncFallbackNodesCache && !isSynced && isFallback {
		return true
	}
	if cache == syncedNodesCache && isSynced && !isFallback {
		return true
	}
	if cache == outOfSyncNodesCache && !isSynced && !isFallback {
		return true
	}

	return false
}

func getCacheKey(cache cacheType, shardID uint32) string {
	return fmt.Sprintf("%s_%d", cache, shardID)
}

// IsInterfaceNil returns true if there is no value under the interface
func (nh *nodesHolder) IsInterfaceNil() bool {
	return nh == nil
}

func (nh *nodesHolder) printNodesInShardsUnprotected() {
	nodesByType := make(map[uint32]map[cacheType][]*data.NodeData)

	// populate nodesByType map
	for shard, nodes := range nh.allNodes {
		if nodesByType[shard] == nil {
			nodesByType[shard] = make(map[cacheType][]*data.NodeData)
		}

		for _, node := range nodes {
			cache := getCacheType(node)
			nodesByType[shard][cache] = append(nodesByType[shard][cache], node)
		}
	}

	printHeader := nh.availabilityProvider.GetDescriptionForAvailability(nh.availability)
	for shard, nodesByCache := range nodesByType {
		log.Info(fmt.Sprintf("shard %d %s", shard, printHeader),
			"synced observers", getNodesListAsString(nodesByCache[syncedNodesCache]),
			"synced fallback observers", getNodesListAsString(nodesByCache[syncedFallbackNodesCache]),
			"out of sync observers", getNodesListAsString(nodesByCache[outOfSyncNodesCache]),
			"out of sync fallback observers", getNodesListAsString(nodesByCache[outOfSyncFallbackNodesCache]))
	}
}

func getCacheType(node *data.NodeData) cacheType {
	if node.IsFallback {
		if node.IsSynced {
			return syncedFallbackNodesCache
		}
		return outOfSyncFallbackNodesCache
	}
	if node.IsSynced {
		return syncedNodesCache
	}
	return outOfSyncNodesCache
}

func getNodesListAsString(nodes []*data.NodeData) string {
	addressesString := ""
	for _, node := range nodes {
		addressesString += fmt.Sprintf("%s, ", node.Address)
	}

	return strings.TrimSuffix(addressesString, ", ")
}

func cloneNodesSlice(input []*data.NodeData) []*data.NodeData {
	clonedSlice := make([]*data.NodeData, len(input))
	for idx, node := range input {
		clonedNodeData := *node
		clonedSlice[idx] = &clonedNodeData
	}

	return clonedSlice
}

func computeInitialNodeList(regularNodes []*data.NodeData, fallbackNodes []*data.NodeData) map[uint32][]*data.NodeData {
	// clone the original maps as not to affect the input
	clonedRegularNodes := cloneNodesSlice(regularNodes)
	clonedFallbackNodes := cloneNodesSlice(fallbackNodes)

	mapToReturn := make(map[uint32][]*data.NodeData)
	// since this function is called at constructor level, consider that all the nodes are active
	for _, node := range clonedRegularNodes {
		node.IsSynced = true
		mapToReturn[node.ShardId] = append(mapToReturn[node.ShardId], node)
	}
	for _, node := range clonedFallbackNodes {
		node.IsSynced = true
		mapToReturn[node.ShardId] = append(mapToReturn[node.ShardId], node)
	}
	return mapToReturn
}
