package mock

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// FacadeStub is the mock implementation of a node's router handler
type FacadeStub struct {
	IsFaucetEnabledHandler                       func() bool
	GetAccountHandler                            func(address string, options common.AccountQueryOptions) (*data.AccountModel, error)
	GetAccountsHandler                           func(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error)
	GetShardIDForAddressHandler                  func(address string) (uint32, error)
	GetValueForKeyHandler                        func(address string, key string, options common.AccountQueryOptions) (string, error)
	GetKeyValuePairsHandler                      func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTTokenDataCalled                       func(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTNftTokenDataCalled                    func(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsWithRoleCalled                       func(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddressCalled      func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetAllESDTTokensCalled                       func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetTransactionsHandler                       func(address string) ([]data.DatabaseTransaction, error)
	GetTransactionHandler                        func(txHash string, withResults bool, relayedTxHash string) (*transaction.ApiTransactionResult, error)
	GetTransactionsPoolHandler                   func(fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForShardHandler           func(shardID uint32, fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForSenderHandler          func(sender, fields string) (*data.TransactionsPoolForSender, error)
	GetLastPoolNonceForSenderHandler             func(sender string) (uint64, error)
	GetTransactionsPoolNonceGapsForSenderHandler func(sender string) (*data.TransactionsPoolNonceGaps, error)
	SendTransactionHandler                       func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler              func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionHandler                   func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                          func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler                        func(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error)
	GetHeartbeatDataHandler                      func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler                   func() (map[string]*data.ValidatorApiResponse, error)
	AuctionListHandler                           func() ([]*data.AuctionListValidatorAPIResponse, error)
	TransactionCostRequestHandler                func(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatusHandler                  func(txHash string, sender string) (string, error)
	GetProcessedTransactionStatusHandler         func(txHash string) (*data.ProcessStatusResponse, error)
	GetConfigMetricsHandler                      func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler                     func(shardID uint32) (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsHandler                     func(tokenType string) (*data.GenericAPIResponse, error)
	GetEnableEpochsMetricsHandler                func() (*data.GenericAPIResponse, error)
	GetEconomicsDataMetricsHandler               func() (*data.GenericAPIResponse, error)
	GetDirectStakedInfoCalled                    func() (*data.GenericAPIResponse, error)
	GetDelegatedInfoCalled                       func() (*data.GenericAPIResponse, error)
	GetRatingsConfigCalled                       func() (*data.GenericAPIResponse, error)
	GetTransactionByHashAndSenderAddressHandler  func(txHash string, sndAddr string, withResults bool) (*transaction.ApiTransactionResult, int, error)
	GetBlockByHashCalled                         func(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled                        func(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlocksByRoundCalled                       func(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error)
	GetInternalBlockByHashCalled                 func(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonceCalled                func(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalMiniBlockByHashCalled             func(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error)
	GetInternalStartOfEpochMetaBlockCalled       func(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalStartOfEpochValidatorsInfoCalled  func(epoch uint32) (*data.ValidatorsInfoApiResponse, error)
	GetHyperBlockByHashCalled                    func(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled                   func(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	ReloadObserversCalled                        func() data.NodesReloadResponse
	ReloadFullHistoryObserversCalled             func() data.NodesReloadResponse
	GetProofCalled                               func(string, string) (*data.GenericAPIResponse, error)
	GetProofDataTrieCalled                       func(string, string, string) (*data.GenericAPIResponse, error)
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
	GetNodesVersionsCalled                       func() (*data.GenericAPIResponse, error)
	GetAlteredAccountsByNonceCalled              func(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error)
	GetAlteredAccountsByHashCalled               func(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error)
	GetTriesStatisticsCalled                     func(shardID uint32) (*data.TrieStatisticsAPIResponse, error)
	GetEpochStartDataCalled                      func(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error)
	GetCodeHashCalled                            func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetGuardianDataCalled                        func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	IsDataTrieMigratedCalled                     func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetWaitingEpochsLeftForPublicKeyCalled       func(publicKey string) (*data.WaitingEpochsLeftApiResponse, error)
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
func (f *FacadeStub) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTsWithRoleCalled != nil {
		return f.GetESDTsWithRoleCalled(address, role, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetESDTsRoles -
func (f *FacadeStub) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTsRolesCalled != nil {
		return f.GetESDTsRolesCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetNFTTokenIDsRegisteredByAddress -
func (f *FacadeStub) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetNFTTokenIDsRegisteredByAddressCalled != nil {
		return f.GetNFTTokenIDsRegisteredByAddressCalled(address, options)
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

// GetRatingsConfig -
func (f *FacadeStub) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	return f.GetRatingsConfigCalled()
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
	if f.ValidatorStatisticsHandler != nil {
		return f.ValidatorStatisticsHandler()
	}

	return nil, nil
}

// AuctionList -
func (f *FacadeStub) AuctionList() ([]*data.AuctionListValidatorAPIResponse, error) {
	if f.AuctionListHandler != nil {
		return f.AuctionListHandler()
	}

	return nil, nil
}

// GetAccount -
func (f *FacadeStub) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	return f.GetAccountHandler(address, options)
}

// GetAccounts -
func (f *FacadeStub) GetAccounts(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error) {
	return f.GetAccountsHandler(addresses, options)
}

// GetKeyValuePairs -
func (f *FacadeStub) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return f.GetKeyValuePairsHandler(address, options)
}

// GetValueForKey -
func (f *FacadeStub) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	return f.GetValueForKeyHandler(address, key, options)
}

// GetGuardianData -
func (f *FacadeStub) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return f.GetGuardianDataCalled(address, options)
}

// GetShardIDForAddress -
func (f *FacadeStub) GetShardIDForAddress(address string) (uint32, error) {
	return f.GetShardIDForAddressHandler(address)
}

// GetESDTTokenData -
func (f *FacadeStub) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTTokenDataCalled != nil {
		return f.GetESDTTokenDataCalled(address, key, options)
	}

	return nil, nil
}

// GetAllESDTTokens -
func (f *FacadeStub) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetAllESDTTokensCalled != nil {
		return f.GetAllESDTTokensCalled(address, options)
	}

	return nil, nil
}

// GetESDTNftTokenData -
func (f *FacadeStub) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.GetESDTNftTokenDataCalled != nil {
		return f.GetESDTNftTokenDataCalled(address, key, nonce, options)
	}

	return nil, nil
}

// IsOldStorageForToken -
func (f *FacadeStub) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	if f.IsOldStorageForTokenCalled != nil {
		return f.IsOldStorageForTokenCalled(tokenID, nonce)
	}

	return false, nil
}

// GetTransactions -
func (f *FacadeStub) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return f.GetTransactionsHandler(address)
}

// GetTransactionByHashAndSenderAddress -
func (f *FacadeStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error) {
	return f.GetTransactionByHashAndSenderAddressHandler(txHash, sndAddr, withEvents)
}

// GetTransaction -
func (f *FacadeStub) GetTransaction(txHash string, withResults bool, relayedTxHash string) (*transaction.ApiTransactionResult, error) {
	return f.GetTransactionHandler(txHash, withResults, relayedTxHash)
}

// GetTransactionsPool -
func (f *FacadeStub) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	if f.GetTransactionsPoolHandler != nil {
		return f.GetTransactionsPoolHandler(fields)
	}

	return nil, nil
}

// GetTransactionsPoolForShard -
func (f *FacadeStub) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	if f.GetTransactionsPoolForShardHandler != nil {
		return f.GetTransactionsPoolForShardHandler(shardID, fields)
	}

	return nil, nil
}

// GetTransactionsPoolForSender -
func (f *FacadeStub) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	if f.GetTransactionsPoolForSenderHandler != nil {
		return f.GetTransactionsPoolForSenderHandler(sender, fields)
	}

	return nil, nil
}

// GetLastPoolNonceForSender -
func (f *FacadeStub) GetLastPoolNonceForSender(sender string) (uint64, error) {
	if f.GetLastPoolNonceForSenderHandler != nil {
		return f.GetLastPoolNonceForSenderHandler(sender)
	}

	return 0, nil
}

// GetTransactionsPoolNonceGapsForSender -
func (f *FacadeStub) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	if f.GetTransactionsPoolNonceGapsForSenderHandler != nil {
		return f.GetTransactionsPoolNonceGapsForSenderHandler(sender)
	}

	return nil, nil
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

// GetProcessedTransactionStatus -
func (f *FacadeStub) GetProcessedTransactionStatus(txHash string) (*data.ProcessStatusResponse, error) {
	return f.GetProcessedTransactionStatusHandler(txHash)
}

// SendUserFunds -
func (f *FacadeStub) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery -
func (f *FacadeStub) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData -
func (f *FacadeStub) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// GetBlockByHash -
func (f *FacadeStub) GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return f.GetBlockByHashCalled(shardID, hash, options)
}

// GetBlockByNonce -
func (f *FacadeStub) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return f.GetBlockByNonceCalled(shardID, nonce, options)
}

// GetBlocksByRound -
func (f *FacadeStub) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	if f.GetBlocksByRoundCalled != nil {
		return f.GetBlocksByRoundCalled(round, options)
	}
	return nil, nil
}

// GetInternalBlockByHash -
func (f *FacadeStub) GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return f.GetInternalBlockByHashCalled(shardID, hash, format)
}

// GetInternalBlockByNonce -
func (f *FacadeStub) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return f.GetInternalBlockByNonceCalled(shardID, nonce, format)
}

// GetInternalMiniBlockByHash -
func (f *FacadeStub) GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
	return f.GetInternalMiniBlockByHashCalled(shardID, hash, epoch, format)
}

// GetInternalStartOfEpochMetaBlock -
func (f *FacadeStub) GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return f.GetInternalStartOfEpochMetaBlockCalled(epoch, format)
}

// GetHyperBlockByHash -
func (f *FacadeStub) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByHashCalled(hash, options)
}

// GetHyperBlockByNonce -
func (f *FacadeStub) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return f.GetHyperBlockByNonceCalled(nonce, options)
}

// GetMetrics -
func (f *FacadeStub) GetMetrics() map[string]*data.EndpointMetrics {
	return f.GetMetricsCalled()
}

// GetMetricsForPrometheus -
func (f *FacadeStub) GetMetricsForPrometheus() string {
	return f.GetPrometheusMetricsCalled()
}

// GetGenesisNodesPubKeys -
func (f *FacadeStub) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	return f.GetGenesisNodesPubKeysCalled()
}

// GetGasConfigs -
func (f *FacadeStub) GetGasConfigs() (*data.GenericAPIResponse, error) {
	return f.GetGasConfigsCalled()
}

// GetAboutInfo -
func (f *FacadeStub) GetAboutInfo() (*data.GenericAPIResponse, error) {
	return f.GetAboutInfoCalled()
}

// GetNodesVersions -
func (f *FacadeStub) GetNodesVersions() (*data.GenericAPIResponse, error) {
	return f.GetNodesVersionsCalled()
}

// GetAlteredAccountsByNonce -
func (f *FacadeStub) GetAlteredAccountsByNonce(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	if f.GetAlteredAccountsByNonceCalled != nil {
		return f.GetAlteredAccountsByNonceCalled(shardID, nonce, options)
	}
	return nil, nil
}

// GetAlteredAccountsByHash -
func (f *FacadeStub) GetAlteredAccountsByHash(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	if f.GetAlteredAccountsByHashCalled != nil {
		return f.GetAlteredAccountsByHashCalled(shardID, hash, options)
	}

	return nil, nil
}

// GetTriesStatistics -
func (f *FacadeStub) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	if f.GetTriesStatisticsCalled != nil {
		return f.GetTriesStatisticsCalled(shardID)
	}
	return &data.TrieStatisticsAPIResponse{}, nil
}

// GetEpochStartData -
func (f *FacadeStub) GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error) {
	return f.GetEpochStartDataCalled(epoch, shardID)
}

// GetInternalStartOfEpochValidatorsInfo -
func (f *FacadeStub) GetInternalStartOfEpochValidatorsInfo(epoch uint32) (*data.ValidatorsInfoApiResponse, error) {
	return f.GetInternalStartOfEpochValidatorsInfoCalled(epoch)
}

// GetCodeHash -
func (f *FacadeStub) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return f.GetCodeHashCalled(address, options)
}

// IsDataTrieMigrated -
func (f *FacadeStub) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if f.IsDataTrieMigratedCalled != nil {
		return f.IsDataTrieMigratedCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetWaitingEpochsLeftForPublicKey -
func (f *FacadeStub) GetWaitingEpochsLeftForPublicKey(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
	if f.GetWaitingEpochsLeftForPublicKeyCalled != nil {
		return f.GetWaitingEpochsLeftForPublicKeyCalled(publicKey)
	}
	return &data.WaitingEpochsLeftApiResponse{}, nil
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
