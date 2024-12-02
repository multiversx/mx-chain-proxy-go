package mock

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

// TxNotarizationCheckerMock -
type TxNotarizationCheckerMock struct {
	IsNotarizedCalled func(tx transaction.ApiTransactionResult) bool
}

// IsNotarized -
func (tnc *TxNotarizationCheckerMock) IsNotarized(tx transaction.ApiTransactionResult) bool {
	if tnc.IsNotarizedCalled != nil {
		return tnc.IsNotarizedCalled(tx)
	}

	return false
}

// IsInterfaceNil -
func (tnc *TxNotarizationCheckerMock) IsInterfaceNil() bool {
	return tnc == nil
}
