package txcost

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionCostProcessor_IndexOutOfBounds(t *testing.T) {
	t.Parallel()

	coreProc := &mock.ProcessorStub{}
	newTxCostProcessor, _ := NewTransactionCostProcessor(
		coreProc, &mock.PubKeyConverterMock{})
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
		coreProc, &mock.PubKeyConverterMock{})
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
	newTxCostProcessor.scrsToExecute = append(newTxCostProcessor.scrsToExecute, &smartContractResult.SmartContractResult{
		GasLimit: 200,
	})

	res := &data.TxCostResponseData{
		TxCost: 500,
	}

	expectedGas := uint64(1300)
	newTxCostProcessor.prepareGasUsed(0, 0, res)
	require.Equal(t, expectedGas, res.TxCost)
	require.Equal(t, "", res.RetMessage)
}
