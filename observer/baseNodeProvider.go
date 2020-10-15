package observer

import (
	"sort"
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type baseNodeProvider struct {
	mutNodes sync.RWMutex
	nodes    map[uint32][]*data.NodeData
	allNodes []*data.NodeData
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

func initAllNodesSlice(nodesOnShards map[uint32][]*data.NodeData) []*data.NodeData {
	sliceToReturn := make([]*data.NodeData, 0)
	shardIDs := make([]uint32, 0)
	for shardID := range nodesOnShards {
		shardIDs = append(shardIDs, shardID)
	}
	sort.SliceStable(shardIDs, func(i, j int) bool {
		return shardIDs[i] < shardIDs[j]
	})

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
