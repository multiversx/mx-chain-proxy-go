package observer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type baseNodeProvider struct {
	mutNodes              sync.RWMutex
	nodesMap              map[uint32][]*data.NodeData
	configurationFilePath string
	syncedNodes           []*data.NodeData
	outOfSyncNodes        []*data.NodeData
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
	bnp.syncedNodes = initAllNodesSlice(newNodes)
	bnp.outOfSyncNodes = make([]*data.NodeData, 0)
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

	syncedNodesMap := nodesSliceToShardedMap(syncedNodes)

	if len(bnp.syncedNodes) == len(syncedNodes) && len(outOfSyncNodes) == 0 {
		// early exit as all the nodes are in sync
		return
	}

	for _, outOfSyncNode := range outOfSyncNodes {
		if len(syncedNodesMap[outOfSyncNode.ShardId]) < 1 {
			log.Warn("cannot remove observer as not enough will remain in shard",
				"address", outOfSyncNode.Address,
				"shard", outOfSyncNode.ShardId)
			continue
		}

		bnp.removeNodeUnprotected(outOfSyncNode)
	}

	bnp.addSyncedNodesUnprotected(syncedNodes)
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
		if bnp.isReceivedSyncedNodeExistent(node) {
			continue
		}

		bnp.syncedNodes = append(bnp.syncedNodes, node)
	}

	bnp.nodesMap = nodesSliceToShardedMap(bnp.syncedNodes)
}

func (bnp *baseNodeProvider) isReceivedSyncedNodeExistent(receivedNode *data.NodeData) bool {
	for _, node := range bnp.syncedNodes {
		if node.Address == receivedNode.Address && node.ShardId == receivedNode.ShardId {
			return true
		}
	}

	return false
}

func (bnp *baseNodeProvider) addToOutOfSyncUnprotected(node *data.NodeData) {
	if bnp.outOfSyncNodes == nil {
		bnp.outOfSyncNodes = make([]*data.NodeData, 0)
	}

	for _, oosNode := range bnp.outOfSyncNodes {
		if oosNode.Address == node.Address && oosNode.ShardId == node.ShardId {
			log.Warn("[programming error] -> node is already in out of sync list", "address", node.Address, "shard ID", node.ShardId)
			return
		}
	}

	bnp.outOfSyncNodes = append(bnp.outOfSyncNodes, node)
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
		log.Warn("out of sync observer to remove from sharded map not found", "address", node.Address, "shard ID", node.ShardId)
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
	nodeIndex := -1
	for idx, node := range bnp.syncedNodes {
		if node.Address == nodeToRemove.Address && node.ShardId == nodeToRemove.ShardId {
			nodeIndex = idx
			break
		}
	}

	if nodeIndex == -1 {
		log.Warn("out of sync observer to remove from all nodes not found", "address", nodeToRemove.Address, "shard ID", nodeToRemove.ShardId)
		return
	}

	copy(bnp.syncedNodes[nodeIndex:], bnp.syncedNodes[nodeIndex+1:])
	bnp.syncedNodes[len(bnp.syncedNodes)-1] = nil
	bnp.syncedNodes = bnp.syncedNodes[:len(bnp.syncedNodes)-1]
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
	bnp.syncedNodes = initAllNodesSlice(newNodes)
	bnp.mutNodes.Unlock()

	return data.NodesReloadResponse{
		OkRequest:   true,
		Description: prepareReloadResponseMessage(newNodes),
		Error:       "",
	}
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

func initAllNodesSlice(nodesOnShards map[uint32][]*data.NodeData) []*data.NodeData {
	sliceToReturn := make([]*data.NodeData, 0)
	shardIDs := getSortedSliceIDsSlice(nodesOnShards)

	finishedShards := make(map[uint32]struct{})
	for i := 0; ; i++ {
		for _, shardID := range shardIDs {
			if i >= len(nodesOnShards[shardID]) {
				finishedShards[shardID] = struct{}{}
				continue
			}

			sliceToReturn = append(sliceToReturn, nodesOnShards[shardID][i])
		}

		if len(finishedShards) == len(nodesOnShards) {
			break
		}
	}

	return sliceToReturn
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
