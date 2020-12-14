package groups

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountsFacadeHandler interface defines methods that can be used from facade context variable
type AccountsFacadeHandler interface {
	GetAccount(address string) (*data.Account, int, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, int, error)
	GetShardIDForAddress(address string) (uint32, int, error)
	GetValueForKey(address string, key string) (string, int, error)
	GetAllESDTTokens(address string) (*data.GenericAPIResponse, int, error)
	GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, int, error)
}

// BlocksFacadeHandler interface defines methods that can be used from facade context variable
type BlocksFacadeHandler interface {
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, int, error)
	GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, int, error)
}

// BlockAtlasFacadeHandler interface defines methods that can be used from facade context variable
type BlockAtlasFacadeHandler interface {
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, int, error)
}

// HyperBlockFacadeHandler defines the actions needed for fetching the hyperblocks from the nodes
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, int, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, int, error)
}

// NetworkFacadeHandler interface defines methods that can be used from facade context variable
type NetworkFacadeHandler interface {
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, int, error)
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetEconomicsDataMetrics() (*data.GenericAPIResponse, error)
}

// NodeFacadeHandler interface defines methods that can be used from facade context variable
type NodeFacadeHandler interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
}

// TransactionFacadeHandler interface defines methods that can be used from facade context variable
type TransactionFacadeHandler interface {
	SendTransaction(tx *data.Transaction) (string, int, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, int, error)
	SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, int, error)
	IsFaucetEnabled() bool
	SendUserFunds(receiver string, value *big.Int) (int, error)
	TransactionCostRequest(tx *data.Transaction) (string, int, error)
	GetTransactionStatus(txHash string, sender string) (string, int, error)
	GetTransaction(txHash string, withResults bool) (*data.FullTransaction, int, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error)
}

// ValidatorFacadeHandler interface defines methods that can be used from facade context variable
type ValidatorFacadeHandler interface {
	ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error)
}

// VmValuesFacadeHandler interface defines methods that can be used from `elrondFacade` context variable
type VmValuesFacadeHandler interface {
	ExecuteSCQuery(*data.SCQuery) (*vm.VMOutputApi, int, error)
}
