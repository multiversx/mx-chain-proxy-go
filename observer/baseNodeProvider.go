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
	nodes                 map[uint32][]*data.NodeData
	configurationFilePath string
	allNodes              []*data.NodeData
}

func (bop *baseNodeProvider) initNodesMaps(nodes []*data.NodeData) error {
	if len(nodes) == 0 {
		return ErrEmptyObserversList
	}

	newNodes := make(map[uint32][]*data.NodeData)
	for _, observer := range nodes {
		shardId := observer.ShardId
		newNodes[shardId] = append(newNodes[shardId], observer)
	}

	bop.mutNodes.Lock()
	bop.nodes = newNodes
	bop.allNodes = initAllNodesSlice(newNodes)
	bop.mutNodes.Unlock()

	return nil
}

// ReloadNodes will reload the observers or the full history observers
func (bop *baseNodeProvider) ReloadNodes(nodesType data.NodeType) data.NodesReloadResponse {
	bop.mutNodes.RLock()
	oldNodes := bop.nodes
	bop.mutNodes.RUnlock()

	newConfig, err := loadMainConfig(bop.configurationFilePath)
	if err != nil {
		return data.NodesReloadResponse{
			OkRequest:   true,
			Description: "not reloaded",
			Error:       "cannot load configuration file at " + bop.configurationFilePath,
		}
	}

	nodes := newConfig.Observers
	if nodesType == data.FullHistoryNode {
		nodes = newConfig.FullHistoryNodes
	}

	newNodes := nodesSliceToShardedMap(nodes)
	numOldShardsCount := len(oldNodes)
	numNewShardsCount := len(newNodes)

	if numOldShardsCount != numNewShardsCount {
		return data.NodesReloadResponse{
			OkRequest:   false,
			Description: "not reloaded",
			Error:       fmt.Sprintf("different number of shards. before: %d, now: %d", numOldShardsCount, numNewShardsCount),
		}
	}

	bop.mutNodes.Lock()
	bop.nodes = newNodes
	bop.allNodes = initAllNodesSlice(newNodes)
	bop.mutNodes.Unlock()

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
