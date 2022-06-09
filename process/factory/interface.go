package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/observer"
)

// Processor defines what a processor should be able to do
type Processor interface {
	ComputeShardId(addressBuff []byte) (uint32, error)
	CallGetRestEndPoint(address string, path string, value interface{}) (int, error)
	CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error)
	GetObserversOnePerShard() ([]*data.NodeData, error)
	GetShardIDs() []uint32
	GetFullHistoryNodesOnePerShard() ([]*data.NodeData, error)
	GetObservers(shardID uint32) ([]*data.NodeData, error)
	GetAllObservers() ([]*data.NodeData, error)
	GetFullHistoryNodes(shardID uint32) ([]*data.NodeData, error)
	GetAllFullHistoryNodes() ([]*data.NodeData, error)
	GetShardCoordinator() sharding.Coordinator
	GetPubKeyConverter() core.PubkeyConverter
	GetObserverProvider() observer.NodesProviderHandler
	GetFullHistoryNodesProvider() observer.NodesProviderHandler
	IsInterfaceNil() bool
}

// PrivateKeysLoaderHandler defines what a component which handles loading of the private keys file should do
type PrivateKeysLoaderHandler interface {
	PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error)
}
