package factory

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestTxNotarizationChecker(t *testing.T) {
	t.Parallel()

	tnc := NewTxNotarizationChecker()
	require.False(t, tnc.IsInterfaceNil())
}

func TestTxNotarizationChecker_IsNotarized(t *testing.T) {
	t.Parallel()

	t.Run("tx is notarized", func(t *testing.T) {
		tnc := NewTxNotarizationChecker()
		tx := transaction.ApiTransactionResult{
			NotarizedAtSourceInMetaNonce:      1,
			NotarizedAtDestinationInMetaNonce: 1,
		}
		isNotarized := tnc.IsNotarized(tx)
		require.True(t, isNotarized)
	})
	t.Run("tx is not notarized", func(t *testing.T) {
		tnc := NewTxNotarizationChecker()
		isNotarized := tnc.IsNotarized(transaction.ApiTransactionResult{})
		require.False(t, isNotarized)
	})
}
