package v1_0

import (
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

// TransactionProcessorV1_0 represents the transaction processor that handles actions for version v1.0
type TransactionProcessorV1_0 struct {
	*process.TransactionProcessor
}

// GetTransaction overrides the base function in order to change it's logic for the version v1.0
func (tpv11 *TransactionProcessorV1_0) GetTransaction(txHash string) (*data.FullTransaction, error) {
	originalResponse, err := tpv11.TransactionProcessor.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}

	originalResponse.Status = changeStatus(string(originalResponse.Status))
	return originalResponse, nil
}

func changeStatus(input string) transaction.TxStatus {
	// TODO: change this after switch to new elrond-go version and use constants
	switch input {
	case "success":
		return "Success"
	case "fail":
		return "Not Executed"
	case "invalid":
		return "Invalid"
	case "pending":
		return "Pending"
	default:
		return transaction.TxStatus(input)
	}
}
