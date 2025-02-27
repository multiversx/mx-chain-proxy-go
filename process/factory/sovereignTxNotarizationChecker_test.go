package factory

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestSovereignTxNotarizationChecker(t *testing.T) {
	t.Parallel()

	tnc := NewSovereignTxNotarizationChecker()
	require.False(t, tnc.IsInterfaceNil())
}

func TestSovereignTxNotarizationChecker_IsNotarized(t *testing.T) {
	t.Parallel()

	tnc := NewSovereignTxNotarizationChecker()
	isNotarized := tnc.IsNotarized(transaction.ApiTransactionResult{})
	require.True(t, isNotarized)
}
