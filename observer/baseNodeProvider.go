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
	shardIds               []uint32
	configurationFilePath  string
	syncedNodes            []*data.NodeData
	outOfSyncNodes         []*data.NodeData
	syncedFallbackNodes    []*data.NodeData
	outOfSyncFallbackNodes []*data.NodeData
	lastSyncedNodes        map[uint32]*data.NodeData
}

func (bnp *baseNodeProvider) initNodes(nodes []*data.NodeData) error {
	if len(nodes) == 0 {
		return ErrEmptyObserversList
	}

	newNodes := make(map[uint32][]*data.NodeData)
	for _, observer := range nodes {
		shardId := observer.ShardId
		newNodes[shardId] = append(newNodes[shardId], observer)
	}

	bnp.mutNodes.Lock()
	bnp.shardIds = getSortedShardIDsSlice(newNodes)
	bnp.syncedNodes, bnp.syncedFallbackNodes = initAllNodesSlice(newNodes)
	bnp.outOfSyncNodes = make([]*data.NodeData, 0)
	bnp.outOfSyncFallbackNodes = make([]*data.NodeData, 0)
	bnp.lastSyncedNodes = make(map[uint32]*data.NodeData)
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
	for _, node := range bnp.syncedFallbackNodes {
		nodesSlice = append(nodesSlice, node)
	}
	for _, node := range bnp.outOfSyncFallbackNodes {
		nodesSlice = append(nodesSlice, node)
	}

	return nodesSlice
}

// UpdateNodesBasedOnSyncState will handle the nodes lists, by removing out of sync observers or by adding back observers
// that were previously removed because they were out of sync.
func (bnp *baseNodeProvider) UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData) {
	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()

	syncedNodes, syncedFallbackNodes, outOfSyncNodes, err := computeSyncedAndOutOfSyncNodes(nodesWithSyncStatus, bnp.shardIds)
	if err != nil {
		log.Error("cannot update nodes based on sync state", "error", err)
		return
	}

	sameNumOfSynced := len(bnp.syncedNodes) == len(syncedNodes)
	sameNumOfSyncedFallback := len(bnp.syncedFallbackNodes) == len(syncedFallbackNodes)
	if sameNumOfSynced && sameNumOfSyncedFallback && len(outOfSyncNodes) == 0 {
		bnp.printSyncedNodesInShardsUnprotected()
		// early exit as all the nodes are in sync
		return
	}

	syncedNodesMap := nodesSliceToShardedMap(syncedNodes)
	syncedFallbackNodesMap := nodesSliceToShardedMap(syncedFallbackNodes)

	bnp.removeOutOfSyncNodesUnprotected(outOfSyncNodes, syncedNodesMap, syncedFallbackNodesMap)
	bnp.addSyncedNodesUnprotected(syncedNodes, syncedFallbackNodes)
	bnp.printSyncedNodesInShardsUnprotected()
}

func (bnp *baseNodeProvider) printSyncedNodesInShardsUnprotected() {
	inSyncAddresses := make(map[uint32][]string, 0)
	for _, syncedNode := range bnp.syncedNodes {
		inSyncAddresses[syncedNode.ShardId] = append(inSyncAddresses[syncedNode.ShardId], syncedNode.Address)
	}

	inSyncFallbackAddresses := make(map[uint32][]string, 0)
	for _, syncedFallbackNode := range bnp.syncedFallbackNodes {
		inSyncFallbackAddresses[syncedFallbackNode.ShardId] = append(inSyncFallbackAddresses[syncedFallbackNode.ShardId], syncedFallbackNode.Address)
	}

	for _, shardID := range bnp.shardIds {
		totalNumOfActiveNodes := len(inSyncAddresses[shardID]) + len(inSyncFallbackAddresses[shardID])
		// if none of them is active, use the backup if exists
		hasBackup := bnp.lastSyncedNodes[shardID] != nil
		if totalNumOfActiveNodes == 0 && hasBackup {
			totalNumOfActiveNodes++
			inSyncAddresses[shardID] = append(inSyncAddresses[shardID], bnp.lastSyncedNodes[shardID].Address)
		}
		log.Info(fmt.Sprintf("shard %d active nodes", shardID),
			"observers count", totalNumOfActiveNodes,
			"addresses", strings.Join(inSyncAddresses[shardID], ", "),
			"fallback addresses", strings.Join(inSyncFallbackAddresses[shardID], ", "))
	}
}

func computeSyncedAndOutOfSyncNodes(nodes []*data.NodeData, shardIDs []uint32) ([]*data.NodeData, []*data.NodeData, []*data.NodeData, error) {
	tempSyncedNodesMap := make(map[uint32][]*data.NodeData)
	tempSyncedFallbackNodesMap := make(map[uint32][]*data.NodeData)
	tempNotSyncedNodesMap := make(map[uint32][]*data.NodeData)

	for _, node := range nodes {
		if node.IsSynced {
			if node.IsFallback {
				tempSyncedFallbackNodesMap[node.ShardId] = append(tempSyncedFallbackNodesMap[node.ShardId], node)
			} else {
				tempSyncedNodesMap[node.ShardId] = append(tempSyncedNodesMap[node.ShardId], node)
			}
			continue
		}

		tempNotSyncedNodesMap[node.ShardId] = append(tempNotSyncedNodesMap[node.ShardId], node)
	}

	syncedNodes := make([]*data.NodeData, 0)
	syncedFallbackNodes := make([]*data.NodeData, 0)
	notSyncedNodes := make([]*data.NodeData, 0)
	for _, shardID := range shardIDs {
		syncedNodes = append(syncedNodes, tempSyncedNodesMap[shardID]...)
		syncedFallbackNodes = append(syncedFallbackNodes, tempSyncedFallbackNodesMap[shardID]...)
		notSyncedNodes = append(notSyncedNodes, tempNotSyncedNodesMap[shardID]...)

		totalLen := len(tempSyncedNodesMap[shardID]) + len(tempSyncedFallbackNodesMap[shardID]) + len(tempNotSyncedNodesMap[shardID])
		if totalLen == 0 {
			return nil, nil, nil, fmt.Errorf("%w for shard %d - no synced or not synced node", ErrWrongObserversConfiguration, shardID)
		}
	}

	return syncedNodes, syncedFallbackNodes, notSyncedNodes, nil
}

func (bnp *baseNodeProvider) addSyncedNodesUnprotected(receivedSyncedNodes []*data.NodeData, receivedSyncedFallbackNodes []*data.NodeData) {
	syncedNodesPerShard := make(map[uint32][]string)
	for _, node := range receivedSyncedNodes {
		bnp.removeFromOutOfSyncIfNeededUnprotected(node)
		syncedNodesPerShard[node.ShardId] = append(syncedNodesPerShard[node.ShardId], node.Address)
		if bnp.isReceivedSyncedNodeExistent(node) {
			continue
		}

		bnp.syncedNodes = append(bnp.syncedNodes, node)
	}

	for _, node := range receivedSyncedFallbackNodes {
		bnp.removeFromOutOfSyncIfNeededUnprotected(node)
		if bnp.isReceivedSyncedNodeExistentAsFallback(node) {
			continue
		}

		bnp.syncedFallbackNodes = append(bnp.syncedFallbackNodes, node)
	}

	// if there are at least one synced node regular received, clean the backup list
	for _, shardId := range bnp.shardIds {
		if len(syncedNodesPerShard[shardId]) != 0 {
			delete(bnp.lastSyncedNodes, shardId)
		}
	}
}

func (bnp *baseNodeProvider) removeFromOutOfSyncIfNeededUnprotected(node *data.NodeData) {
	source := &bnp.outOfSyncNodes
	if node.IsFallback {
		source = &bnp.outOfSyncFallbackNodes
	}

	if *source == nil {
		return
	}

	bnp.removeNodeFromListUnprotected(node, source)
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
	syncedFallbackNodesMap map[uint32][]*data.NodeData,
) {
	if len(outOfSyncNodes) == 0 {
		bnp.outOfSyncNodes = make([]*data.NodeData, 0)
		bnp.outOfSyncFallbackNodes = make([]*data.NodeData, 0)
		return
	}

	for _, outOfSyncNode := range outOfSyncNodes {
		hasOneSyncedNode := len(syncedNodesMap[outOfSyncNode.ShardId]) >= 1
		hasEnoughSyncedFallbackNodes := len(syncedFallbackNodesMap[outOfSyncNode.ShardId]) > 1
		canDeleteFallbackNode := hasOneSyncedNode || hasEnoughSyncedFallbackNodes
		if outOfSyncNode.IsFallback && canDeleteFallbackNode {
			bnp.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// if trying to delete last fallback, use last known synced node
		// if backup node does not exist, keep fallback
		hasBackup := bnp.lastSyncedNodes[outOfSyncNode.ShardId] != nil
		if outOfSyncNode.IsFallback && hasBackup {
			bnp.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		hasEnoughSyncedNodes := len(syncedNodesMap[outOfSyncNode.ShardId]) >= 1
		if hasEnoughSyncedNodes {
			bnp.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// trying to remove last synced node
		// if fallbacks are available, save this one as backup and use fallbacks
		// else, keep using this one
		// save this last regular observer as backup in case fallbacks go offline
		// also, if this is the old fallback observer which didn't get synced, keep it in list
		wasSyncedAtPreviousStep := bnp.isReceivedSyncedNodeExistent(outOfSyncNode)
		isBackupObserver := bnp.lastSyncedNodes[outOfSyncNode.ShardId] == outOfSyncNode
		if !outOfSyncNode.IsFallback && wasSyncedAtPreviousStep || isBackupObserver {
			log.Info("backup observer updated",
				"address", outOfSyncNode.Address,
				"is fallback", outOfSyncNode.IsFallback,
				"shard", outOfSyncNode.ShardId)
			bnp.lastSyncedNodes[outOfSyncNode.ShardId] = outOfSyncNode
		}
		hasOneSyncedFallbackNode := len(syncedFallbackNodesMap[outOfSyncNode.ShardId]) >= 1
		if hasOneSyncedFallbackNode {
			bnp.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// safe to delete regular observer, as it is already in lastSyncedNodes map
		if !outOfSyncNode.IsFallback {
			bnp.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// this is a fallback node, with no synced nodes.
		// save it as backup and delete it from its list
		bnp.lastSyncedNodes[outOfSyncNode.ShardId] = outOfSyncNode
		bnp.removeNodeUnprotected(outOfSyncNode)
	}
}

func (bnp *baseNodeProvider) removeNodeUnprotected(node *data.NodeData) {
	bnp.removeNodeFromSyncedNodesUnprotected(node)
	bnp.addToOutOfSyncUnprotected(node)
}

func (bnp *baseNodeProvider) removeNodeFromSyncedNodesUnprotected(nodeToRemove *data.NodeData) {
	source := &bnp.syncedNodes
	if nodeToRemove.IsFallback {
		source = &bnp.syncedFallbackNodes
	}

	bnp.removeNodeFromListUnprotected(nodeToRemove, source)
}

func (bnp *baseNodeProvider) removeNodeFromListUnprotected(nodeToRemove *data.NodeData, source *[]*data.NodeData) {
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
	numOldShardsCount := len(bnp.shardIds)
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
	bnp.shardIds = getSortedShardIDsSlice(newNodes)
	bnp.syncedNodes, bnp.syncedFallbackNodes = initAllNodesSlice(newNodes)
	bnp.lastSyncedNodes = make(map[uint32]*data.NodeData)
	bnp.mutNodes.Unlock()

	return data.NodesReloadResponse{
		OkRequest:   true,
		Description: prepareReloadResponseMessage(newNodes),
		Error:       "",
	}
}

func (bnp *baseNodeProvider) getSyncedNodesForShardUnprotected(shardId uint32) ([]*data.NodeData, error) {
	syncedNodes := make([]*data.NodeData, 0)
	for _, node := range bnp.syncedNodes {
		if node.ShardId == shardId {
			syncedNodes = append(syncedNodes, node)
		}
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	for _, node := range bnp.syncedFallbackNodes {
		if node.ShardId == shardId {
			syncedNodes = append(syncedNodes, node)
		}
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	backupNode, hasBackup := bnp.lastSyncedNodes[shardId]
	if hasBackup {
		return []*data.NodeData{backupNode}, nil
	}

	return nil, ErrShardNotAvailable
}

func (bnp *baseNodeProvider) getSyncedNodesUnprotected() ([]*data.NodeData, error) {
	syncedNodes := make([]*data.NodeData, 0)
	for _, node := range bnp.syncedNodes {
		syncedNodes = append(syncedNodes, node)
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	for _, node := range bnp.syncedFallbackNodes {
		syncedNodes = append(syncedNodes, node)
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	for _, backupNode := range bnp.lastSyncedNodes {
		syncedNodes = append(syncedNodes, backupNode)
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	return nil, ErrEmptyObserversList
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
	shardIDs := getSortedShardIDsSlice(nodesOnShards)

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

func getSortedShardIDsSlice(nodesOnShards map[uint32][]*data.NodeData) []uint32 {
	shardIDs := make([]uint32, 0)
	for shardID := range nodesOnShards {
		shardIDs = append(shardIDs, shardID)
	}
	sort.SliceStable(shardIDs, func(i, j int) bool {
		return shardIDs[i] < shardIDs[j]
	})

	return shardIDs
}
