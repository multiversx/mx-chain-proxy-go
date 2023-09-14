package holder

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

var (
	errEmptyShardIDsList  = errors.New("empty shard IDs list")
	errWrongConfiguration = errors.New("wrong observers configuration")
	log                   = logger.GetOrCreate("observer/holder")
)

type nodesHolder struct {
	mut                    sync.RWMutex
	syncedNodes            []*data.NodeData
	outOfSyncNodes         []*data.NodeData
	syncedFallbackNodes    []*data.NodeData
	outOfSyncFallbackNodes []*data.NodeData
	lastSyncedNodes        map[uint32]*data.NodeData
	shardIDs               []uint32
	availability           data.ObserverDataAvailabilityType
}

// NewNodesHolder will return a new instance of a nodesHolder
func NewNodesHolder(syncedNodes []*data.NodeData, fallbackNodes []*data.NodeData, shardIDs []uint32, availability data.ObserverDataAvailabilityType) (*nodesHolder, error) {
	if len(shardIDs) == 0 {
		return nil, errEmptyShardIDsList
	}

	return &nodesHolder{
		syncedNodes:            syncedNodes,
		outOfSyncNodes:         make([]*data.NodeData, 0),
		syncedFallbackNodes:    fallbackNodes,
		outOfSyncFallbackNodes: make([]*data.NodeData, 0),
		lastSyncedNodes:        make(map[uint32]*data.NodeData),
		shardIDs:               shardIDs,
		availability:           availability,
	}, nil
}

// UpdateNodes will update the internal maps based on the provided nodes
func (nh *nodesHolder) UpdateNodes(nodesWithSyncStatus []*data.NodeData) {
	if len(nodesWithSyncStatus) == 0 {
		return
	}
	syncedNodes, syncedFallbackNodes, outOfSyncNodes, err := computeSyncedAndOutOfSyncNodes(nodesWithSyncStatus, nh.shardIDs, nh.availability)
	if err != nil {
		log.Error("cannot update nodes based on sync state", "error", err)
		return
	}

	sameNumOfSynced := len(nh.syncedNodes) == len(syncedNodes)
	sameNumOfSyncedFallback := len(nh.syncedFallbackNodes) == len(syncedFallbackNodes)
	if sameNumOfSynced && sameNumOfSyncedFallback && len(outOfSyncNodes) == 0 {
		nh.printSyncedNodesInShardsUnprotected()
		// early exit as all the nodes are in sync
		return
	}

	syncedNodesMap := nodesSliceToShardedMap(syncedNodes)
	syncedFallbackNodesMap := nodesSliceToShardedMap(syncedFallbackNodes)

	nh.removeOutOfSyncNodesUnprotected(outOfSyncNodes, syncedNodesMap, syncedFallbackNodesMap)
	nh.addSyncedNodesUnprotected(syncedNodes, syncedFallbackNodes)
	nh.printSyncedNodesInShardsUnprotected()
}

// GetSyncedNodes returns all the synced nodes
func (nh *nodesHolder) GetSyncedNodes() []*data.NodeData {
	nh.mut.RLock()
	defer nh.mut.RUnlock()

	return copyNodes(nh.syncedNodes)
}

// GetSyncedFallbackNodes returns all the synced fallback nodes
func (nh *nodesHolder) GetSyncedFallbackNodes() []*data.NodeData {
	nh.mut.RLock()
	defer nh.mut.RUnlock()

	return copyNodes(nh.syncedFallbackNodes)
}

// GetOutOfSyncNodes returns all the out of sync nodes
func (nh *nodesHolder) GetOutOfSyncNodes() []*data.NodeData {
	nh.mut.RLock()
	defer nh.mut.RUnlock()

	return copyNodes(nh.outOfSyncNodes)
}

// GetOutOfSyncFallbackNodes returns all the out of sync fallback nodes
func (nh *nodesHolder) GetOutOfSyncFallbackNodes() []*data.NodeData {
	nh.mut.RLock()
	defer nh.mut.RUnlock()

	return copyNodes(nh.outOfSyncFallbackNodes)
}

// GetLastSyncedNodes returns the internal map of the last synced nodes
func (nh *nodesHolder) GetLastSyncedNodes() map[uint32]*data.NodeData {
	mapCopy := make(map[uint32]*data.NodeData, 0)
	nh.mut.RLock()
	for key, value := range nh.lastSyncedNodes {
		mapCopy[key] = value
	}
	nh.mut.RUnlock()

	return mapCopy
}

// IsInterfaceNil returns true if there is no value under the interface
func (nh *nodesHolder) IsInterfaceNil() bool {
	return nh == nil
}

func copyNodes(nodes []*data.NodeData) []*data.NodeData {
	sliceCopy := make([]*data.NodeData, 0, len(nodes))
	for _, node := range nodes {
		sliceCopy = append(sliceCopy, node)
	}

	return sliceCopy
}

func (nh *nodesHolder) printSyncedNodesInShardsUnprotected() {
	inSyncAddresses := make(map[uint32][]string, 0)
	for _, syncedNode := range nh.syncedNodes {
		inSyncAddresses[syncedNode.ShardId] = append(inSyncAddresses[syncedNode.ShardId], syncedNode.Address)
	}

	inSyncFallbackAddresses := make(map[uint32][]string, 0)
	for _, syncedFallbackNode := range nh.syncedFallbackNodes {
		inSyncFallbackAddresses[syncedFallbackNode.ShardId] = append(inSyncFallbackAddresses[syncedFallbackNode.ShardId], syncedFallbackNode.Address)
	}

	for _, shardID := range nh.shardIDs {
		totalNumOfActiveNodes := len(inSyncAddresses[shardID]) + len(inSyncFallbackAddresses[shardID])
		// if none of them is active, use the backup if exists
		hasBackup := nh.lastSyncedNodes[shardID] != nil
		if totalNumOfActiveNodes == 0 && hasBackup {
			totalNumOfActiveNodes++
			inSyncAddresses[shardID] = append(inSyncAddresses[shardID], nh.lastSyncedNodes[shardID].Address)
		}
		nodesType := "regular active nodes"
		if nh.availability == data.AvailabilityRecent {
			nodesType = "snapshotless active nodes"
		}
		log.Info(fmt.Sprintf("shard %d %s", shardID, nodesType),
			"observers count", totalNumOfActiveNodes,
			"addresses", strings.Join(inSyncAddresses[shardID], ", "),
			"fallback addresses", strings.Join(inSyncFallbackAddresses[shardID], ", "))
	}
}

func computeSyncedAndOutOfSyncNodes(nodes []*data.NodeData, shardIDs []uint32, availability data.ObserverDataAvailabilityType) ([]*data.NodeData, []*data.NodeData, []*data.NodeData, error) {
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
			if availability != data.AvailabilityRecent {
				return nil, nil, nil, fmt.Errorf("%w for shard %d - no synced or not synced node", errWrongConfiguration, shardID)
			}
		}
	}

	return syncedNodes, syncedFallbackNodes, notSyncedNodes, nil
}

func (nh *nodesHolder) addSyncedNodesUnprotected(receivedSyncedNodes []*data.NodeData, receivedSyncedFallbackNodes []*data.NodeData) {
	syncedNodesPerShard := make(map[uint32][]string)
	for _, node := range receivedSyncedNodes {
		nh.removeFromOutOfSyncIfNeededUnprotected(node)
		syncedNodesPerShard[node.ShardId] = append(syncedNodesPerShard[node.ShardId], node.Address)
		if nh.isReceivedSyncedNodeExistent(node) {
			continue
		}

		nh.syncedNodes = append(nh.syncedNodes, node)
	}

	for _, node := range receivedSyncedFallbackNodes {
		nh.removeFromOutOfSyncIfNeededUnprotected(node)
		if nh.isReceivedSyncedNodeExistentAsFallback(node) {
			continue
		}

		nh.syncedFallbackNodes = append(nh.syncedFallbackNodes, node)
	}

	// if there is at least one synced node regular received, clean the backup list
	for _, shardId := range nh.shardIDs {
		if len(syncedNodesPerShard[shardId]) != 0 {
			delete(nh.lastSyncedNodes, shardId)
		}
	}
}

func (nh *nodesHolder) removeFromOutOfSyncIfNeededUnprotected(node *data.NodeData) {
	if node.IsFallback {
		nh.removeFallbackFromOutOfSyncListUnprotected(node)
		return
	}

	nh.removeRegularFromOutOfSyncListUnprotected(node)
}

func (nh *nodesHolder) isReceivedSyncedNodeExistent(receivedNode *data.NodeData) bool {
	for _, node := range nh.syncedNodes {
		if node.Address == receivedNode.Address && node.ShardId == receivedNode.ShardId {
			return true
		}
	}

	return false
}

func (nh *nodesHolder) isReceivedSyncedNodeExistentAsFallback(receivedNode *data.NodeData) bool {
	for _, node := range nh.syncedFallbackNodes {
		if node.Address == receivedNode.Address && node.ShardId == receivedNode.ShardId {
			return true
		}
	}

	return false
}

func (nh *nodesHolder) addToOutOfSyncUnprotected(node *data.NodeData) {
	if node.IsFallback {
		nh.addFallbackToOutOfSyncUnprotected(node)
		return
	}

	nh.addRegularToOutOfSyncUnprotected(node)
}

func (nh *nodesHolder) addRegularToOutOfSyncUnprotected(node *data.NodeData) {
	for _, oosNode := range nh.outOfSyncNodes {
		if oosNode.Address == node.Address && oosNode.ShardId == node.ShardId {
			return
		}
	}

	nh.outOfSyncNodes = append(nh.outOfSyncNodes, node)
}

func (nh *nodesHolder) addFallbackToOutOfSyncUnprotected(node *data.NodeData) {
	for _, oosNode := range nh.outOfSyncFallbackNodes {
		if oosNode.Address == node.Address && oosNode.ShardId == node.ShardId {
			return
		}
	}

	nh.outOfSyncFallbackNodes = append(nh.outOfSyncFallbackNodes, node)
}

func (nh *nodesHolder) removeOutOfSyncNodesUnprotected(
	outOfSyncNodes []*data.NodeData,
	syncedNodesMap map[uint32][]*data.NodeData,
	syncedFallbackNodesMap map[uint32][]*data.NodeData,
) {
	minSyncedNodes := 1
	if nh.availability == data.AvailabilityRecent {
		minSyncedNodes = 0 // allow the snapshotless list to be empty so regular observers can be used
	}
	if len(outOfSyncNodes) == 0 {
		nh.outOfSyncNodes = make([]*data.NodeData, 0)
		nh.outOfSyncFallbackNodes = make([]*data.NodeData, 0)
		return
	}

	for _, outOfSyncNode := range outOfSyncNodes {
		hasOneSyncedNode := len(syncedNodesMap[outOfSyncNode.ShardId]) >= minSyncedNodes
		hasEnoughSyncedFallbackNodes := len(syncedFallbackNodesMap[outOfSyncNode.ShardId]) > minSyncedNodes
		canDeleteFallbackNode := hasOneSyncedNode || hasEnoughSyncedFallbackNodes
		if outOfSyncNode.IsFallback && canDeleteFallbackNode {
			nh.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// if trying to delete last fallback, use last known synced node
		// if backup node does not exist, keep fallback
		hasBackup := nh.lastSyncedNodes[outOfSyncNode.ShardId] != nil
		if outOfSyncNode.IsFallback && hasBackup {
			nh.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		hasEnoughSyncedNodes := len(syncedNodesMap[outOfSyncNode.ShardId]) >= minSyncedNodes
		if hasEnoughSyncedNodes {
			nh.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// trying to remove last synced node
		// if fallbacks are available, save this one as backup and use fallbacks
		// else, keep using this one
		// save this last regular observer as backup in case fallbacks go offline
		// also, if this is the old fallback observer which didn't get synced, keep it in list
		wasSyncedAtPreviousStep := nh.isReceivedSyncedNodeExistent(outOfSyncNode)
		isBackupObserver := nh.lastSyncedNodes[outOfSyncNode.ShardId] == outOfSyncNode
		isRegularSyncedBefore := !outOfSyncNode.IsFallback && wasSyncedAtPreviousStep
		if isRegularSyncedBefore || isBackupObserver {
			log.Info("backup observer updated",
				"address", outOfSyncNode.Address,
				"is fallback", outOfSyncNode.IsFallback,
				"shard", outOfSyncNode.ShardId)
			nh.lastSyncedNodes[outOfSyncNode.ShardId] = outOfSyncNode
		}
		hasOneSyncedFallbackNode := len(syncedFallbackNodesMap[outOfSyncNode.ShardId]) >= minSyncedNodes
		if hasOneSyncedFallbackNode {
			nh.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// safe to delete regular observer, as it is already in lastSyncedNodes map
		if !outOfSyncNode.IsFallback {
			nh.removeNodeUnprotected(outOfSyncNode)
			continue
		}

		// this is a fallback node, with no synced nodes.
		// save it as backup and delete it from its list
		nh.lastSyncedNodes[outOfSyncNode.ShardId] = outOfSyncNode
		nh.removeNodeUnprotected(outOfSyncNode)
	}
}

func (nh *nodesHolder) removeNodeUnprotected(node *data.NodeData) {
	nh.removeNodeFromSyncedNodesUnprotected(node)
	nh.addToOutOfSyncUnprotected(node)
}

func (nh *nodesHolder) removeNodeFromSyncedNodesUnprotected(nodeToRemove *data.NodeData) {
	if nodeToRemove.IsFallback {
		nh.removeFallbackFromSyncedListUnprotected(nodeToRemove)
		return
	}

	nh.removeRegularFromSyncedListUnprotected(nodeToRemove)
}

func (nh *nodesHolder) removeRegularFromSyncedListUnprotected(nodeToRemove *data.NodeData) {
	nodeIndex := getIndexFromList(nodeToRemove, nh.syncedNodes)
	if nodeIndex == -1 {
		return
	}

	copy(nh.syncedNodes[nodeIndex:], nh.syncedNodes[nodeIndex+1:])
	nh.syncedNodes[len(nh.syncedNodes)-1] = nil
	nh.syncedNodes = nh.syncedNodes[:len(nh.syncedNodes)-1]
}

func (nh *nodesHolder) removeFallbackFromSyncedListUnprotected(nodeToRemove *data.NodeData) {
	nodeIndex := getIndexFromList(nodeToRemove, nh.syncedFallbackNodes)
	if nodeIndex == -1 {
		return
	}

	copy(nh.syncedFallbackNodes[nodeIndex:], nh.syncedFallbackNodes[nodeIndex+1:])
	nh.syncedFallbackNodes[len(nh.syncedFallbackNodes)-1] = nil
	nh.syncedFallbackNodes = nh.syncedFallbackNodes[:len(nh.syncedFallbackNodes)-1]
}

func (nh *nodesHolder) removeRegularFromOutOfSyncListUnprotected(nodeToRemove *data.NodeData) {
	nodeIndex := getIndexFromList(nodeToRemove, nh.outOfSyncNodes)
	if nodeIndex == -1 {
		return
	}

	copy(nh.outOfSyncNodes[nodeIndex:], nh.outOfSyncNodes[nodeIndex+1:])
	nh.outOfSyncNodes[len(nh.outOfSyncNodes)-1] = nil
	nh.outOfSyncNodes = nh.outOfSyncNodes[:len(nh.outOfSyncNodes)-1]
}

func (nh *nodesHolder) removeFallbackFromOutOfSyncListUnprotected(nodeToRemove *data.NodeData) {
	nodeIndex := getIndexFromList(nodeToRemove, nh.outOfSyncFallbackNodes)
	if nodeIndex == -1 {
		return
	}

	copy(nh.outOfSyncFallbackNodes[nodeIndex:], nh.outOfSyncFallbackNodes[nodeIndex+1:])
	nh.outOfSyncFallbackNodes[len(nh.outOfSyncFallbackNodes)-1] = nil
	nh.outOfSyncFallbackNodes = nh.outOfSyncFallbackNodes[:len(nh.outOfSyncFallbackNodes)-1]
}

func getIndexFromList(providedNode *data.NodeData, list []*data.NodeData) int {
	nodeIndex := -1
	for idx, node := range list {
		if node.Address == providedNode.Address && node.ShardId == providedNode.ShardId {
			nodeIndex = idx
			break
		}
	}

	return nodeIndex
}

func nodesSliceToShardedMap(nodes []*data.NodeData) map[uint32][]*data.NodeData {
	newNodes := make(map[uint32][]*data.NodeData)
	for _, node := range nodes {
		shardId := node.ShardId
		newNodes[shardId] = append(newNodes[shardId], node)
	}

	return newNodes
}
