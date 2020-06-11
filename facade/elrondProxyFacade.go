package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ElrondProxyFacade implements the facade used in api calls
type ElrondProxyFacade struct {
	accountProc    AccountProcessor
	txProc         TransactionProcessor
	scQueryService SCQueryService
	heartbeatProc  HeartbeatProcessor
	valStatsProc   ValidatorStatisticsProcessor
	faucetProc     FaucetProcessor
	nodeStatusProc NodeStatusProcessor
	blockProc      BlockProcessor
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	scQueryService SCQueryService,
	heartbeatProc HeartbeatProcessor,
	valStatsProc ValidatorStatisticsProcessor,
	faucetProc FaucetProcessor,
	nodeStatusProc NodeStatusProcessor,
	blockProc BlockProcessor,
) (*ElrondProxyFacade, error) {

	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}
	if scQueryService == nil {
		return nil, ErrNilSCQueryService
	}
	if heartbeatProc == nil {
		return nil, ErrNilHeartbeatProcessor
	}
	if valStatsProc == nil {
		return nil, ErrNilValidatorStatisticsProcessor
	}
	if faucetProc == nil {
		return nil, ErrNilFaucetProcessor
	}
	if nodeStatusProc == nil {
		return nil, ErrNilNodeStatusProcessor
	}
	if blockProc == nil {
		return nil, ErrNilBlockProcessor
	}

	return &ElrondProxyFacade{
		accountProc:    accountProc,
		txProc:         txProc,
		scQueryService: scQueryService,
		heartbeatProc:  heartbeatProc,
		valStatsProc:   valStatsProc,
		faucetProc:     faucetProc,
		nodeStatusProc: nodeStatusProc,
		blockProc:      blockProc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ElrondProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}

// GetValueForKey returns the value for the given address and key
func (epf *ElrondProxyFacade) GetValueForKey(address string, key string) (string, error) {
	return epf.accountProc.GetValueForKey(address, key)
}

// GetTransactions returns transactions by address
func (epf *ElrondProxyFacade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return epf.accountProc.GetTransactions(address)
}

// SendTransaction should sends the transaction to the correct observer
func (epf *ElrondProxyFacade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return epf.txProc.SendTransaction(tx)
}

// SendMultipleTransactions should send the transactions to the correct observers
func (epf *ElrondProxyFacade) SendMultipleTransactions(txs []*data.Transaction) (data.ResponseMultipleTransactions, error) {
	return epf.txProc.SendMultipleTransactions(txs)
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (epf *ElrondProxyFacade) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return epf.txProc.TransactionCostRequest(tx)
}

// GetTransactionStatus should return transaction status
func (epf *ElrondProxyFacade) GetTransactionStatus(txHash string) (string, error) {
	return epf.txProc.GetTransactionStatus(txHash)
}

// GetTransaction should return a transaction by hash
func (epf *ElrondProxyFacade) GetTransaction(txHash string) (*transaction.ApiTransactionResult, error) {
	return epf.txProc.GetTransaction(txHash)
}

// GetTransactionByHashAndSenderAddress should return a transaction by hash and sender address
func (epf *ElrondProxyFacade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string) (*transaction.ApiTransactionResult, int, error) {
	return epf.txProc.GetTransactionByHashAndSenderAddress(txHash, sndAddr)
}

// SendUserFunds should send a transaction to load one user's account with extra funds from an account in the pem file
func (epf *ElrondProxyFacade) SendUserFunds(receiver string, value *big.Int) error {
	senderSk, senderPk, err := epf.faucetProc.SenderDetailsFromPem(receiver)
	if err != nil {
		return err
	}

	senderAccount, err := epf.accountProc.GetAccount(senderPk)
	if err != nil {
		return err
	}

	tx, err := epf.faucetProc.GenerateTxForSendUserFunds(senderSk, senderPk, senderAccount.Nonce, receiver, value)
	if err != nil {
		return err
	}

	_, _, err = epf.txProc.SendTransaction(tx)
	return err
}

// ExecuteSCQuery retrieves data from existing SC trie through the use of a VM
func (epf *ElrondProxyFacade) ExecuteSCQuery(query *data.SCQuery) (*vmcommon.VMOutput, error) {
	return epf.scQueryService.ExecuteQuery(query)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (epf *ElrondProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return epf.heartbeatProc.GetHeartbeatData()
}

// GetNetworkConfigMetrics retrieves the node's configuration's metrics
func (epf *ElrondProxyFacade) GetNetworkConfigMetrics() (map[string]interface{}, error) {
	return epf.nodeStatusProc.GetNetworkConfigMetrics()
}

// GetNetworkStatusMetrics retrieves the node's network metrics for a given shard
func (epf *ElrondProxyFacade) GetNetworkStatusMetrics(shardID uint32) (map[string]interface{}, error) {
	return epf.nodeStatusProc.GetNetworkStatusMetrics(shardID)
}

// ValidatorStatistics will return the statistics from an observer
func (epf *ElrondProxyFacade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	valStats, err := epf.valStatsProc.GetValidatorStatistics()
	if err != nil {
		return nil, err
	}

	return valStats.Statistics, nil
}

// GetBlockByShardIDAndNonce returns block by shardID and nonce
func (epf *ElrondProxyFacade) GetBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error) {
	return epf.blockProc.GetBlockByShardIDAndNonce(shardID, nonce)
}
