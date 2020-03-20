package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ElrondProxyFacade implements the facade used in api calls
type ElrondProxyFacade struct {
	accountProc    AccountProcessor
	txProc         TransactionProcessor
	scQueryService SCQueryService
	heartbeatProc  HeartbeatProcessor
	faucetProc     FaucetProcessor
	nodeStatusProc NodeStatusProcessor
	web3Proc       ProcessorWeb3
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	scQueryService SCQueryService,
	heartbeatProc HeartbeatProcessor,
	faucetProc FaucetProcessor,
	nodeStatusProc NodeStatusProcessor,
	web3Proc ProcessorWeb3,
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
	if faucetProc == nil {
		return nil, ErrNilFaucetProcessor
	}
	if nodeStatusProc == nil {
		return nil, ErrNilNodeStatusProcessor
	}

	return &ElrondProxyFacade{
		accountProc:    accountProc,
		txProc:         txProc,
		scQueryService: scQueryService,
		heartbeatProc:  heartbeatProc,
		faucetProc:     faucetProc,
		nodeStatusProc: nodeStatusProc,
		web3Proc:       web3Proc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ElrondProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}

// SendTransaction should sends the transaction to the correct observer
func (epf *ElrondProxyFacade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return epf.txProc.SendTransaction(tx)
}

// SendMultipleTransactions should send the transactions to the correct observers
func (epf *ElrondProxyFacade) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	return epf.txProc.SendMultipleTransactions(txs)
}

// SendTransactionCostRequest should return how many gas units a transaction will cost
func (epf *ElrondProxyFacade) SendTransactionCostRequest(tx *data.Transaction) (string, error) {
	return epf.txProc.SendTransactionCostRequest(tx)
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
func (epf *ElrondProxyFacade) ExecuteSCQuery(query *process.SCQuery) (*vmcommon.VMOutput, error) {
	return epf.scQueryService.ExecuteQuery(query)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (epf *ElrondProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return epf.heartbeatProc.GetHeartbeatData()
}

// GetHeartbeatData retrieves the node status from one observer
func (epf *ElrondProxyFacade) GetNodeStatusData(shardId string) (map[string]interface{}, error) {
	return epf.nodeStatusProc.GetNodeStatusData(shardId)
}

// ValidatorStatistics will return the statistics from an observer
func (epf *ElrondProxyFacade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return epf.accountProc.ValidatorStatistics()
}

func (epf *ElrondProxyFacade) PrepareDataForRequest(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	return epf.web3Proc.PrepareDataForRequest(r)
}
