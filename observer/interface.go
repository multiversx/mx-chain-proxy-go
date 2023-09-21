package observer

import "github.com/multiversx/mx-chain-proxy-go/data"

// NodesProviderHandler defines what a nodes provider should be able to do
type NodesProviderHandler interface {
	GetNodesByShardId(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllNodes(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData)
	GetAllNodesWithSyncState() []*data.NodeData
	ReloadNodes(nodesType data.NodeType) data.NodesReloadResponse
	IsInterfaceNil() bool
}

// NodesHolder defines the actions of a component that is able to hold nodes
type NodesHolder interface {
	UpdateNodes(nodesWithSyncStatus []*data.NodeData)
	GetSyncedNodes(shardID uint32) []*data.NodeData
	GetSyncedFallbackNodes(shardID uint32) []*data.NodeData
	GetOutOfSyncNodes(shardID uint32) []*data.NodeData
	GetOutOfSyncFallbackNodes(shardID uint32) []*data.NodeData
	Count() int
	IsInterfaceNil() bool
}
