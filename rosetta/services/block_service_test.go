package services

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestBlockAPIService_BlockByIndex(t *testing.T) {
	t.Parallel()

	blockIndex := int64(10)
	round := uint64(12)
	blockHash := "hash-hash-hash"
	prevBlockHash := "prev-hash-hash-hash"

	txHash := "txHash"
	fullTxNormal := &data.FullTransaction{
		Hash:     txHash,
		Type:     string(transaction.TxTypeNormal),
		Receiver: "erd1uml89f3lqqfxan67dnnlytd0r3mz3v684zxdhqq60gs5u7qa9yjqa5dgqp",
		Sender:   "erd18f33a94auxr4v8v23wu8gwv7mzf408jsskktvj4lcmcrv4v5jmqs5x3kdn",
		Value:    "1234",
		GasLimit: 100,
		GasPrice: 10,
	}

	rewardHash := "rewardHash"
	rewardTx := &data.FullTransaction{
		Hash:     rewardHash,
		Receiver: "erd1uml89f3lqqfxan67dnnlytd0r3mz3v684zxdhqq60gs5u7qa9yjqa5dgqp",
		Value:    "1111",
		Type:     string(transaction.TxTypeReward),
	}

	invalidTxHash := "invalidTx"
	invalidTx := &data.FullTransaction{
		Hash:     invalidTxHash,
		Sender:   "erd1uml89f3lqqfxan67dnnlytd0r3mz3v684zxdhqq60gs5u7qa9yjqa5dgqp",
		GasLimit: 100,
		GasPrice: 10,
		Type:     string(transaction.TxTypeInvalid),
	}

	elrondProviderMock := &mocks.ElrondProviderMock{
		GetBlockByNonceCalled: func(nonce int64) (*data.Hyperblock, error) {
			return &data.Hyperblock{
				Nonce:         uint64(blockIndex),
				Hash:          blockHash,
				Round:         round,
				PrevBlockHash: prevBlockHash,
				Transactions: []*data.FullTransaction{
					fullTxNormal, rewardTx, invalidTx,
				},
			}, nil
		},
	}

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
	}
	cfg := &configuration.Configuration{}
	blockAPIService := NewBlockAPIService(elrondProviderMock, cfg, networkCfg)
	tp := newTransactionParser(elrondProviderMock, cfg, networkCfg)

	expectedBlock := &types.Block{
		BlockIdentifier: &types.BlockIdentifier{
			Index: blockIndex,
			Hash:  blockHash,
		},
		ParentBlockIdentifier: &types.BlockIdentifier{
			Index: blockIndex - 1,
			Hash:  prevBlockHash,
		},
		Timestamp: 0,
		Transactions: []*types.Transaction{
			{
				TransactionIdentifier: &types.TransactionIdentifier{Hash: txHash},
				Operations: []*types.Operation{
					{
						OperationIdentifier: &types.OperationIdentifier{
							Index: 0,
						},
						Type:   opTransfer,
						Status: &OpStatusSuccess,
						Account: &types.AccountIdentifier{
							Address: fullTxNormal.Sender,
						},
						Amount: &types.Amount{
							Value:    "-" + fullTxNormal.Value,
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
						Status: &OpStatusSuccess,
						Account: &types.AccountIdentifier{
							Address: fullTxNormal.Receiver,
						},
						Amount: &types.Amount{
							Value:    fullTxNormal.Value,
							Currency: nil,
						},
					},
					{
						OperationIdentifier: &types.OperationIdentifier{
							Index: 2,
						},
						Type:   opFee,
						Status: &OpStatusSuccess,
						Account: &types.AccountIdentifier{
							Address: fullTxNormal.Sender,
						},
						Amount: &types.Amount{
							Value:    "-" + tp.computeTxFee(fullTxNormal).String(),
							Currency: nil,
						},
					},
				},
			},
			{

				TransactionIdentifier: &types.TransactionIdentifier{
					Hash: rewardTx.Hash,
				},
				Operations: []*types.Operation{
					{
						OperationIdentifier: &types.OperationIdentifier{
							Index: 0,
						},
						Type:   opReward,
						Status: &OpStatusSuccess,
						Account: &types.AccountIdentifier{
							Address: rewardTx.Receiver,
						},
						Amount: &types.Amount{
							Value:    rewardTx.Value,
							Currency: tp.config.Currency,
						},
					},
				},
			},
			{
				TransactionIdentifier: &types.TransactionIdentifier{
					Hash: invalidTx.Hash,
				},
				Operations: []*types.Operation{
					{
						OperationIdentifier: &types.OperationIdentifier{
							Index: 0,
						},
						Type:   opInvalid,
						Status: &OpStatusSuccess,
						Account: &types.AccountIdentifier{
							Address: invalidTx.Sender,
						},
						Amount: &types.Amount{
							Value:    "-" + tp.computeTxFee(invalidTx).String(),
							Currency: tp.config.Currency,
						},
					},
				},
			},
		},
		Metadata: objectsMap{
			"epoch": uint32(0),
			"round": round,
		},
	}
	blockResponse, err := blockAPIService.Block(context.Background(), &types.BlockRequest{
		NetworkIdentifier: nil,
		BlockIdentifier: &types.PartialBlockIdentifier{
			Index: &blockIndex,
			Hash:  nil,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, expectedBlock, blockResponse.Block)
}
