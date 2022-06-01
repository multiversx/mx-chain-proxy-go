package services

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRosettaTxFromUnsignedTxSendFunds(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	tp := newTransactionParser(&mocks.ElrondProviderMock{}, cfg, networkCfg)

	tx := &data.FullTransaction{
		Hash:     "hash-hash",
		Receiver: "receiverAddress",
		Sender:   "senderAddress",
		Value:    "1234",
	}

	rosettaTx, ok := tp.createRosettaTxFromUnsignedTx(tx)
	assert.True(t, ok)
	assert.Equal(t, &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: tx.Hash,
		},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opScResult,
				Status: &OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: tx.Sender,
				},
				Amount: &types.Amount{
					Value:    "-" + tx.Value,
					Currency: tp.config.Currency,
				},
			},
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 1,
				},
				Type:   opScResult,
				Status: &OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: tx.Receiver,
				},
				Amount: &types.Amount{
					Value:    tx.Value,
					Currency: tp.config.Currency,
				},
			},
		},
	}, rosettaTx)
}

func TestCreateRosettaTxFromUnsignedTxWithBadValueShouldBeIgnored(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	tp := newTransactionParser(&mocks.ElrondProviderMock{}, cfg, networkCfg)

	rosettaTx, ok := tp.createRosettaTxFromUnsignedTx(&data.FullTransaction{
		Hash:     "hash-hash",
		Receiver: "receiverAddress",
		GasLimit: 1000,
		Value:    "0",
	})
	require.False(t, ok)
	require.Nil(t, rosettaTx)

	rosettaTx, ok = tp.createRosettaTxFromUnsignedTx(&data.FullTransaction{
		Hash:     "hash-hash",
		Receiver: "receiverAddress",
		GasLimit: 1000,
		Value:    "-1",
	})
	require.False(t, ok)
	require.Nil(t, rosettaTx)
}

func TestCreateRosettaTxFromUnsignedTxRefundGas(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	tp := newTransactionParser(&mocks.ElrondProviderMock{}, cfg, networkCfg)

	tx := &data.FullTransaction{
		Hash:     "hash-hash",
		Sender:   "senderAddress",
		Receiver: "receiverAddress",
		GasLimit: 1000,
		Value:    "1234",
		Nonce:    1,
		IsRefund: true,
	}

	rosettaTx, ok := tp.createRosettaTxFromUnsignedTx(tx)
	assert.True(t, ok)
	assert.Equal(t, &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: tx.Hash,
		},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opScResult,
				Status: &OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: tx.Receiver,
				},
				Amount: &types.Amount{
					Value:    tx.Value,
					Currency: tp.config.Currency,
				},
			},
		},
	}, rosettaTx)
}

func TestCreateOperationsFromPreparedTx(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	tp := newTransactionParser(&mocks.ElrondProviderMock{}, cfg, networkCfg)

	preparedTx := &data.Transaction{
		Value:    "12345",
		Receiver: "receiver",
		Sender:   "sender",
	}

	expectedOperations := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: preparedTx.Sender,
			},
			Amount: &types.Amount{
				Value:    "-" + preparedTx.Value,
				Currency: tp.config.Currency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{Index: 0},
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: preparedTx.Receiver,
			},
			Amount: &types.Amount{
				Value:    preparedTx.Value,
				Currency: tp.config.Currency,
			},
		},
	}

	operations := tp.createOperationsFromPreparedTx(preparedTx)
	assert.Equal(t, expectedOperations, operations)
}
