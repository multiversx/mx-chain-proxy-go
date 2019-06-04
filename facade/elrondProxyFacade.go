package facade

import "github.com/ElrondNetwork/elrond-proxy-go/data"

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
		return nil, ErrNilAccountProccessor
	}
	//TODO check txProc when implemented

	return &ElrondProxyFacade{
		accountProc: accountProc,
		txProc:      txProc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ElrondProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}
