package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// TransactionCostHandlerStub -
type TransactionCostHandlerStub struct {
	RezolveCostRequestCalled func(tx *data.Transaction) (*data.TxCostResponseData, error)
}

// RezolveCostRequest -
func (tchs *TransactionCostHandlerStub) RezolveCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	if tchs.RezolveCostRequestCalled != nil {
		return tchs.RezolveCostRequestCalled(tx)
	}

	return nil, nil
}
