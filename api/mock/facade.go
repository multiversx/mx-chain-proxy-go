package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node's router handler
type Facade struct {
	IsFaucetEnabledHandler                      func() bool
	GetAccountHandler                           func(address string) (*data.Account, int, error)
	GetShardIDForAddressHandler                 func(address string) (uint32, int, error)
	GetValueForKeyHandler                       func(address string, key string) (string, int, error)
	GetESDTTokenDataCalled                      func(address string, key string) (*data.GenericAPIResponse, int, error)
	GetAllESDTTokensCalled                      func(address string) (*data.GenericAPIResponse, int, error)
	GetTransactionsHandler                      func(address string) ([]data.DatabaseTransaction, int, error)
	GetTransactionHandler                       func(txHash string, withResults bool) (*data.FullTransaction, int, error)
	SendTransactionHandler                      func(tx *data.Transaction) (string, int, error)
	SendMultipleTransactionsHandler             func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, int, error)
	SimulateTransactionHandler                  func(tx *data.Transaction) (*data.GenericAPIResponse, int, error)
	SendUserFundsCalled                         func(receiver string, value *big.Int) (int, error)
	ExecuteSCQueryHandler                       func(query *data.SCQuery) (*vm.VMOutputApi, int, error)
	GetHeartbeatDataHandler                     func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler                  func() (map[string]*data.ValidatorApiResponse, error)
	TransactionCostRequestHandler               func(tx *data.Transaction) (string, int, error)
	GetTransactionStatusHandler                 func(txHash string, sender string) (string, int, error)
	GetConfigMetricsHandler                     func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler                    func(shardID uint32) (*data.GenericAPIResponse, int, error)
	GetEconomicsDataMetricsHandler              func() (*data.GenericAPIResponse, error)
	GetBlockByShardIDAndNonceHandler            func(shardID uint32, nonce uint64) (data.AtlasBlock, int, error)
	GetTransactionByHashAndSenderAddressHandler func(txHash string, sndAddr string, withResults bool) (*data.FullTransaction, int, error)
	GetBlockByHashCalled                        func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, int, error)
	GetBlockByNonceCalled                       func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, int, error)
	GetHyperBlockByHashCalled                   func(hash string) (*data.HyperblockApiResponse, int, error)
	GetHyperBlockByNonceCalled                  func(nonce uint64) (*data.HyperblockApiResponse, int, error)
}

// IsFaucetEnabled -
func (f *Facade) IsFaucetEnabled() bool {
	if f.IsFaucetEnabledHandler != nil {
		return f.IsFaucetEnabledHandler()
	}

	return true
}

// GetNetworkStatusMetrics -
func (f *Facade) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, int, error) {
	if f.GetNetworkMetricsHandler != nil {
		return f.GetNetworkMetricsHandler(shardID)
	}

	return nil, 0, nil
}

// GetNetworkConfigMetrics -
func (f *Facade) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if f.GetConfigMetricsHandler != nil {
		return f.GetConfigMetricsHandler()
	}

	return nil, nil
}

// GetEconomicsDataMetrics -
func (f *Facade) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	if f.GetEconomicsDataMetricsHandler != nil {
		return f.GetEconomicsDataMetricsHandler()
	}

	return &data.GenericAPIResponse{}, nil
}

// ValidatorStatistics -
func (f *Facade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return f.ValidatorStatisticsHandler()
}

// GetAccount -
func (f *Facade) GetAccount(address string) (*data.Account, int, error) {
	return f.GetAccountHandler(address)
}

// GetValueForKey -
func (f *Facade) GetValueForKey(address string, key string) (string, int, error) {
	return f.GetValueForKeyHandler(address, key)
}

// GetShardIDForAddress -
func (f *Facade) GetShardIDForAddress(address string) (uint32, int, error) {
	return f.GetShardIDForAddressHandler(address)
}

// GetESDTTokenData -
func (f *Facade) GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, int, error) {
	if f.GetESDTTokenDataCalled != nil {
		return f.GetESDTTokenDataCalled(address, key)
	}

	return nil, 0, nil
}

// GetAllESDTTokens -
func (f *Facade) GetAllESDTTokens(address string) (*data.GenericAPIResponse, int, error) {
	if f.GetAllESDTTokensCalled != nil {
		return f.GetAllESDTTokensCalled(address)
	}

	return nil, 0, nil
}

// GetTransactions -
func (f *Facade) GetTransactions(address string) ([]data.DatabaseTransaction, int, error) {
	return f.GetTransactionsHandler(address)
}

// GetTransactionByHashAndSenderAddress -
func (f *Facade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error) {
	return f.GetTransactionByHashAndSenderAddressHandler(txHash, sndAddr, withEvents)
}

// GetTransaction -
func (f *Facade) GetTransaction(txHash string, withResults bool) (*data.FullTransaction, int, error) {
	return f.GetTransactionHandler(txHash, withResults)
}

// SendTransaction -
func (f *Facade) SendTransaction(tx *data.Transaction) (string, int, error) {
	return f.SendTransactionHandler(tx)
}

// SimulateTransaction -
func (f *Facade) SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, int, error) {
	return f.SimulateTransactionHandler(tx)
}

// GetAddressConverter -
func (f *Facade) GetAddressConverter() (core.PubkeyConverter, error) {
	return nil, nil
}

// SendMultipleTransactions -
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, int, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// TransactionCostRequest -
func (f *Facade) TransactionCostRequest(tx *data.Transaction) (string, int, error) {
	return f.TransactionCostRequestHandler(tx)
}

// GetTransactionStatus -
func (f *Facade) GetTransactionStatus(txHash string, sender string) (string, int, error) {
	return f.GetTransactionStatusHandler(txHash, sender)
}

// SendUserFunds -
func (f *Facade) SendUserFunds(receiver string, value *big.Int) (int, error) {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery -
func (f *Facade) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, int, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData -
func (f *Facade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// GetAtlasBlockByShardIDAndNonce -
func (f *Facade) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, int, error) {
	return f.GetBlockByShardIDAndNonceHandler(shardID, nonce)
}

// GetBlockByHash -
func (f *Facade) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, int, error) {
	return f.GetBlockByHashCalled(shardID, hash, withTxs)
}

// GetBlockByNonce -
func (f *Facade) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, int, error) {
	return f.GetBlockByNonceCalled(shardID, nonce, withTxs)
}

// GetHyperBlockByHash -
func (f *Facade) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, int, error) {
	return f.GetHyperBlockByHashCalled(hash)
}

// GetHyperBlockByNonce -
func (f *Facade) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, int, error) {
	return f.GetHyperBlockByNonceCalled(nonce)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
