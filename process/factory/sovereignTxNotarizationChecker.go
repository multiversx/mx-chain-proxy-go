package factory

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

type sovereignTxNotarizationChecker struct{}

// NewSovereignTxNotarizationChecker creates a new sovereign tx notarization checker
func NewSovereignTxNotarizationChecker() *sovereignTxNotarizationChecker {
	return &sovereignTxNotarizationChecker{}
}

// IsNotarized returns true
func (stnc *sovereignTxNotarizationChecker) IsNotarized(_ transaction.ApiTransactionResult) bool {
	return true
}

// IsInterfaceNil returns true if there is no value under the interface
func (stnc *sovereignTxNotarizationChecker) IsInterfaceNil() bool {
	return stnc == nil
}
