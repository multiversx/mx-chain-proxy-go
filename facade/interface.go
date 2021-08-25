package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ActionsProcessor defines what an actions processor should do
type ActionsProcessor interface {
	ReloadObservers() data.NodesReloadResponse
	ReloadFullHistoryObservers() data.NodesReloadResponse
}

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
	GetShardIDForAddress(address string) (uint32, error)
	GetValueForKey(address string, key string) (string, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
	GetAllESDTTokens(address string) (*data.GenericAPIResponse, error)
	GetKeyValuePairs(address string) (*data.GenericAPIResponse, error)
	GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error)
	GetESDTsWithRole(address string, role string) (*data.GenericAPIResponse, error)
	GetESDTsRoles(address string) (*data.GenericAPIResponse, error)
	GetESDTNftTokenData(address string, key string, nonce uint64) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddress(address string) (*data.GenericAPIResponse, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatus(txHash string, sender string) (string, error)
	GetTransaction(txHash string, withEvents bool) (*data.FullTransaction, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error)
	ComputeTransactionHash(tx *data.Transaction) (string, error)
}

// ProofProcessor defines what a proof request processor should do
type ProofProcessor interface {
	GetProof(rootHash string, address string) (*data.GenericAPIResponse, error)
	GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error)
	VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error)
}

// SCQueryService defines how data should be get from a SC account
type SCQueryService interface {
	ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, error)
}

// HeartbeatProcessor defines what a heartbeat processor should do
type HeartbeatProcessor interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
}

// ValidatorStatisticsProcessor defines what a validator statistics processor should do
type ValidatorStatisticsProcessor interface {
	GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error)
}

// ESDTSuppliesProcessor defines what an esdt supplies processor should do
type ESDTSuppliesProcessor interface {
	GetESDTSupply(token string) (*data.ESDTSupplyResponse, error)
}

// NodeStatusProcessor defines what a node status processor should do
type NodeStatusProcessor interface {
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error)
	GetEconomicsDataMetrics() (*data.GenericAPIResponse, error)
	GetLatestFullySynchronizedHyperblockNonce() (uint64, error)
	GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error)
	GetEnableEpochsMetrics() (*data.GenericAPIResponse, error)
	GetDirectStakedInfo() (*data.GenericAPIResponse, error)
	GetDelegatedInfo() (*data.GenericAPIResponse, error)
}

// BlockProcessor defines what a block processor should do
type BlockProcessor interface {
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error)
}

// FaucetProcessor defines what a component which will handle faucets should do
type FaucetProcessor interface {
	IsEnabled() bool
	SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error)
	GenerateTxForSendUserFunds(
		senderSk crypto.PrivateKey,
		senderPk string,
		senderNonce uint64,
		receiver string,
		value *big.Int,
		chainID string,
		version uint32,
	) (*data.Transaction, error)
}
