package services

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func parseTxsFromHyperBlock(hyperBlock *data.Hyperblock) []*types.Transaction {
	txs := make([]*types.Transaction, 0)
	for _, eTx := range hyperBlock.Transactions {
		switch eTx.MiniBlockType {
		case block.TxBlock.String():
			txs = append(txs, createRosettaTxFromMoveBalance(eTx))
		case block.RewardsBlock.String():
			txs = append(txs, createRosettaTxFromReward(eTx))
		case block.SmartContractResultBlock.String():
			tx, ok := createRosettaTxFromUnsignedTx(eTx)
			if !ok {
				continue
			}

			txs = append(txs, tx)
		default:
			continue
		}
	}

	return txs
}

func createRosettaTxFromUnsignedTx(eTx *data.FullTransaction) (*types.Transaction, bool) {
	if eTx.Value == "0" {
		return nil, false
	}

	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: eTx.Hash,
		},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opScResult,
				Status: OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: eTx.Receiver,
				},
				Amount: &types.Amount{
					Value:    eTx.Value,
					Currency: ElrondCurrency,
				},
			},
		},
	}, true
}

func createRosettaTxFromReward(eTx *data.FullTransaction) *types.Transaction {
	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: eTx.Hash,
		},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opReward,
				Status: OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: eTx.Receiver,
				},
				Amount: &types.Amount{
					Value:    eTx.Value,
					Currency: ElrondCurrency,
				},
			},
		},
	}
}

func createRosettaTxFromMoveBalance(eTx *data.FullTransaction) *types.Transaction {
	tx := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: eTx.Hash,
		},
	}

	operations := make([]*types.Operation, 0)

	// check if transaction have value
	if eTx.Value != "0" {
		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type:   opTransfer,
			Status: OpStatusSuccess,
			Account: &types.AccountIdentifier{
				Address: eTx.Sender,
			},
			Amount: &types.Amount{
				Value:    "-" + eTx.Value,
				Currency: ElrondCurrency,
			},
		})

		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{Index: 0},
			},
			Type:   opTransfer,
			Status: OpStatusSuccess,
			Account: &types.AccountIdentifier{
				Address: eTx.Receiver,
			},
			Amount: &types.Amount{
				Value:    eTx.Value,
				Currency: ElrondCurrency,
			},
		})
	}

	// check if transaction have fee
	if eTx.GasLimit != 0 {
		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 2,
			},
			Type:   opFee,
			Status: OpStatusSuccess,
			Account: &types.AccountIdentifier{
				Address: eTx.Sender,
			},
			Amount: &types.Amount{
				Value:    "-" + computeTxFee(eTx),
				Currency: ElrondCurrency,
			},
		})
	}

	if len(operations) != 0 {
		tx.Operations = operations
	}

	return tx
}

func createOperationsFromPreparedTx(tx *data.Transaction) []*types.Operation {
	operations := make([]*types.Operation, 0)

	operations = append(operations, &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 0,
		},
		Type:   opTransfer,
		Status: OpStatusSuccess,
		Account: &types.AccountIdentifier{
			Address: tx.Sender,
		},
		Amount: &types.Amount{
			Value:    "-" + tx.Value,
			Currency: ElrondCurrency,
		},
	})

	operations = append(operations, &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 1,
		},
		RelatedOperations: []*types.OperationIdentifier{
			{Index: 0},
		},
		Type:   opTransfer,
		Status: OpStatusSuccess,
		Account: &types.AccountIdentifier{
			Address: tx.Receiver,
		},
		Amount: &types.Amount{
			Value:    tx.Value,
			Currency: ElrondCurrency,
		},
	})

	return operations
}

func computeTxFee(eTx *data.FullTransaction) string {
	fee := big.NewInt(0).SetUint64(eTx.GasPrice)
	fee.Mul(fee, big.NewInt(0).SetUint64(eTx.GasLimit))

	return fee.String()
}
