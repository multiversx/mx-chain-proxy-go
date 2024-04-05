package observer

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/mapCounters"
)

// circularQueueNodesProvider will handle the providing of observers in a circular queue way, guaranteeing the
// balancing of them
type circularQueueNodesProvider struct {
	*baseNodeProvider
	positionsHolder CounterMapsHolder
}

// NewCircularQueueNodesProvider returns a new instance of circularQueueNodesProvider
func NewCircularQueueNodesProvider(observers []*data.NodeData, configurationFilePath string) (*circularQueueNodesProvider, error) {
	bop := &baseNodeProvider{
		configurationFilePath: configurationFilePath,
	}

	err := bop.initNodes(observers)
	if err != nil {
		return nil, err
	}

	return &circularQueueNodesProvider{
		baseNodeProvider: bop,
		positionsHolder:  mapCounters.NewMapCountersHolder(),
	}, nil
}

// GetNodesByShardId will return a slice of observers for the given shard
func (cqnp *circularQueueNodesProvider) GetNodesByShardId(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	cqnp.mutNodes.Lock()
	defer cqnp.mutNodes.Unlock()

	syncedNodesForShard, err := cqnp.getSyncedNodesForShardUnprotected(shardId, dataAvailability)
	if err != nil {
		return nil, err
	}

	position, err := cqnp.positionsHolder.ComputeShardPosition(dataAvailability, shardId, uint32(len(syncedNodesForShard)))
	if err != nil {
		return nil, err
	}

	sliceToRet := append(syncedNodesForShard[position:], syncedNodesForShard[:position]...)

	return sliceToRet, nil
}

// GetAllNodes will return a slice containing all observers
func (cqnp *circularQueueNodesProvider) GetAllNodes(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	cqnp.mutNodes.Lock()
	defer cqnp.mutNodes.Unlock()

	allNodes, err := cqnp.getSyncedNodesUnprotected(dataAvailability)
	if err != nil {
		return nil, err
	}

	position, err := cqnp.positionsHolder.ComputeAllNodesPosition(dataAvailability, uint32(len(allNodes)))
	if err != nil {
		return nil, err
	}

	sliceToRet := append(allNodes[position:], allNodes[:position]...)

	return sliceToRet, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (cqnp *circularQueueNodesProvider) IsInterfaceNil() bool {
	return cqnp == nil
}
