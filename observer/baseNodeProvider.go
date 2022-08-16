package observer

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type baseNodeProvider struct {
	mutNodes               sync.RWMutex
	nodesMap               map[uint32][]*data.NodeData
	configurationFilePath  string
	syncedNodes            []*data.NodeData
	outOfSyncNodes         []*data.NodeData
	syncedFallbackNodes    []*data.NodeData
	outOfSyncFallbackNodes []*data.NodeData
}

func (bnp *baseNodeProvider) initNodesMaps(nodes []*data.NodeData) error {
	if len(nodes) == 0 {
		return ErrEmptyObserversList
	}

	newNodes := make(map[uint32][]*data.NodeData)
	for _, observer := range nodes {
		shardId := observer.ShardId
		newNodes[shardId] = append(newNodes[shardId], observer)
	}

	bnp.mutNodes.Lock()
	bnp.nodesMap = newNodes
	bnp.syncedNodes, bnp.syncedFallbackNodes = initAllNodesSlice(newNodes)
	bnp.outOfSyncNodes = make([]*data.NodeData, 0)
	bnp.outOfSyncFallbackNodes = make([]*data.NodeData, 0)
	bnp.mutNodes.Unlock()

	return nil
}

// GetAllNodesWithSyncState will return the merged list of active observers and out of sync observers
func (bnp *baseNodeProvider) GetAllNodesWithSyncState() []*data.NodeData {
	bnp.mutNodes.RLock()
	defer bnp.mutNodes.RUnlock()

	nodesSlice := make([]*data.NodeData, 0)
	for _, node := range bnp.syncedNodes {
		nodesSlice = append(nodesSlice, node)
	}
	for _, node := range bnp.outOfSyncNodes {
		nodesSlice = append(nodesSlice, node)
	}

	return nodesSlice
}

// UpdateNodesBasedOnSyncState will handle the nodes lists, by removing out of sync observers or by adding back observers
// that were previously removed because they were out of sync.
func (bnp *baseNodeProvider) UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData) {
	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()

	shardIDs := getSortedSliceIDsSlice(bnp.nodesMap)
	syncedNodes, outOfSyncNodes, err := computeSyncedAndOutOfSyncNodes(nodesWithSyncStatus, shardIDs)
	if err != nil {
		log.Error("cannot update nodes based on sync state", "error", err)
		return
	}

	totalSyncedNodes := len(bnp.syncedNodes) + len(bnp.syncedFallbackNodes)
	if totalSyncedNodes == len(syncedNodes) && len(outOfSyncNodes) == 0 {
		bnp.printSyncedNodesInShardsUnprotected()
		// early exit as all the nodes are in sync
		return
	}

	syncedNodesMap := nodesSliceToShardedMap(syncedNodes)

	bnp.removeOutOfSyncNodesUnprotected(outOfSyncNodes, syncedNodesMap)
	bnp.addSyncedNodesUnprotected(syncedNodes)
	bnp.printSyncedNodesInShardsUnprotected()
}

func (bnp *baseNodeProvider) printSyncedNodesInShardsUnprotected() {
	for shardID, nodes := range bnp.nodesMap {
		inSyncAddresses := make([]string, 0)
		fallbackAddresses := make([]string, 0)
		for _, node := range nodes {
			inSyncAddresses = append(inSyncAddresses, node.Address)
			if node.IsFallback {
				fallbackAddresses = append(fallbackAddresses, node.Address)
			}
		}
		log.Info(fmt.Sprintf("shard %d active nodes", shardID),
			"observers count", len(nodes),
			"addresses", strings.Join(inSyncAddresses, ", "),
			"fallback addresses", strings.Join(fallbackAddresses, ", "))
	}
}

func computeSyncedAndOutOfSyncNodes(nodes []*data.NodeData, shardIDs []uint32) ([]*data.NodeData, []*data.NodeData, error) {
	tempSyncedNodesMap := make(map[uint32][]*data.NodeData)
	tempNotSyncedNodesMap := make(map[uint32][]*data.NodeData)

	for _, node := range nodes {
		if node.IsSynced {
			tempSyncedNodesMap[node.ShardId] = append(tempSyncedNodesMap[node.ShardId], node)
			continue
		}

		tempNotSyncedNodesMap[node.ShardId] = append(tempNotSyncedNodesMap[node.ShardId], node)
	}

	syncedNodes := make([]*data.NodeData, 0)
	notSyncedNodes := make([]*data.NodeData, 0)
	for _, shardID := range shardIDs {
		if len(tempSyncedNodesMap[shardID]) > 0 {
			syncedNodes = append(syncedNodes, tempSyncedNodesMap[shardID]...)
			notSyncedNodes = append(notSyncedNodes, tempNotSyncedNodesMap[shardID]...)
			continue
		}

		if len(tempNotSyncedNodesMap[shardID]) == 0 {
			return nil, nil, fmt.Errorf("%w for shard %d - no synced or not synced node", ErrWrongObserversConfiguration, shardID)
		}

		syncedNodes = append(syncedNodes, tempNotSyncedNodesMap[shardID][0])
		notSyncedNodes = append(notSyncedNodes, tempNotSyncedNodesMap[shardID][1:]...)
	}

	return syncedNodes, notSyncedNodes, nil
}

func (bnp *baseNodeProvider) addSyncedNodesUnprotected(receivedSyncedNodes []*data.NodeData) {
	for _, node := range receivedSyncedNodes {
		if bnp.isReceivedSyncedNodeExistent(node) || bnp.isReceivedSyncedNodeExistentAsFallback(node) {
			continue
		}

		if node.IsFallback {
			bnp.syncedFallbackNodes = append(bnp.syncedFallbackNodes, node)
			continue
		}

		bnp.syncedNodes = append(bnp.syncedNodes, node)
	}

	bnp.nodesMap = nodesSliceToShardedMap(bnp.syncedNodes)
	syncedFallbackNodesMap := nodesSliceToShardedMap(bnp.syncedFallbackNodes)
	for shardId := range syncedFallbackNodesMap {
		bnp.nodesMap[shardId] = append(bnp.nodesMap[shardId], syncedFallbackNodesMap[shardId]...)
	}
}

func (bnp *baseNodeProvider) isReceivedSyncedNodeExistent(receivedNode *data.NodeData) bool {
	for _, node := range bnp.syncedNodes {
		if node.Address == receivedNode.Address && node.ShardId == receivedNode.ShardId {
			return true
		}
	}

	return false
}

func (bnp *baseNodeProvider) isReceivedSyncedNodeExistentAsFallback(receivedNode *data.NodeData) bool {
	for _, node := range bnp.syncedFallbackNodes {
		if node.Address == receivedNode.Address && node.ShardId == receivedNode.ShardId {
			return true
		}
	}

	return false
}

func (bnp *baseNodeProvider) addToOutOfSyncUnprotected(node *data.NodeData) {
	source := &bnp.outOfSyncNodes
	if node.IsFallback {
		source = &bnp.outOfSyncFallbackNodes
	}

	if *source == nil {
		*source = make([]*data.NodeData, 0)
	}

	for _, oosNode := range *source {
		if oosNode.Address == node.Address && oosNode.ShardId == node.ShardId {
			return
		}
	}

	*source = append(*source, node)
}

func (bnp *baseNodeProvider) removeOutOfSyncNodesUnprotected(
	outOfSyncNodes []*data.NodeData,
	syncedNodesMap map[uint32][]*data.NodeData,
) {
	if len(outOfSyncNodes) == 0 {
		bnp.outOfSyncNodes = make([]*data.NodeData, 0)
		bnp.outOfSyncFallbackNodes = make([]*data.NodeData, 0)
		return
	}

	syncedNodes, syncedFallbackNodes := splitSyncedNodes(syncedNodesMap)

	for _, outOfSyncNode := range outOfSyncNodes {
		hasSyncedNodes := len(syncedNodes[outOfSyncNode.ShardId]) >= 1
		if !outOfSyncNode.IsFallback && !hasSyncedNodes {
			hasSyncedFallbackNodes := len(syncedFallbackNodes[outOfSyncNode.ShardId]) >= 1
			if !hasSyncedFallbackNodes {
				log.Warn("cannot remove observer as not enough will remain in shard",
					"address", outOfSyncNode.Address,
					"shard", outOfSyncNode.ShardId)
				continue
			}

			log.Info("not enough nodes in shard, using synced fallback nodes", "shard", outOfSyncNode.ShardId)
			bnp.moveFallbackNodesToSyncedUnprotected(syncedFallbackNodes[outOfSyncNode.ShardId])
		}

		bnp.removeNodeUnprotected(outOfSyncNode)
	}
}

func (bnp *baseNodeProvider) moveFallbackNodesToSyncedUnprotected(syncedFallbackNodes []*data.NodeData) {
	for _, node := range syncedFallbackNodes {
		if bnp.isReceivedSyncedNodeExistent(node) {
			continue
		}

		bnp.syncedNodes = append(bnp.syncedNodes, node)
	}
}

func (bnp *baseNodeProvider) removeNodeUnprotected(node *data.NodeData) {
	bnp.removeNodeFromShardedMapUnprotected(node)
	bnp.removeNodeFromSyncedNodesUnprotected(node)
	bnp.addToOutOfSyncUnprotected(node)
}

func (bnp *baseNodeProvider) removeNodeFromShardedMapUnprotected(node *data.NodeData) {
	nodeIndex := -1
	nodesInShard := bnp.nodesMap[node.ShardId]
	if len(nodesInShard) == 0 {
		log.Error("no observer in shard", "shard ID", node.ShardId)
		return
	}

	for idx, nodeInShard := range nodesInShard {
		if node.ShardId == nodeInShard.ShardId && node.Address == nodeInShard.Address {
			nodeIndex = idx
			break
		}
	}

	if nodeIndex == -1 {
		return
	}

	copy(nodesInShard[nodeIndex:], nodesInShard[nodeIndex+1:])
	nodesInShard[len(nodesInShard)-1] = nil
	nodesInShard = nodesInShard[:len(nodesInShard)-1]

	bnp.nodesMap[node.ShardId] = nodesInShard
	log.Info("updated observers sharded map after removing out of sync observer",
		"address", node.Address,
		"shard ID", node.ShardId,
		"num observers left in shard", len(nodesInShard))
}

func (bnp *baseNodeProvider) removeNodeFromSyncedNodesUnprotected(nodeToRemove *data.NodeData) {
	source := &bnp.syncedNodes
	if nodeToRemove.IsFallback {
		source = &bnp.syncedFallbackNodes
	}

	nodeIndex := -1
	for idx, node := range *source {
		if node.Address == nodeToRemove.Address && node.ShardId == nodeToRemove.ShardId {
			nodeIndex = idx
			break
		}
	}

	if nodeIndex == -1 {
		return
	}

	copy((*source)[nodeIndex:], (*source)[nodeIndex+1:])
	(*source)[len(*source)-1] = nil
	*source = (*source)[:len(*source)-1]
}

// ReloadNodes will reload the observers or the full history observers
func (bnp *baseNodeProvider) ReloadNodes(nodesType data.NodeType) data.NodesReloadResponse {
	bnp.mutNodes.RLock()
	numOldShardsCount := len(bnp.nodesMap)
	bnp.mutNodes.RUnlock()

	newConfig, err := loadMainConfig(bnp.configurationFilePath)
	if err != nil {
		return data.NodesReloadResponse{
			OkRequest:   true,
			Description: "not reloaded",
			Error:       "cannot load configuration file at " + bnp.configurationFilePath,
		}
	}

	nodes := newConfig.Observers
	if nodesType == data.FullHistoryNode {
		nodes = newConfig.FullHistoryNodes
	}

	newNodes := nodesSliceToShardedMap(nodes)
	numNewShardsCount := len(newNodes)

	if numOldShardsCount != numNewShardsCount {
		return data.NodesReloadResponse{
			OkRequest:   false,
			Description: "not reloaded",
			Error:       fmt.Sprintf("different number of shards. before: %d, now: %d", numOldShardsCount, numNewShardsCount),
		}
	}

	bnp.mutNodes.Lock()
	bnp.nodesMap = newNodes
	bnp.syncedNodes, bnp.syncedFallbackNodes = initAllNodesSlice(newNodes)
	bnp.mutNodes.Unlock()

	return data.NodesReloadResponse{
		OkRequest:   true,
		Description: prepareReloadResponseMessage(newNodes),
		Error:       "",
	}
}

func (bnp *baseNodeProvider) getSyncedNodesForShardUnprotected(shardId uint32) ([]*data.NodeData, error) {
	nodesForShard, ok := bnp.nodesMap[shardId]
	if !ok {
		return nil, ErrShardNotAvailable
	}

	syncedNodes := make([]*data.NodeData, 0)
	for _, node := range nodesForShard {
		syncedNodes = append(syncedNodes, node)
	}

	return syncedNodes, nil
}

func loadMainConfig(filepath string) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func nodesSliceToShardedMap(nodes []*data.NodeData) map[uint32][]*data.NodeData {
	newNodes := make(map[uint32][]*data.NodeData)
	for _, observer := range nodes {
		shardId := observer.ShardId
		newNodes[shardId] = append(newNodes[shardId], observer)
	}

	return newNodes
}

func prepareReloadResponseMessage(newNodes map[uint32][]*data.NodeData) string {
	retString := "Reloaded configuration. New configuration: "
	for shardID, nodesInShard := range newNodes {
		retString += fmt.Sprintf("{[Shard %d]:", shardID)
		for _, node := range nodesInShard {
			retString += fmt.Sprintf("[%s]", node.Address)
		}
		retString += "}"
	}

	return retString
}

func initAllNodesSlice(nodesOnShards map[uint32][]*data.NodeData) ([]*data.NodeData, []*data.NodeData) {
	sliceToReturn := make([]*data.NodeData, 0)
	fallbackNodes := make([]*data.NodeData, 0)
	shardIDs := getSortedSliceIDsSlice(nodesOnShards)

	finishedShards := make(map[uint32]struct{})
	for i := 0; ; i++ {
		for _, shardID := range shardIDs {
			if i >= len(nodesOnShards[shardID]) {
				finishedShards[shardID] = struct{}{}
				continue
			}

			node := nodesOnShards[shardID][i]
			if node.IsFallback {
				fallbackNodes = append(fallbackNodes, node)
			} else {
				sliceToReturn = append(sliceToReturn, node)
			}
		}

		if len(finishedShards) == len(nodesOnShards) {
			break
		}
	}

	return sliceToReturn, fallbackNodes
}

func getSortedSliceIDsSlice(nodesOnShards map[uint32][]*data.NodeData) []uint32 {
	shardIDs := make([]uint32, 0)
	for shardID := range nodesOnShards {
		shardIDs = append(shardIDs, shardID)
	}
	sort.SliceStable(shardIDs, func(i, j int) bool {
		return shardIDs[i] < shardIDs[j]
	})

	return shardIDs
}

func splitSyncedNodes(nodes map[uint32][]*data.NodeData) (map[uint32][]*data.NodeData, map[uint32][]*data.NodeData) {
	syncedNodes := make(map[uint32][]*data.NodeData, 0)
	syncedFallbackNodes := make(map[uint32][]*data.NodeData, 0)
	for shardId, nodesInShard := range nodes {
		for _, node := range nodesInShard {
			if node.IsFallback {
				syncedFallbackNodes[shardId] = append(syncedFallbackNodes[shardId], node)
				continue
			}

			syncedNodes[shardId] = append(syncedNodes[shardId], node)
		}
	}

	return syncedNodes, syncedFallbackNodes
}
