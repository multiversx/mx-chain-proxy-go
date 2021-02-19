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
	GetAccountHandler                           func(address string) (*data.Account, error)
	GetShardIDForAddressHandler                 func(address string) (uint32, error)
	GetValueForKeyHandler                       func(address string, key string) (string, error)
	GetKeyValuePairsHandler                     func(address string) (*data.GenericAPIResponse, error)
	GetESDTTokenDataCalled                      func(address string, key string) (*data.GenericAPIResponse, error)
	GetAllESDTTokensCalled                      func(address string) (*data.GenericAPIResponse, error)
	GetTransactionsHandler                      func(address string) ([]data.DatabaseTransaction, error)
	GetTransactionHandler                       func(txHash string, withResults bool) (*data.FullTransaction, error)
	SendTransactionHandler                      func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler             func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionHandler                  func(tx *data.Transaction) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                         func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler                       func(query *data.SCQuery) (*vm.VMOutputApi, error)
	GetHeartbeatDataHandler                     func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler                  func() (map[string]*data.ValidatorApiResponse, error)
	TransactionCostRequestHandler               func(tx *data.Transaction) (string, error)
	GetTransactionStatusHandler                 func(txHash string, sender string) (string, error)
	GetConfigMetricsHandler                     func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler                    func(shardID uint32) (*data.GenericAPIResponse, error)
	GetEconomicsDataMetricsHandler              func() (*data.GenericAPIResponse, error)
	GetBlockByShardIDAndNonceHandler            func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetTransactionByHashAndSenderAddressHandler func(txHash string, sndAddr string, withResults bool) (*data.FullTransaction, int, error)
	GetBlockByHashCalled                        func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled                       func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetHyperBlockByHashCalled                   func(hash string) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled                  func(nonce uint64) (*data.HyperblockApiResponse, error)
	ReloadObserversCalled                       func() data.NodesReloadResponse
	ReloadFullHistoryObserversCalled            func() data.NodesReloadResponse
	GetTotalStakedCalled                        func() (*data.GenericAPIResponse, error)
}

// IsFaucetEnabled -
func (f *Facade) IsFaucetEnabled() bool {
	if f.IsFaucetEnabledHandler != nil {
		return f.IsFaucetEnabledHandler()
	}

	return true
}

// ReloadObservers -
func (f *Facade) ReloadObservers() data.NodesReloadResponse {
	if f.ReloadObserversCalled != nil {
		return f.ReloadObserversCalled()
	}

	return data.NodesReloadResponse{}
}

// ReloadFullHistoryObservers -
func (f *Facade) ReloadFullHistoryObservers() data.NodesReloadResponse {
	if f.ReloadFullHistoryObserversCalled != nil {
		return f.ReloadFullHistoryObserversCalled()
	}

	return data.NodesReloadResponse{}
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
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// GetKeyValuePairs -
func (f *Facade) GetKeyValuePairs(address string) (*data.GenericAPIResponse, error) {
	return f.GetKeyValuePairsHandler(address)
}

// GetValueForKey -
func (f *Facade) GetValueForKey(address string, key string) (string, error) {
	return f.GetValueForKeyHandler(address, key)
}

// GetShardIDForAddress -
func (f *Facade) GetShardIDForAddress(address string) (uint32, error) {
	return f.GetShardIDForAddressHandler(address)
}

// GetESDTTokenData -
func (f *Facade) GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error) {
	if f.GetESDTTokenDataCalled != nil {
		return f.GetESDTTokenDataCalled(address, key)
	}

	return nil, nil
}

// GetAllESDTTokens -
func (f *Facade) GetAllESDTTokens(address string) (*data.GenericAPIResponse, error) {
	if f.GetAllESDTTokensCalled != nil {
		return f.GetAllESDTTokensCalled(address)
	}

	return nil, nil
}

// GetTransactions -
func (f *Facade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return f.GetTransactionsHandler(address)
}

// GetTransactionByHashAndSenderAddress -
func (f *Facade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error) {
	return f.GetTransactionByHashAndSenderAddressHandler(txHash, sndAddr, withEvents)
}

// GetTransaction -
func (f *Facade) GetTransaction(txHash string, withResults bool) (*data.FullTransaction, error) {
	return f.GetTransactionHandler(txHash, withResults)
}

// SendTransaction -
func (f *Facade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SimulateTransaction -
func (f *Facade) SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, error) {
	return f.SimulateTransactionHandler(tx)
}

// GetAddressConverter -
func (f *Facade) GetAddressConverter() (core.PubkeyConverter, error) {
	return nil, nil
}

// SendMultipleTransactions -
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// TransactionCostRequest -
func (f *Facade) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return f.TransactionCostRequestHandler(tx)
}

// GetTransactionStatus -
func (f *Facade) GetTransactionStatus(txHash string, sender string) (string, error) {
	return f.GetTransactionStatusHandler(txHash, sender)
}

// SendUserFunds -
func (f *Facade) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery -
func (f *Facade) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData -
func (f *Facade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// GetAtlasBlockByShardIDAndNonce -
func (f *Facade) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return f.GetBlockByShardIDAndNonceHandler(shardID, nonce)
}

// GetBlockByHash -
func (f *Facade) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	return f.GetBlockByHashCalled(shardID, hash, withTxs)
}

// GetBlockByNonce -
func (f *Facade) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	return f.GetBlockByNonceCalled(shardID, nonce, withTxs)
}

// GetHyperBlockByHash -
func (f *Facade) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByHashCalled(hash)
}

// GetHyperBlockByNonce -
func (f *Facade) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByNonceCalled(nonce)
}

// GetTotalStaked -
func (f *Facade) GetTotalStaked() (*data.GenericAPIResponse, error) {
	if f.GetTotalStakedCalled != nil {
		return f.GetTotalStakedCalled()
	}

	return nil, nil
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
