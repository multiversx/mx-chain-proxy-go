package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ElrondProxyFacade implements the facade used in api calls
type ElrondProxyFacade struct {
	accountProc   AccountProcessor
	txProc        TransactionProcessor
	getValuesProc GetValuesProcessor
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	getValuesProc GetValuesProcessor,
) (*ElrondProxyFacade, error) {

	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}
	if getValuesProc == nil {
		return nil, ErrNilGetValueProcessor
	}

	return &ElrondProxyFacade{
		accountProc:   accountProc,
		txProc:        txProc,
		getValuesProc: getValuesProc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ElrondProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}

// SendTransaction should sends the transaction to the correct observer
func (epf *ElrondProxyFacade) SendTransaction(
	nonce uint64,
	sender string,
	receiver string,
	value *big.Int,
	data string,
	signature []byte,
	gasPrice uint64,
	gasLimit uint64,
) (string, error) {

	return epf.txProc.SendTransaction(nonce, sender, receiver, value, data,
		signature, gasPrice, gasLimit)
}

// SendUserFunds should send a transaction to load one user's account with extra funds from the observer
func (epf *ElrondProxyFacade) SendUserFunds(receiver string) error {
	return epf.txProc.SendUserFunds(receiver)
}

// GetDataValue retrieves data from existing SC trie
func (epf *ElrondProxyFacade) GetDataValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return epf.getValuesProc.GetDataValue(address, funcName, argsBuff...)
}
