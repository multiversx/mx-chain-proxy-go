package observer

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// simpleNodesProvider will handle the providing of observers in a simple way, in the order in which they were
// provided in the config file.
type simpleNodesProvider struct {
	*baseNodeProvider
}

// NewSimpleNodesProvider will return a new instance of simpleNodesProvider
func NewSimpleNodesProvider(observers []*data.NodeData, configurationFilePath string) (*simpleNodesProvider, error) {
	bop := &baseNodeProvider{
		configurationFilePath: configurationFilePath,
	}

	err := bop.initNodes(observers)
	if err != nil {
		return nil, err
	}

	return &simpleNodesProvider{
		baseNodeProvider: bop,
	}, nil
}

// GetNodesByShardId will return a slice of the nodes for the given shard
func (snp *simpleNodesProvider) GetNodesByShardId(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	snp.mutNodes.RLock()
	defer snp.mutNodes.RUnlock()

	return snp.getSyncedNodesForShardUnprotected(shardId, dataAvailability)
}

// GetAllNodes will return a slice containing all the nodes
func (snp *simpleNodesProvider) GetAllNodes(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	snp.mutNodes.RLock()
	defer snp.mutNodes.RUnlock()

	return snp.getSyncedNodesUnprotected(dataAvailability)
}

// IsInterfaceNil returns true if there is no value under the interface
func (snp *simpleNodesProvider) IsInterfaceNil() bool {
	return snp == nil
}
