package process

import (
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer"
)

// Processor defines what a processor should be able to do
type Processor interface {
	GetObservers(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllObservers(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetObserversOnePerShard(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetFullHistoryNodesOnePerShard(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetFullHistoryNodes(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllFullHistoryNodes(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetShardIDs() []uint32
	ComputeShardId(addressBuff []byte) (uint32, error)
	CallGetRestEndPoint(address string, path string, value interface{}) (int, error)
	CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error)
	GetShardCoordinator() common.Coordinator
	GetPubKeyConverter() core.PubkeyConverter
	GetObserverProvider() observer.NodesProviderHandler
	GetFullHistoryNodesProvider() observer.NodesProviderHandler
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

// GenericApiResponseCacheHandler will define what a real economic metrics cacher should do
type GenericApiResponseCacheHandler interface {
	Load() (*data.GenericAPIResponse, error)
	Store(response *data.GenericAPIResponse)
	IsInterfaceNil() bool
}

// TransactionCostHandler will define what a real transaction cost handler should do
type TransactionCostHandler interface {
	ResolveCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error)
}

// LogsMergerHandler will define what a real merge logs handler should do
type LogsMergerHandler interface {
	MergeLogEvents(logSource *transaction.ApiLogs, logDestination *transaction.ApiLogs) *transaction.ApiLogs
	IsInterfaceNil() bool
}

// SCQueryService defines how data should be get from a SC account
type SCQueryService interface {
	ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error)
	IsInterfaceNil() bool
}

// StatusMetricsProvider defines what a status metrics provider should do
type StatusMetricsProvider interface {
	GetAll() map[string]*data.EndpointMetrics
	GetMetricsForPrometheus() string
	IsInterfaceNil() bool
}

// HttpClient defines an interface for the http client
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
