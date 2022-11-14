package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node's router handler
type Facade struct {
	IsFaucetEnabledHandler                       func() bool
	GetAccountHandler                            func(address string, options common.AccountQueryOptions) (*data.AccountModel, error)
	GetShardIDForAddressHandler                  func(address string) (uint32, error)
	GetValueForKeyHandler                        func(address string, key string, options common.AccountQueryOptions) (string, error)
	GetKeyValuePairsHandler                      func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTTokenDataCalled                       func(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTNftTokenDataCalled                    func(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsWithRoleCalled                       func(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddressCalled      func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetAllESDTTokensCalled                       func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetTransactionsHandler                       func(address string) ([]data.DatabaseTransaction, error)
	GetTransactionHandler                        func(txHash string, withResults bool) (*transaction.ApiTransactionResult, error)
	GetTransactionsPoolHandler                   func(fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForShardHandler           func(shardID uint32, fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForSenderHandler          func(sender, fields string) (*data.TransactionsPoolForSender, error)
	GetLastPoolNonceForSenderHandler             func(sender string) (uint64, error)
	GetTransactionsPoolNonceGapsForSenderHandler func(sender string) (*data.TransactionsPoolNonceGaps, error)
	SendTransactionHandler                       func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler              func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionHandler                   func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                          func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler                        func(query *data.SCQuery) (*vm.VMOutputApi, error)
	GetHeartbeatDataHandler                      func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler                   func() (map[string]*data.ValidatorApiResponse, error)
	TransactionCostRequestHandler                func(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatusHandler                  func(txHash string, sender string) (string, error)
	GetConfigMetricsHandler                      func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler                     func(shardID uint32) (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsHandler                     func(tokenType string) (*data.GenericAPIResponse, error)
	GetEnableEpochsMetricsHandler                func() (*data.GenericAPIResponse, error)
	GetEconomicsDataMetricsHandler               func() (*data.GenericAPIResponse, error)
	GetDirectStakedInfoCalled                    func() (*data.GenericAPIResponse, error)
	GetDelegatedInfoCalled                       func() (*data.GenericAPIResponse, error)
	GetRatingsConfigCalled                       func() (*data.GenericAPIResponse, error)
	GetBlockByShardIDAndNonceHandler             func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetTransactionByHashAndSenderAddressHandler  func(txHash string, sndAddr string, withResults bool) (*transaction.ApiTransactionResult, int, error)
	GetBlockByHashCalled                         func(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled                        func(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlocksByRoundCalled                       func(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error)
	GetInternalBlockByHashCalled                 func(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonceCalled                func(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalMiniBlockByHashCalled             func(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error)
	GetInternalStartOfEpochMetaBlockCalled       func(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetHyperBlockByHashCalled                    func(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled                   func(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	ReloadObserversCalled                        func() data.NodesReloadResponse
	ReloadFullHistoryObserversCalled             func() data.NodesReloadResponse
	GetProofCalled                               func(string, string) (*data.GenericAPIResponse, error)
	GetProofCurrentRootHashCalled                func(string) (*data.GenericAPIResponse, error)
	VerifyProofCalled                            func(string, string, []string) (*data.GenericAPIResponse, error)
	GetESDTsRolesCalled                          func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTSupplyCalled                          func(token string) (*data.ESDTSupplyResponse, error)
	GetMetricsCalled                             func() map[string]*data.EndpointMetrics
	GetPrometheusMetricsCalled                   func() string
	GetGenesisNodesPubKeysCalled                 func() (*data.GenericAPIResponse, error)
	GetGasConfigsCalled                          func() (*data.GenericAPIResponse, error)
	IsOldStorageForTokenCalled                   func(tokenID string, nonce uint64) (bool, error)
	GetAboutInfoCalled                           func() (*data.GenericAPIResponse, error)
	GetTriesStatisticsCalled                     func(shardID uint32) (*data.TrieStatisticsAPIResponse, error)
	GetEpochStartDataCalled                      func(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error)
}

// GetProof -
func (f *Facade) GetProof(rootHash string, address string) (*data.GenericAPIResponse, error) {
	if f.GetProofCalled != nil {
		return f.GetProofCalled(rootHash, address)
	}

	return nil, nil
}

// GetProofCurrentRootHash -
func (f *Facade) GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error) {
	if f.GetProofCurrentRootHashCalled != nil {
		return f.GetProofCurrentRootHashCalled(address)
	}

	return nil, nil
}

// VerifyProof -
func (f *Facade) VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error) {
	if f.VerifyProofCalled != nil {
		return f.VerifyProofCalled(rootHash, address, proof)
	}

	return nil, nil
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

// GetAllIssuedESDTs -
func (f *Facade) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	if f.GetAllIssuedESDTsHandler != nil {
		return f.GetAllIssuedESDTsHandler(tokenType)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetESDTsWithRole -
func (f *Facade) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTsWithRoleCalled != nil {
		return f.GetESDTsWithRoleCalled(address, role, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetESDTsRoles -
func (f *Facade) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTsRolesCalled != nil {
		return f.GetESDTsRolesCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetNFTTokenIDsRegisteredByAddress -
func (f *Facade) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetNFTTokenIDsRegisteredByAddressCalled != nil {
		return f.GetNFTTokenIDsRegisteredByAddressCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetDirectStakedInfo -
func (f *Facade) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	if f.GetDirectStakedInfoCalled != nil {
		return f.GetDirectStakedInfoCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetDelegatedInfo -
func (f *Facade) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	if f.GetDelegatedInfoCalled != nil {
		return f.GetDelegatedInfoCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetEnableEpochsMetrics -
func (f *Facade) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	return f.GetEnableEpochsMetricsHandler()
}

// GetRatingsConfig -
func (f *Facade) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	return f.GetRatingsConfigCalled()
}

// GetESDTSupply -
func (f *Facade) GetESDTSupply(token string) (*data.ESDTSupplyResponse, error) {
	if f.GetESDTSupplyCalled != nil {
		return f.GetESDTSupplyCalled(token)
	}

	return nil, nil
}

// ValidatorStatistics -
func (f *Facade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return f.ValidatorStatisticsHandler()
}

// GetAccount -
func (f *Facade) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	return f.GetAccountHandler(address, options)
}

// GetKeyValuePairs -
func (f *Facade) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return f.GetKeyValuePairsHandler(address, options)
}

// GetValueForKey -
func (f *Facade) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	return f.GetValueForKeyHandler(address, key, options)
}

// GetShardIDForAddress -
func (f *Facade) GetShardIDForAddress(address string) (uint32, error) {
	return f.GetShardIDForAddressHandler(address)
}

// GetESDTTokenData -
func (f *Facade) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTTokenDataCalled != nil {
		return f.GetESDTTokenDataCalled(address, key, options)
	}

	return nil, nil
}

// GetAllESDTTokens -
func (f *Facade) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetAllESDTTokensCalled != nil {
		return f.GetAllESDTTokensCalled(address, options)
	}

	return nil, nil
}

// GetESDTNftTokenData -
func (f *Facade) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTNftTokenDataCalled != nil {
		return f.GetESDTNftTokenDataCalled(address, key, nonce, options)
	}

	return nil, nil
}

// IsOldStorageForToken -
func (f *Facade) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	if f.IsOldStorageForTokenCalled != nil {
		return f.IsOldStorageForTokenCalled(tokenID, nonce)
	}

	return false, nil
}

// GetTransactions -
func (f *Facade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return f.GetTransactionsHandler(address)
}

// GetTransactionByHashAndSenderAddress -
func (f *Facade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error) {
	return f.GetTransactionByHashAndSenderAddressHandler(txHash, sndAddr, withEvents)
}

// GetTransaction -
func (f *Facade) GetTransaction(txHash string, withResults bool) (*transaction.ApiTransactionResult, error) {
	return f.GetTransactionHandler(txHash, withResults)
}

// GetTransactionsPool -
func (f *Facade) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	if f.GetTransactionsPoolHandler != nil {
		return f.GetTransactionsPoolHandler(fields)
	}

	return nil, nil
}

// GetTransactionsPoolForShard -
func (f *Facade) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	if f.GetTransactionsPoolForShardHandler != nil {
		return f.GetTransactionsPoolForShardHandler(shardID, fields)
	}

	return nil, nil
}

// GetTransactionsPoolForSender -
func (f *Facade) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	if f.GetTransactionsPoolForSenderHandler != nil {
		return f.GetTransactionsPoolForSenderHandler(sender, fields)
	}

	return nil, nil
}

// GetLastPoolNonceForSender -
func (f *Facade) GetLastPoolNonceForSender(sender string) (uint64, error) {
	if f.GetLastPoolNonceForSenderHandler != nil {
		return f.GetLastPoolNonceForSenderHandler(sender)
	}

	return 0, nil
}

// GetTransactionsPoolNonceGapsForSender -
func (f *Facade) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	if f.GetTransactionsPoolNonceGapsForSenderHandler != nil {
		return f.GetTransactionsPoolNonceGapsForSenderHandler(sender)
	}

	return nil, nil
}

// SendTransaction -
func (f *Facade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SimulateTransaction -
func (f *Facade) SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
	return f.SimulateTransactionHandler(tx, checkSignature)
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
func (f *Facade) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
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
func (f *Facade) GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return f.GetBlockByHashCalled(shardID, hash, options)
}

// GetBlockByNonce -
func (f *Facade) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return f.GetBlockByNonceCalled(shardID, nonce, options)
}

// GetBlocksByRound -
func (f *Facade) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	if f.GetBlocksByRoundCalled != nil {
		return f.GetBlocksByRoundCalled(round, options)
	}
	return nil, nil
}

// GetInternalBlockByHash -
func (f *Facade) GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return f.GetInternalBlockByHashCalled(shardID, hash, format)
}

// GetInternalBlockByNonce -
func (f *Facade) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return f.GetInternalBlockByNonceCalled(shardID, nonce, format)
}

// GetInternalMiniBlockByHash -
func (f *Facade) GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
	return f.GetInternalMiniBlockByHashCalled(shardID, hash, epoch, format)
}

// GetInternalStartOfEpochMetaBlock -
func (f *Facade) GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return f.GetInternalStartOfEpochMetaBlockCalled(epoch, format)
}

// GetHyperBlockByHash -
func (f *Facade) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByHashCalled(hash, options)
}

// GetHyperBlockByNonce -
func (f *Facade) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByNonceCalled(nonce, options)
}

// GetMetrics -
func (f *Facade) GetMetrics() map[string]*data.EndpointMetrics {
	return f.GetMetricsCalled()
}

// GetMetricsForPrometheus -
func (f *Facade) GetMetricsForPrometheus() string {
	return f.GetPrometheusMetricsCalled()
}

// GetGenesisNodesPubKeys -
func (f *Facade) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	return f.GetGenesisNodesPubKeysCalled()
}

// GetGasConfigs -
func (f *Facade) GetGasConfigs() (*data.GenericAPIResponse, error) {
	return f.GetGasConfigsCalled()
}

// GetAboutInfo -
func (f *Facade) GetAboutInfo() (*data.GenericAPIResponse, error) {
	return f.GetAboutInfoCalled()
}

// GetTriesStatistics -
func (f *Facade) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	if f.GetTriesStatisticsCalled != nil {
		return f.GetTriesStatisticsCalled(shardID)
	}
	return &data.TrieStatisticsAPIResponse{}, nil
}

// GetEpochStartData -
func (f *Facade) GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error) {
	return f.GetEpochStartDataCalled(epoch, shardID)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
