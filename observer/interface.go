package observer

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// NodesProviderHandler defines what a nodes provider should be able to do
type NodesProviderHandler interface {
	GetNodesByShardId(shardId uint32) ([]*data.NodeData, error)
	GetAllNodes() ([]*data.NodeData, error)
	UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData)
	GetAllNodesWithSyncState() []*data.NodeData
	ReloadNodes(nodesType data.NodeType) data.NodesReloadResponse
	IsInterfaceNil() bool
}
