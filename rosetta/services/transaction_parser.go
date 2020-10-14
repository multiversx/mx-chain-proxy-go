package services

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type transactionsParser struct {
	config         *configuration.Configuration
	networkConfig  *provider.NetworkConfig
	elrondProvider provider.ElrondProviderHandler
}

func newTransactionParser(
	provider provider.ElrondProviderHandler,
	cfg *configuration.Configuration,
	networkConfig *provider.NetworkConfig,
) *transactionsParser {
	return &transactionsParser{
		config:         cfg,
		networkConfig:  networkConfig,
		elrondProvider: provider,
	}
}

func (tp *transactionsParser) parseTxsFromHyperBlock(hyperBlock *data.Hyperblock) []*types.Transaction {
	txs := make([]*types.Transaction, 0)
	for _, eTx := range hyperBlock.Transactions {
		tx, ok := tp.parseTx(eTx, false)
		if !ok {
			continue
		}
		txs = append(txs, tx)
	}

	return txs
}

func (tp *transactionsParser) parseTx(eTx *data.FullTransaction, isInPool bool) (*types.Transaction, bool) {
	switch eTx.Type {
	case string(transaction.TxTypeNormal):
		return tp.createRosettaTxFromMoveBalance(eTx, isInPool), true
	case string(transaction.TxTypeReward):
		return tp.createRosettaTxFromReward(eTx), true
	case string(transaction.TxTypeUnsigned):
		return tp.createRosettaTxFromUnsignedTx(eTx)
	case string(transaction.TxTypeInvalid):
		return tp.createRosettaTxFromInvalidTx(eTx), true
	default:
		return nil, false
	}
}

func (tp *transactionsParser) createRosettaTxFromUnsignedTx(eTx *data.FullTransaction) (*types.Transaction, bool) {
	// TODO check if we have a SCR that calls another contract
	if eTx.Value == "0" {
		return nil, false
	}

	switch {
	case eTx.GasLimit != 0:
		// we have a SCR with gas refund
		return tp.createRosettaTxWithGasRefund(eTx)
	case eTx.Sender != eTx.Receiver:
		// we have a SCR with send funds
		return tp.createRosettaTxUnsignedTxSendFunds(eTx)
	default:
		return nil, false
	}
}

func (tp *transactionsParser) createRosettaTxWithGasRefund(eTx *data.FullTransaction) (*types.Transaction, bool) {
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
					Currency: tp.config.Currency,
				},
			},
		},
	}, true
}

func (tp *transactionsParser) createRosettaTxUnsignedTxSendFunds(
	eTx *data.FullTransaction,
) (*types.Transaction, bool) {
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
					Address: eTx.Sender,
				},
				Amount: &types.Amount{
					Value:    "-" + eTx.Value,
					Currency: tp.config.Currency,
				},
			},
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
					Currency: tp.config.Currency,
				},
			},
		},
	}, true
}

func (tp *transactionsParser) createRosettaTxFromReward(eTx *data.FullTransaction) *types.Transaction {
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
					Currency: tp.config.Currency,
				},
			},
		},
	}
}

func (tp *transactionsParser) createRosettaTxFromMoveBalance(eTx *data.FullTransaction, isInPool bool) *types.Transaction {
	tx := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: eTx.Hash,
		},
	}

	operations := make([]*types.Operation, 0)

	// check if transaction has value
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
				Currency: tp.config.Currency,
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
				Currency: tp.config.Currency,
			},
		})
	}

	// check if transaction has fee and transaction is not in pool
	if eTx.GasLimit != 0 && !isInPool {
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
				Value:    "-" + tp.computeTxFee(eTx),
				Currency: tp.config.Currency,
			},
		})
	}

	if len(operations) != 0 {
		tx.Operations = operations
	}

	return tx
}

func (tp *transactionsParser) createOperationsFromPreparedTx(tx *data.Transaction) []*types.Operation {
	operations := make([]*types.Operation, 0)

	operations = append(operations, &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 0,
		},
		Type: opTransfer,
		Account: &types.AccountIdentifier{
			Address: tx.Sender,
		},
		Amount: &types.Amount{
			Value:    "-" + tx.Value,
			Currency: tp.config.Currency,
		},
	})

	operations = append(operations, &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: 1,
		},
		RelatedOperations: []*types.OperationIdentifier{
			{Index: 0},
		},
		Type: opTransfer,
		Account: &types.AccountIdentifier{
			Address: tx.Receiver,
		},
		Amount: &types.Amount{
			Value:    tx.Value,
			Currency: tp.config.Currency,
		},
	})

	return operations
}

func (tp *transactionsParser) createRosettaTxFromInvalidTx(eTx *data.FullTransaction) *types.Transaction {
	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: eTx.Hash,
		},
		Operations: []*types.Operation{
			{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opInvalid,
				Status: OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: eTx.Sender,
				},
				Amount: &types.Amount{
					Value:    "-" + tp.computeTxFee(eTx),
					Currency: tp.config.Currency,
				},
			},
		},
	}
}

func (tp *transactionsParser) computeTxFee(eTx *data.FullTransaction) string {
	// errors should never appear because transaction is included in a block and receiver address is a valid address
	decodedAddr, _ := tp.elrondProvider.DecodeAddress(eTx.Receiver)

	if core.IsSmartContractAddress(decodedAddr) {
		// we have a smart contract call
		fee := big.NewInt(0)
		fee.Mul(big.NewInt(0).SetUint64(eTx.GasPrice), big.NewInt(0).SetUint64(eTx.GasLimit))

		return fee.String()
	}

	gasPrice := eTx.GasPrice
	gasLimit := tp.networkConfig.MinGasLimit + uint64(len(eTx.Data))*tp.networkConfig.GasPerDataByte

	fee := big.NewInt(0).SetUint64(gasPrice)
	fee.Mul(fee, big.NewInt(0).SetUint64(gasLimit))

	return fee.String()
}
