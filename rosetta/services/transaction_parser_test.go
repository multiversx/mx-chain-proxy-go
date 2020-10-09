package services

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRosettaTxFromUnsignedTx(t *testing.T) {
	t.Parallel()

	networkCfg := &client.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	tp := newTransactionParser(cfg, networkCfg)

	tx := &data.FullTransaction{
		Hash:     "hash-hash",
		Receiver: "receiverAddress",
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
				Status: OpStatusSuccess,
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

	networkCfg := &client.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	tp := newTransactionParser(cfg, networkCfg)

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
