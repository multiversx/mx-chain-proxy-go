package factory

import (
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

type txNotarizationChecker struct{}

// NewTxNotarizationChecker creates a new tx notarization checker
func NewTxNotarizationChecker() *txNotarizationChecker {
	return &txNotarizationChecker{}
}

// IsNotarized returns if tx is notarized
func (tnc *txNotarizationChecker) IsNotarized(tx transaction.ApiTransactionResult) bool {
	return tx.NotarizedAtSourceInMetaNonce > 0 && tx.NotarizedAtDestinationInMetaNonce > 0
}

// IsInterfaceNil returns true if there is no value under the interface
func (tnc *txNotarizationChecker) IsInterfaceNil() bool {
	return tnc == nil
}
