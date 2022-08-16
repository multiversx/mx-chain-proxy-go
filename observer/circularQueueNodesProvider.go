package observer

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// circularQueueNodesProvider will handle the providing of observers in a circular queue way, guaranteeing the
// balancing of them
type circularQueueNodesProvider struct {
	*baseNodeProvider
	countersMap        map[uint32]uint32
	counterForAllNodes uint32
	mutCounters        sync.RWMutex
}

// NewCircularQueueNodesProvider returns a new instance of circularQueueNodesProvider
func NewCircularQueueNodesProvider(observers []*data.NodeData, configurationFilePath string) (*circularQueueNodesProvider, error) {
	bop := &baseNodeProvider{
		configurationFilePath: configurationFilePath,
	}

	err := bop.initNodesMaps(observers)
	if err != nil {
		return nil, err
	}

	countersMap := make(map[uint32]uint32)
	return &circularQueueNodesProvider{
		baseNodeProvider:   bop,
		countersMap:        countersMap,
		counterForAllNodes: 0,
	}, nil
}

// GetNodesByShardId will return a slice of observers for the given shard
func (cqnp *circularQueueNodesProvider) GetNodesByShardId(shardId uint32) ([]*data.NodeData, error) {
	cqnp.mutNodes.Lock()
	defer cqnp.mutNodes.Unlock()

	syncedNodesForShard, err := cqnp.getSyncedNodesForShardUnprotected(shardId)
	if err != nil {
		return nil, err
	}

	position := cqnp.computeCounterForShard(shardId, uint32(len(syncedNodesForShard)))
	sliceToRet := append(syncedNodesForShard[position:], syncedNodesForShard[:position]...)

	return sliceToRet, nil
}

// GetAllNodes will return a slice containing all observers
func (cqnp *circularQueueNodesProvider) GetAllNodes() ([]*data.NodeData, error) {
	cqnp.mutNodes.Lock()
	defer cqnp.mutNodes.Unlock()

	allNodes := cqnp.syncedNodes

	position := cqnp.computeCounterForAllNodes(uint32(len(allNodes)))
	sliceToRet := append(allNodes[position:], allNodes[:position]...)

	return sliceToRet, nil
}

func (cqnp *circularQueueNodesProvider) computeCounterForShard(shardID uint32, lenNodes uint32) uint32 {
	cqnp.mutCounters.Lock()
	defer cqnp.mutCounters.Unlock()

	cqnp.countersMap[shardID]++
	cqnp.countersMap[shardID] %= lenNodes

	return cqnp.countersMap[shardID]
}

func (cqnp *circularQueueNodesProvider) computeCounterForAllNodes(lenNodes uint32) uint32 {
	cqnp.mutCounters.Lock()
	defer cqnp.mutCounters.Unlock()

	cqnp.counterForAllNodes++
	cqnp.counterForAllNodes %= lenNodes

	return cqnp.counterForAllNodes
}

// IsInterfaceNil returns true if there is no value under the interface
func (cqnp *circularQueueNodesProvider) IsInterfaceNil() bool {
	return cqnp == nil
}
