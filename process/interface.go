package process

import (
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Processor defines what a processor should be able to do
type Processor interface {
	GetObservers(shardID uint32) ([]*data.NodeData, error)
	GetAllObservers() ([]*data.NodeData, error)
	GetObserversOnePerShard() ([]*data.NodeData, error)
	GetFullHistoryNodesOnePerShard() ([]*data.NodeData, error)
	GetFullHistoryNodes(shardID uint32) ([]*data.NodeData, error)
	GetAllFullHistoryNodes() ([]*data.NodeData, error)
	GetShardIDs() []uint32
	ComputeShardId(addressBuff []byte) (uint32, error)
	CallGetRestEndPoint(address string, path string, value interface{}) (int, error)
	CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error)
	IsInterfaceNil() bool
}

// ExternalStorageConnector defines what a external storage connector should be able to do
type ExternalStorageConnector interface {
	GetTransactionsByAddress(address string) ([]data.DatabaseTransaction, error)
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	IsInterfaceNil() bool
}

// PrivateKeysLoaderHandler defines what a component which handles loading of the private keys file should do
type PrivateKeysLoaderHandler interface {
	PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error)
}

// HeartbeatCacheHandler will define what a real heartbeat cacher should do
type HeartbeatCacheHandler interface {
	LoadHeartbeats() (*data.HeartbeatResponse, error)
	StoreHeartbeats(hbts *data.HeartbeatResponse) error
	IsInterfaceNil() bool
}

// ValidatorStatisticsCacheHandler will define what a real validator statistics cacher should do
type ValidatorStatisticsCacheHandler interface {
	LoadValStats() (map[string]*data.ValidatorApiResponse, error)
	StoreValStats(valStats map[string]*data.ValidatorApiResponse) error
	IsInterfaceNil() bool
}
