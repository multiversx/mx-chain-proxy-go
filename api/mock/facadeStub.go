package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeStub is the stub implementation of the proxy routes' handler
type FacadeStub struct {
	IsFaucetEnabledHandler                      func() bool
	GetAccountHandler                           func(address string) (*data.Account, error)
	GetShardIDForAddressHandler                 func(address string) (uint32, error)
	GetValueForKeyHandler                       func(address string, key string) (string, error)
	GetKeyValuePairsHandler                     func(address string) (*data.GenericAPIResponse, error)
	GetESDTTokenDataCalled                      func(address string, key string) (*data.GenericAPIResponse, error)
	GetESDTNftTokenDataCalled                   func(address string, key string, nonce uint64) (*data.GenericAPIResponse, error)
	GetESDTsWithRoleCalled                      func(address string, role string) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddressCalled     func(address string) (*data.GenericAPIResponse, error)
	GetAllESDTTokensCalled                      func(address string) (*data.GenericAPIResponse, error)
	GetTransactionsHandler                      func(address string) ([]data.DatabaseTransaction, error)
	GetTransactionHandler                       func(txHash string, withResults bool) (*data.FullTransaction, error)
	SendTransactionHandler                      func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler             func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionHandler                  func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                         func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler                       func(query *data.SCQuery) (*vm.VMOutputApi, error)
	GetHeartbeatDataHandler                     func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler                  func() (map[string]*data.ValidatorApiResponse, error)
	TransactionCostRequestHandler               func(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatusHandler                 func(txHash string, sender string) (string, error)
	GetConfigMetricsHandler                     func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler                    func(shardID uint32) (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsHandler                    func(tokenType string) (*data.GenericAPIResponse, error)
	GetEnableEpochsMetricsHandler               func() (*data.GenericAPIResponse, error)
	GetEconomicsDataMetricsHandler              func() (*data.GenericAPIResponse, error)
	GetDirectStakedInfoCalled                   func() (*data.GenericAPIResponse, error)
	GetDelegatedInfoCalled                      func() (*data.GenericAPIResponse, error)
	GetBlockByShardIDAndNonceHandler            func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetTransactionByHashAndSenderAddressHandler func(txHash string, sndAddr string, withResults bool) (*data.FullTransaction, int, error)
	GetBlockByHashCalled                        func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled                       func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetHyperBlockByHashCalled                   func(hash string) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled                  func(nonce uint64) (*data.HyperblockApiResponse, error)
	ReloadObserversCalled                       func() data.NodesReloadResponse
	ReloadFullHistoryObserversCalled            func() data.NodesReloadResponse
	GetProofCalled                              func(string, string) (*data.GenericAPIResponse, error)
	GetProofDataTrieCalled                      func(string, string, string) (*data.GenericAPIResponse, error)
	GetProofCurrentRootHashCalled               func(string) (*data.GenericAPIResponse, error)
	VerifyProofCalled                           func(string, string, []string) (*data.GenericAPIResponse, error)
	GetESDTsRolesCalled                         func(address string) (*data.GenericAPIResponse, error)
	GetESDTSupplyCalled                         func(token string) (*data.ESDTSupplyResponse, error)
}

// GetProof -
func (f *FacadeStub) GetProof(rootHash string, address string) (*data.GenericAPIResponse, error) {
	if f.GetProofCalled != nil {
		return f.GetProofCalled(rootHash, address)
	}

	return nil, nil
}

// GetProofDataTrie -
func (f *FacadeStub) GetProofDataTrie(rootHash string, address string, key string) (*data.GenericAPIResponse, error) {
	if f.GetProofDataTrieCalled != nil {
		return f.GetProofDataTrieCalled(rootHash, address, key)
	}

	return nil, nil
}

// GetProofCurrentRootHash -
func (f *FacadeStub) GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error) {
	if f.GetProofCurrentRootHashCalled != nil {
		return f.GetProofCurrentRootHashCalled(address)
	}

	return nil, nil
}

// VerifyProof -
func (f *FacadeStub) VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error) {
	if f.VerifyProofCalled != nil {
		return f.VerifyProofCalled(rootHash, address, proof)
	}

	return nil, nil
}

// IsFaucetEnabled -
func (f *FacadeStub) IsFaucetEnabled() bool {
	if f.IsFaucetEnabledHandler != nil {
		return f.IsFaucetEnabledHandler()
	}

	return true
}

// ReloadObservers -
func (f *FacadeStub) ReloadObservers() data.NodesReloadResponse {
	if f.ReloadObserversCalled != nil {
		return f.ReloadObserversCalled()
	}

	return data.NodesReloadResponse{}
}

// ReloadFullHistoryObservers -
func (f *FacadeStub) ReloadFullHistoryObservers() data.NodesReloadResponse {
	if f.ReloadFullHistoryObserversCalled != nil {
		return f.ReloadFullHistoryObserversCalled()
	}

	return data.NodesReloadResponse{}
}

// GetNetworkStatusMetrics -
func (f *FacadeStub) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	if f.GetNetworkMetricsHandler != nil {
		return f.GetNetworkMetricsHandler(shardID)
	}

	return nil, nil
}

// GetNetworkConfigMetrics -
func (f *FacadeStub) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if f.GetConfigMetricsHandler != nil {
		return f.GetConfigMetricsHandler()
	}

	return nil, nil
}

// GetEconomicsDataMetrics -
func (f *FacadeStub) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	if f.GetEconomicsDataMetricsHandler != nil {
		return f.GetEconomicsDataMetricsHandler()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetAllIssuedESDTs -
func (f *FacadeStub) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	if f.GetAllIssuedESDTsHandler != nil {
		return f.GetAllIssuedESDTsHandler(tokenType)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetESDTsWithRole -
func (f *FacadeStub) GetESDTsWithRole(address string, role string) (*data.GenericAPIResponse, error) {
	if f.GetESDTsWithRoleCalled != nil {
		return f.GetESDTsWithRoleCalled(address, role)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetESDTsRoles -
func (f *FacadeStub) GetESDTsRoles(address string) (*data.GenericAPIResponse, error) {
	if f.GetESDTsRolesCalled != nil {
		return f.GetESDTsRolesCalled(address)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetNFTTokenIDsRegisteredByAddress -
func (f *FacadeStub) GetNFTTokenIDsRegisteredByAddress(address string) (*data.GenericAPIResponse, error) {
	if f.GetNFTTokenIDsRegisteredByAddressCalled != nil {
		return f.GetNFTTokenIDsRegisteredByAddressCalled(address)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetDirectStakedInfo -
func (f *FacadeStub) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	if f.GetDirectStakedInfoCalled != nil {
		return f.GetDirectStakedInfoCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetDelegatedInfo -
func (f *FacadeStub) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	if f.GetDelegatedInfoCalled != nil {
		return f.GetDelegatedInfoCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetEnableEpochsMetrics -
func (f *FacadeStub) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	return f.GetEnableEpochsMetricsHandler()
}

// GetESDTSupply -
func (f *FacadeStub) GetESDTSupply(token string) (*data.ESDTSupplyResponse, error) {
	if f.GetESDTSupplyCalled != nil {
		return f.GetESDTSupplyCalled(token)
	}

	return nil, nil
}

// ValidatorStatistics -
func (f *FacadeStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return f.ValidatorStatisticsHandler()
}

// GetAccount -
func (f *FacadeStub) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// GetKeyValuePairs -
func (f *FacadeStub) GetKeyValuePairs(address string) (*data.GenericAPIResponse, error) {
	return f.GetKeyValuePairsHandler(address)
}

// GetValueForKey -
func (f *FacadeStub) GetValueForKey(address string, key string) (string, error) {
	return f.GetValueForKeyHandler(address, key)
}

// GetShardIDForAddress -
func (f *FacadeStub) GetShardIDForAddress(address string) (uint32, error) {
	return f.GetShardIDForAddressHandler(address)
}

// GetESDTTokenData -
func (f *FacadeStub) GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error) {
	if f.GetESDTTokenDataCalled != nil {
		return f.GetESDTTokenDataCalled(address, key)
	}

	return nil, nil
}

// GetAllESDTTokens -
func (f *FacadeStub) GetAllESDTTokens(address string) (*data.GenericAPIResponse, error) {
	if f.GetAllESDTTokensCalled != nil {
		return f.GetAllESDTTokensCalled(address)
	}

	return nil, nil
}

// GetESDTNftTokenData -
func (f *FacadeStub) GetESDTNftTokenData(address string, key string, nonce uint64) (*data.GenericAPIResponse, error) {
	if f.GetESDTNftTokenDataCalled != nil {
		return f.GetESDTNftTokenDataCalled(address, key, nonce)
	}

	return nil, nil
}

// GetTransactions -
func (f *FacadeStub) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return f.GetTransactionsHandler(address)
}

// GetTransactionByHashAndSenderAddress -
func (f *FacadeStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error) {
	return f.GetTransactionByHashAndSenderAddressHandler(txHash, sndAddr, withEvents)
}

// GetTransaction -
func (f *FacadeStub) GetTransaction(txHash string, withResults bool) (*data.FullTransaction, error) {
	return f.GetTransactionHandler(txHash, withResults)
}

// SendTransaction -
func (f *FacadeStub) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SimulateTransaction -
func (f *FacadeStub) SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
	return f.SimulateTransactionHandler(tx, checkSignature)
}

// GetAddressConverter -
func (f *FacadeStub) GetAddressConverter() (core.PubkeyConverter, error) {
	return nil, nil
}

// SendMultipleTransactions -
func (f *FacadeStub) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// TransactionCostRequest -
func (f *FacadeStub) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	return f.TransactionCostRequestHandler(tx)
}

// GetTransactionStatus -
func (f *FacadeStub) GetTransactionStatus(txHash string, sender string) (string, error) {
	return f.GetTransactionStatusHandler(txHash, sender)
}

// SendUserFunds -
func (f *FacadeStub) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery -
func (f *FacadeStub) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData -
func (f *FacadeStub) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// GetAtlasBlockByShardIDAndNonce -
func (f *FacadeStub) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return f.GetBlockByShardIDAndNonceHandler(shardID, nonce)
}

// GetBlockByHash -
func (f *FacadeStub) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	return f.GetBlockByHashCalled(shardID, hash, withTxs)
}

// GetBlockByNonce -
func (f *FacadeStub) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	return f.GetBlockByNonceCalled(shardID, nonce, withTxs)
}

// GetHyperBlockByHash -
func (f *FacadeStub) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByHashCalled(hash)
}

// GetHyperBlockByNonce -
func (f *FacadeStub) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByNonceCalled(nonce)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
