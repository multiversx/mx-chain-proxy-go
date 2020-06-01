package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Facade is the mock implementation of a node's router handler
type Facade struct {
	GetAccountHandler                func(address string) (*data.Account, error)
	GetTransactionsHandler           func(address string) ([]data.DatabaseTransaction, error)
	SendTransactionHandler           func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler  func(txs []*data.Transaction) (data.ResponseMultipleTransactions, error)
	SendUserFundsCalled              func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler            func(query *data.SCQuery) (*vmcommon.VMOutput, error)
	GetHeartbeatDataHandler          func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler       func() (map[string]*data.ValidatorApiResponse, error)
	TransactionCostRequestHandler    func(tx *data.Transaction) (string, error)
	GetShardStatusHandler            func(shardID uint32) (*data.GenericAPIResponse, error)
	GetTransactionStatusHandler      func(txHash string) (string, error)
	GetConfigMetricsHandler          func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler         func(shardID uint32) (*data.GenericAPIResponse, error)
	GetBlockByShardIDAndNonceHandler func(shardID uint32, nonce uint64) (data.ApiBlock, error)
}

// GetNetworkStatusMetrics -
func (f *Facade) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	if f.GetNetworkMetricsHandler != nil {
		return f.GetNetworkMetricsHandler(shardID)
	}

	return nil, nil
}

// GetNetworkConfigMetrics -
func (f *Facade) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if f.GetConfigMetricsHandler != nil {
		return f.GetConfigMetricsHandler()
	}

	return nil, nil
}

// ValidatorStatistics is the mock implementation of a handler's ValidatorStatistics method
func (f *Facade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return f.ValidatorStatisticsHandler()
}

// GetShardStatus --
func (f *Facade) GetShardStatus(shardID uint32) (*data.GenericAPIResponse, error) {
	return f.GetShardStatusHandler(shardID)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// GetTransactions --
func (f *Facade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return f.GetTransactionsHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SendMultipleTransactions is the mock implementation of a handler's SendMultipleTransactions method
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (data.ResponseMultipleTransactions, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// TransactionCostRequest --
func (f *Facade) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return f.TransactionCostRequestHandler(tx)
}

// GetTransactionStatus --
func (f *Facade) GetTransactionStatus(txHash string) (string, error) {
	return f.GetTransactionStatusHandler(txHash)
}

// SendUserFunds is the mock implementation of a handler's SendUserFunds method
func (f *Facade) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery is a mock implementation.
func (f *Facade) ExecuteSCQuery(query *data.SCQuery) (*vmcommon.VMOutput, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData is the mock implementation of a handler's GetHeartbeatData method
func (f *Facade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// GetBlockByShardIDAndNonce -
func (f *Facade) GetBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error) {
	return f.GetBlockByShardIDAndNonceHandler(shardID, nonce)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
