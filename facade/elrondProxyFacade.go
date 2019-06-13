package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ElrondProxyFacade implements the facade used in api calls
type ElrondProxyFacade struct {
	accountProc AccountProcessor
	txProc      TransactionProcessor
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
) (*ElrondProxyFacade, error) {

	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}

	return &ElrondProxyFacade{
		accountProc: accountProc,
		txProc:      txProc,
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
	code string,
	signature []byte,
) (string, error) {

	return epf.txProc.SendTransaction(nonce, sender, receiver, value, code, signature)
}

// SendUserFunds should send a transaction to load one user's account with extra funds from the observer
func (epf *ElrondProxyFacade) SendUserFunds(receiver string) error {
	return epf.txProc.SendUserFunds(receiver)
}
