package observer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/config"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/holder"
)

type baseNodeProvider struct {
	mutNodes              sync.RWMutex
	shardIds              []uint32
	numOfShards           uint32
	configurationFilePath string
	regularNodes          NodesHolder
	snapshotlessNodes     NodesHolder
}

func (bnp *baseNodeProvider) initNodes(nodes []*data.NodeData) error {
	if len(nodes) == 0 {
		return ErrEmptyObserversList
	}

	newNodes := make(map[uint32][]*data.NodeData)
	for _, observer := range nodes {
		shardId := observer.ShardId
		newNodes[shardId] = append(newNodes[shardId], observer)
		isMeta := shardId == core.MetachainShardId
		if isMeta {
			continue
		}

		if shardId >= bnp.numOfShards {
			return fmt.Errorf("%w for observer %s, provided shard %d, number of shards configured %d",
				ErrInvalidShard,
				observer.Address,
				observer.ShardId,
				bnp.numOfShards,
			)
		}
	}

	err := checkNodesInShards(newNodes)
	if err != nil {
		return err
	}

	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()

	bnp.shardIds = getSortedShardIDsSlice(newNodes)
	syncedNodes, syncedFallbackNodes, syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes := initAllNodesSlice(newNodes)
	bnp.regularNodes, err = holder.NewNodesHolder(syncedNodes, syncedFallbackNodes, data.AvailabilityAll)
	if err != nil {
		return err
	}
	bnp.snapshotlessNodes, err = holder.NewNodesHolder(syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes, data.AvailabilityRecent)
	if err != nil {
		return err
	}

	return nil
}

func checkNodesInShards(nodes map[uint32][]*data.NodeData) error {
	for shardID, nodesInShard := range nodes {
		atLeastOneRegularNode := false
		for _, node := range nodesInShard {
			if !node.IsSnapshotless {
				atLeastOneRegularNode = true
				break
			}
		}
		if !atLeastOneRegularNode {
			return fmt.Errorf("observers for shard %d must include at least one historical (non-snapshotless) observer", shardID)
		}
	}

	return nil
}

// GetAllNodesWithSyncState will return the merged list of active observers and out of sync observers
func (bnp *baseNodeProvider) GetAllNodesWithSyncState() []*data.NodeData {
	bnp.mutNodes.RLock()
	defer bnp.mutNodes.RUnlock()

	nodesSlice := make([]*data.NodeData, 0)
	for _, shardID := range bnp.shardIds {
		nodesSlice = append(nodesSlice, bnp.regularNodes.GetSyncedNodes(shardID)...)
		nodesSlice = append(nodesSlice, bnp.regularNodes.GetOutOfSyncNodes(shardID)...)
		nodesSlice = append(nodesSlice, bnp.regularNodes.GetSyncedFallbackNodes(shardID)...)
		nodesSlice = append(nodesSlice, bnp.regularNodes.GetOutOfSyncFallbackNodes(shardID)...)

		nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetSyncedNodes(shardID)...)
		nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetOutOfSyncNodes(shardID)...)
		nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetSyncedFallbackNodes(shardID)...)
		nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetOutOfSyncFallbackNodes(shardID)...)
	}

	return nodesSlice
}

// UpdateNodesBasedOnSyncState will simply call the corresponding function for both regular and snapshotless observers
func (bnp *baseNodeProvider) UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData) {
	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()

	regularNodes, snapshotlessNodes := splitNodesByDataAvailability(nodesWithSyncStatus)
	bnp.regularNodes.UpdateNodes(regularNodes)
	bnp.snapshotlessNodes.UpdateNodes(snapshotlessNodes)
}

// PrintNodesInShards will only print the nodes in shards
func (bnp *baseNodeProvider) PrintNodesInShards() {
	bnp.mutNodes.RLock()
	defer bnp.mutNodes.RUnlock()

	bnp.regularNodes.PrintNodesInShards()
	bnp.snapshotlessNodes.PrintNodesInShards()
}

func splitNodesByDataAvailability(nodes []*data.NodeData) ([]*data.NodeData, []*data.NodeData) {
	regularNodes := make([]*data.NodeData, 0)
	snapshotlessNodes := make([]*data.NodeData, 0)
	for _, node := range nodes {
		if node.IsSnapshotless {
			snapshotlessNodes = append(snapshotlessNodes, node)
		} else {
			regularNodes = append(regularNodes, node)
		}
	}

	return regularNodes, snapshotlessNodes
}

// ReloadNodes will reload the observers or the full history observers
func (bnp *baseNodeProvider) ReloadNodes(nodesType data.NodeType) data.NodesReloadResponse {
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

	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()
	bnp.shardIds = getSortedShardIDsSlice(newNodes)
	syncedNodes, syncedFallbackNodes, syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes := initAllNodesSlice(newNodes)
	bnp.regularNodes, err = holder.NewNodesHolder(syncedNodes, syncedFallbackNodes, data.AvailabilityAll)
	if err != nil {
		log.Error("cannot reload regular nodes: NewNodesHolder", "error", err)
		return data.NodesReloadResponse{
			OkRequest:   true,
			Description: "not reloaded",
			Error:       "cannot create the regular nodes holder: " + err.Error(),
		}
	}

	bnp.snapshotlessNodes, err = holder.NewNodesHolder(syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes, data.AvailabilityRecent)
	if err != nil {
		log.Error("cannot reload snapshotless nodes: NewNodesHolder", "error", err)
		return data.NodesReloadResponse{
			OkRequest:   true,
			Description: "not reloaded",
			Error:       "cannot create the snapshotless nodes holder: " + err.Error(),
		}
	}

	return data.NodesReloadResponse{
		OkRequest:   true,
		Description: prepareReloadResponseMessage(newNodes),
		Error:       "",
	}
}

func (bnp *baseNodeProvider) getSyncedNodesForShardUnprotected(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	var syncedNodes []*data.NodeData

	syncedNodes = bnp.getSyncedNodes(dataAvailability, shardID)
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	fallbackNodesSource := bnp.getFallbackNodes(dataAvailability, shardID)
	if len(fallbackNodesSource) != 0 {
		return fallbackNodesSource, nil
	}

	outOfSyncNodes := bnp.getOutOfSyncNodes(dataAvailability, shardID)
	if len(outOfSyncNodes) > 0 {
		return outOfSyncNodes, nil
	}

	outOfSyncFallbackNodesSource := bnp.getOutOfSyncFallbackNodes(dataAvailability, shardID)
	if len(outOfSyncFallbackNodesSource) != 0 {
		return outOfSyncFallbackNodesSource, nil
	}

	return nil, ErrShardNotAvailable
}

func (bnp *baseNodeProvider) getNodesByType(
	availabilityType data.ObserverDataAvailabilityType,
	shardID uint32,
	getSnapshotlessNodesFunc func(uint32) []*data.NodeData,
	getRegularNodesFunc func(uint32) []*data.NodeData) []*data.NodeData {

	if availabilityType == data.AvailabilityRecent {
		nodes := getSnapshotlessNodesFunc(shardID)
		if len(nodes) > 0 {
			return nodes
		}
	}
	return getRegularNodesFunc(shardID)
}

func (bnp *baseNodeProvider) getSyncedNodes(availabilityType data.ObserverDataAvailabilityType, shardID uint32) []*data.NodeData {
	return bnp.getNodesByType(availabilityType, shardID, bnp.snapshotlessNodes.GetSyncedNodes, bnp.regularNodes.GetSyncedNodes)
}

func (bnp *baseNodeProvider) getFallbackNodes(availabilityType data.ObserverDataAvailabilityType, shardID uint32) []*data.NodeData {
	return bnp.getNodesByType(availabilityType, shardID, bnp.snapshotlessNodes.GetSyncedFallbackNodes, bnp.regularNodes.GetSyncedFallbackNodes)
}

func (bnp *baseNodeProvider) getOutOfSyncNodes(availabilityType data.ObserverDataAvailabilityType, shardID uint32) []*data.NodeData {
	return bnp.getNodesByType(availabilityType, shardID, bnp.snapshotlessNodes.GetOutOfSyncNodes, bnp.regularNodes.GetOutOfSyncNodes)
}

func (bnp *baseNodeProvider) getOutOfSyncFallbackNodes(availabilityType data.ObserverDataAvailabilityType, shardID uint32) []*data.NodeData {
	return bnp.getNodesByType(availabilityType, shardID, bnp.snapshotlessNodes.GetOutOfSyncFallbackNodes, bnp.regularNodes.GetOutOfSyncFallbackNodes)
}

func (bnp *baseNodeProvider) getSyncedNodesUnprotected(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	syncedNodes := make([]*data.NodeData, 0)
	for _, shardId := range bnp.shardIds {
		syncedShardNodes, err := bnp.getSyncedNodesForShardUnprotected(shardId, dataAvailability)
		if err != nil {
			return nil, fmt.Errorf("%w for shard %d", err, shardId)
		}

		syncedNodes = append(syncedNodes, syncedShardNodes...)
	}

	if len(syncedNodes) == 0 {
		return nil, ErrEmptyObserversList
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

func initAllNodesSlice(nodesOnShards map[uint32][]*data.NodeData) ([]*data.NodeData, []*data.NodeData, []*data.NodeData, []*data.NodeData) {
	eligibleNodes := make([]*data.NodeData, 0)
	fallbackNodes := make([]*data.NodeData, 0)
	eligibleSnapshotlessNodes := make([]*data.NodeData, 0)
	fallbackSnapshotlessNodes := make([]*data.NodeData, 0)
	shardIDs := getSortedShardIDsSlice(nodesOnShards)

	finishedShards := make(map[uint32]struct{})
	for i := 0; ; i++ {
		for _, shardID := range shardIDs {
			if i >= len(nodesOnShards[shardID]) {
				finishedShards[shardID] = struct{}{}
				continue
			}

			node := nodesOnShards[shardID][i]
			if node.IsSnapshotless {
				if node.IsFallback {
					fallbackSnapshotlessNodes = append(fallbackSnapshotlessNodes, node)
				} else {
					eligibleSnapshotlessNodes = append(eligibleSnapshotlessNodes, node)
				}
			} else {
				if node.IsFallback {
					fallbackNodes = append(fallbackNodes, node)
				} else {
					eligibleNodes = append(eligibleNodes, node)
				}
			}
		}

		if len(finishedShards) == len(nodesOnShards) {
			break
		}
	}

	return eligibleNodes, fallbackNodes, eligibleSnapshotlessNodes, fallbackSnapshotlessNodes
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
