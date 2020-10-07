package services

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func TestMempoolAPIService_MempoolTransactionCannotFindTxInPool(t *testing.T) {
	t.Parallel()

	elrondClientMock := &mocks.ElrondClientMock{
		GetTransactionByHashFromPoolCalled: func(txHash string) (*data.FullTransaction, bool) {
			return nil, false
		},
	}

	networkCfg := &client.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	mempoolApiService := NewMempoolApiService(elrondClientMock, cfg, networkCfg)

	txHash := "hash-hash-hash"
	txResponse, err := mempoolApiService.MempoolTransaction(context.Background(), &types.MempoolTransactionRequest{
		NetworkIdentifier:     nil,
		TransactionIdentifier: &types.TransactionIdentifier{Hash: txHash},
	})
	require.Equal(t, ErrTransactionIsNotInPool, err)
	require.Nil(t, txResponse)
}

func TestMempoolAPIService_MempoolTransaction(t *testing.T) {
	t.Parallel()

	txHash := "hash-hash-hash"
	fullTx := &data.FullTransaction{
		Hash:     txHash,
		Type:     string(transaction.TxTypeNormal),
		Receiver: "erd1uml89f3lqqfxan67dnnlytd0r3mz3v684zxdhqq60gs5u7qa9yjqa5dgqp",
		Sender:   "erd18f33a94auxr4v8v23wu8gwv7mzf408jsskktvj4lcmcrv4v5jmqs5x3kdn",
		Value:    "1234",
		GasLimit: 100,
		GasPrice: 10,
	}
	elrondClientMock := &mocks.ElrondClientMock{
		GetTransactionByHashFromPoolCalled: func(txHash string) (*data.FullTransaction, bool) {
			return fullTx, true
		},
	}
	networkCfg := &client.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	mempoolApiService := NewMempoolApiService(elrondClientMock, cfg, networkCfg)

	expectedRosettaTx := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{Hash: txHash},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opTransfer,
				Status: OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: fullTx.Sender,
				},
				Amount: &types.Amount{
					Value:    "-" + fullTx.Value,
					Currency: nil,
				},
			},
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 1,
				},
				RelatedOperations: []*types.OperationIdentifier{
					{Index: 0},
				},
				Type:   opTransfer,
				Status: OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: fullTx.Receiver,
				},
				Amount: &types.Amount{
					Value:    fullTx.Value,
					Currency: nil,
				},
			},
		},
	}

	txResponse, err := mempoolApiService.MempoolTransaction(context.Background(), &types.MempoolTransactionRequest{
		NetworkIdentifier:     nil,
		TransactionIdentifier: &types.TransactionIdentifier{Hash: txHash},
	})
	require.Nil(t, err)
	require.Equal(t, expectedRosettaTx, txResponse.Transaction)
}
