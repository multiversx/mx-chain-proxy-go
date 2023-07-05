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
	}

	err := checkNodesInShards(newNodes)
	if err != nil {
		return err
	}

	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()

	bnp.shardIds = getSortedShardIDsSlice(newNodes)
	syncedNodes, syncedFallbackNodes, syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes := initAllNodesSlice(newNodes)
	bnp.regularNodes, err = holder.NewNodesHolder(syncedNodes, syncedFallbackNodes, bnp.shardIds, data.AvailabilityAll)
	if err != nil {
		return err
	}
	bnp.snapshotlessNodes, err = holder.NewNodesHolder(syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes, bnp.shardIds, data.AvailabilityRecent)
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
	nodesSlice = append(nodesSlice, bnp.regularNodes.GetSyncedNodes()...)
	nodesSlice = append(nodesSlice, bnp.regularNodes.GetOutOfSyncNodes()...)
	nodesSlice = append(nodesSlice, bnp.regularNodes.GetSyncedFallbackNodes()...)
	nodesSlice = append(nodesSlice, bnp.regularNodes.GetOutOfSyncFallbackNodes()...)

	nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetSyncedNodes()...)
	nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetOutOfSyncNodes()...)
	nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetSyncedFallbackNodes()...)
	nodesSlice = append(nodesSlice, bnp.snapshotlessNodes.GetOutOfSyncFallbackNodes()...)
	return nodesSlice
}

// UpdateNodesBasedOnSyncState will handle the nodes lists, by removing out of sync observers or by adding back observers
// that were previously removed because they were out of sync.
// If all observers are removed, the last one synced will be saved and the fallbacks will be used.
// If even the fallbacks are out of sync, the last regular observer synced will be used, even though it is out of sync.
// When one or more regular observers are back in sync, the fallbacks will not be used anymore.
func (bnp *baseNodeProvider) UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData) {
	bnp.mutNodes.Lock()
	defer bnp.mutNodes.Unlock()

	regularNodes, snapshotlessNodes := splitNodesByDataAvailability(nodesWithSyncStatus)
	bnp.regularNodes.UpdateNodes(regularNodes)
	bnp.snapshotlessNodes.UpdateNodes(snapshotlessNodes)
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
	defer bnp.mutNodes.Unlock()
	bnp.shardIds = getSortedShardIDsSlice(newNodes)
	syncedNodes, syncedFallbackNodes, syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes := initAllNodesSlice(newNodes)
	bnp.regularNodes, err = holder.NewNodesHolder(syncedNodes, syncedFallbackNodes, bnp.shardIds, data.AvailabilityAll)
	if err != nil {
		log.Error("cannot reload regular nodes: NewNodesHolder", "error", err)
		return data.NodesReloadResponse{
			OkRequest:   true,
			Description: "not reloaded",
			Error:       "cannot create the regular nodes holder: " + err.Error(),
		}
	}
	bnp.snapshotlessNodes, err = holder.NewNodesHolder(syncedSnapshotlessNodes, syncedSnapshotlessFallbackNodes, bnp.shardIds, data.AvailabilityRecent)
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

func (bnp *baseNodeProvider) getSyncedNodesForShardUnprotected(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	var syncedNodesSource []*data.NodeData
	if dataAvailability == data.AvailabilityRecent && len(bnp.snapshotlessNodes.GetSyncedNodes()) > 0 {
		syncedNodesSource = bnp.snapshotlessNodes.GetSyncedNodes()
	} else {
		syncedNodesSource = bnp.regularNodes.GetSyncedNodes()
	}
	syncedNodes := make([]*data.NodeData, 0)
	for _, node := range syncedNodesSource {
		if node.ShardId != shardId {
			continue
		}

		syncedNodes = append(syncedNodes, node)
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	var fallbackNodesSource []*data.NodeData
	if dataAvailability == data.AvailabilityRecent {
		fallbackNodesSource = bnp.snapshotlessNodes.GetSyncedNodes()
	} else {
		fallbackNodesSource = bnp.regularNodes.GetSyncedNodes()
	}
	for _, node := range fallbackNodesSource {
		if node.ShardId == shardId {
			syncedNodes = append(syncedNodes, node)
		}
	}
	if len(syncedNodes) != 0 {
		return syncedNodes, nil
	}

	var lastSyncedNodesMap map[uint32]*data.NodeData
	if dataAvailability == data.AvailabilityAll {
		lastSyncedNodesMap = bnp.regularNodes.GetLastSyncedNodes()
	} else {
		lastSyncedNodesMap = bnp.snapshotlessNodes.GetLastSyncedNodes()
	}
	backupNode, hasBackup := lastSyncedNodesMap[shardId]
	if hasBackup {
		return []*data.NodeData{backupNode}, nil
	}

	return nil, ErrShardNotAvailable
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
