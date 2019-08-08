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
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	vmValuesProc VmValuesProcessor,
	heartbeatProc HeartbeatProcessor,
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

	return &ElrondProxyFacade{
		accountProc:   accountProc,
		txProc:        txProc,
		vmValuesProc:  vmValuesProc,
		heartbeatProc: heartbeatProc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ElrondProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}

// SendTransaction should sends the transaction to the correct observer
func (epf *ElrondProxyFacade) SendTransaction(tx *data.Transaction) (string, error) {
	return epf.txProc.SendTransaction(tx)
}

// SendUserFunds should send a transaction to load one user's account with extra funds from the observer
func (epf *ElrondProxyFacade) SendUserFunds(receiver string, value *big.Int) error {
	return epf.txProc.SendUserFunds(receiver, value)
}

// GetVmValue retrieves data from existing SC trie through the use of a VM
func (epf *ElrondProxyFacade) GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return epf.vmValuesProc.GetVmValue(resType, address, funcName, argsBuff...)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (epf *ElrondProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return epf.heartbeatProc.GetHeartbeatData()
}
