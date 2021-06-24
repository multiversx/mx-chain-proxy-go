package txcost

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionCostProcessor_IndexOutOfBounds(t *testing.T) {
	t.Parallel()

	coreProc := &mock.ProcessorStub{}
	newTxCostProcessor, _ := NewTransactionCostProcessor(
		coreProc, &mock.PubKeyConverterMock{}, "1500000000", "15000000000")
	newTxCostProcessor.responses = append(newTxCostProcessor.responses, &data.ResponseTxCost{})
	newTxCostProcessor.responses = append(newTxCostProcessor.responses, &data.ResponseTxCost{})
	newTxCostProcessor.responses = append(newTxCostProcessor.responses, &data.ResponseTxCost{})

	res := &data.TxCostResponseData{}
	newTxCostProcessor.prepareGasUsed(0, 0, res)
	require.Equal(t, "something went wrong", res.RetMessage)
}

func TestTransactionCostProcessor_PrepareGasUsedShouldWork(t *testing.T) {
	t.Parallel()

	coreProc := &mock.ProcessorStub{}
	newTxCostProcessor, _ := NewTransactionCostProcessor(
		coreProc, &mock.PubKeyConverterMock{}, "1500000000", "15000000000")
	newTxCostProcessor.responses = append(newTxCostProcessor.responses, &data.ResponseTxCost{
		Data: data.TxCostResponseData{
			TxCost: 500,
		},
	})
	newTxCostProcessor.responses = append(newTxCostProcessor.responses, &data.ResponseTxCost{
		Data: data.TxCostResponseData{
			TxCost: 1000,
		},
	})
	newTxCostProcessor.txsFromSCR = append(newTxCostProcessor.txsFromSCR, &data.Transaction{
		GasLimit: 200,
	})

	res := &data.TxCostResponseData{
		TxCost: 500,
	}
	newTxCostProcessor.prepareGasUsed(0, 0, res)
	require.Equal(t, uint64(1300), res.TxCost)
	require.Equal(t, "", res.RetMessage)
}
