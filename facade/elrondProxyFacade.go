package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ElrondProxyFacade implements the facade used in api calls
type ElrondProxyFacade struct {
	accountProc   AccountProcessor
	txProc        TransactionProcessor
	vmValuesProc  VmValuesProcessor
	heartbeatProc HeartbeatProcessor
	faucetProc    FaucetProcessor
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	vmValuesProc VmValuesProcessor,
	heartbeatProc HeartbeatProcessor,
	faucetProc FaucetProcessor,
) (*ElrondProxyFacade, error) {

	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}
	if vmValuesProc == nil {
		return nil, ErrNilVmValueProcessor
	}
	if heartbeatProc == nil {
		return nil, ErrNilHeartbeatProcessor
	}
	if faucetProc == nil {
		return nil, ErrNilFaucetProcessor
	}

	return &ElrondProxyFacade{
		accountProc:   accountProc,
		txProc:        txProc,
		vmValuesProc:  vmValuesProc,
		heartbeatProc: heartbeatProc,
		faucetProc:    faucetProc,
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

// GetVmValue retrieves data from existing SC trie through the use of a VM
func (epf *ElrondProxyFacade) GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return epf.vmValuesProc.GetVmValue(resType, address, funcName, argsBuff...)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (epf *ElrondProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return epf.heartbeatProc.GetHeartbeatData()
}
