package services

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestBlockAPIService_BlockByIndex(t *testing.T) {
	t.Parallel()

	blockIndex := int64(10)
	blockHash := "hash-hash-hash"
	prevBlockHash := "prev-hash-hash-hash"
	txHash := "txHash"

	fullTx := &data.FullTransaction{
		Hash:          txHash,
		MiniBlockType: block.TxBlock.String(),
		Receiver:      "erd1uml89f3lqqfxan67dnnlytd0r3mz3v684zxdhqq60gs5u7qa9yjqa5dgqp",
		Sender:        "erd18f33a94auxr4v8v23wu8gwv7mzf408jsskktvj4lcmcrv4v5jmqs5x3kdn",
		Value:         "1234",
		GasLimit:      100,
		GasPrice:      10,
	}

	elrondClientMock := &mocks.ElrondClientMock{
		GetBlockByNonceCalled: func(nonce int64) (*data.Hyperblock, error) {
			return &data.Hyperblock{
				Nonce:         uint64(blockIndex),
				Hash:          blockHash,
				Round:         10,
				PrevBlockHash: prevBlockHash,
				Transactions: []*data.FullTransaction{
					fullTx,
				},
			}, nil
		},
	}
	cfg := &configuration.Configuration{}
	blockAPIService := NewBlockAPIService(elrondClientMock, cfg)

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
					{
						OperationIdentifier: &types.OperationIdentifier{
							Index: 2,
						},
						Type:   opFee,
						Status: OpStatusSuccess,
						Account: &types.AccountIdentifier{
							Address: fullTx.Sender,
						},
						Amount: &types.Amount{
							Value:    "-" + computeTxFee(fullTx),
							Currency: nil,
						},
					},
				},
			},
		},
		Metadata: objectsMap{
			"epoch": uint32(0),
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
