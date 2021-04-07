package groups

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountsFacadeHandler interface defines methods that can be used from facade context variable
type AccountsFacadeHandler interface {
	GetAccount(address string) (*data.Account, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
	GetShardIDForAddress(address string) (uint32, error)
	GetValueForKey(address string, key string) (string, error)
	GetAllESDTTokens(address string) (*data.GenericAPIResponse, error)
	GetKeyValuePairs(address string) (*data.GenericAPIResponse, error)
	GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error)
	GetESDTNftTokenData(address string, key string, nonce uint64) (*data.GenericAPIResponse, error)
}

// BlocksFacadeHandler interface defines methods that can be used from facade context variable
type BlocksFacadeHandler interface {
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
}

// BlockAtlasFacadeHandler interface defines methods that can be used from facade context variable
type BlockAtlasFacadeHandler interface {
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error)
}

// HyperBlockFacadeHandler defines the actions needed for fetching the hyperblocks from the nodes
type HyperBlockFacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error)
}

// NetworkFacadeHandler interface defines methods that can be used from facade context variable
type NetworkFacadeHandler interface {
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error)
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetEconomicsDataMetrics() (*data.GenericAPIResponse, error)
	GetAllIssuedESDTs() (*data.GenericAPIResponse, error)
}

// NodeFacadeHandler interface defines methods that can be used from facade context variable
type NodeFacadeHandler interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
}

// TransactionFacadeHandler interface defines methods that can be used from facade context variable
type TransactionFacadeHandler interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	IsFaucetEnabled() bool
	SendUserFunds(receiver string, value *big.Int) error
	TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatus(txHash string, sender string) (string, error)
	GetTransaction(txHash string, withResults bool) (*data.FullTransaction, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error)
}

// ProofFacadeHandler interface defines methods that can be used from facade context variable
type ProofFacadeHandler interface {
	GetProof(rootHash []byte, address []byte) ([][]byte, error)
	VerifyProof(rootHash []byte, address []byte, proof [][]byte) (bool, error)
}

// ValidatorFacadeHandler interface defines methods that can be used from facade context variable
type ValidatorFacadeHandler interface {
	ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error)
}

// VmValuesFacadeHandler interface defines methods that can be used from `elrondFacade` context variable
type VmValuesFacadeHandler interface {
	ExecuteSCQuery(*data.SCQuery) (*vm.VMOutputApi, error)
}

// ActionsFacadeHandler interface defines methods that can be used from facade context variable
type ActionsFacadeHandler interface {
	ReloadObservers() data.NodesReloadResponse
	ReloadFullHistoryObservers() data.NodesReloadResponse
}
