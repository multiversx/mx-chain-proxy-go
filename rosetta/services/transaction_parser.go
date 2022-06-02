package services

import (
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
	nodeTxs := filterOutIntrashardContractResultsWhoseOriginalTransactionIsInInvalidMiniblock(hyperBlock.Transactions)
	nodeTxs = filterOutIntrashardRelayedTransactionAlreadyHeldInInvalidMiniblock(nodeTxs)
	// nodeTxs = filterOutIntraMetachainTransactions(nodeTxs)

	txs := make([]*types.Transaction, 0)
	for _, eTx := range nodeTxs {
		tx, ok := tp.parseTx(eTx, false)
		if !ok {
			continue
		}

		// TODO: Should we populate related transactions?
		// populateRelatedTransactions(tx, eTx)
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
	if eTx.Value == "0" {
		return nil, false
	}
	if eTx.Value[0] == '-' {
		return nil, false
	}

	if eTx.IsRefund {
		return tp.createRosettaTxWithGasRefund(eTx)
	} else {
		return tp.createRosettaTxUnsignedTxSendFunds(eTx)
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
				Status: &OpStatusSuccess,
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
				Status: &OpStatusSuccess,
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
					Index: 1,
				},
				Type:   opScResult,
				Status: &OpStatusSuccess,
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
				Status: &OpStatusSuccess,
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
	hasValue := eTx.Value != "0"
	isFromMetachain := eTx.SourceShard == core.MetachainShardId
	isToMetachain := eTx.DestinationShard == core.MetachainShardId

	tx := &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: eTx.Hash,
		},
	}

	operations := make([]*types.Operation, 0)
	operationIndex := int64(0)

	if hasValue {
		if !isFromMetachain {
			operations = append(operations, &types.Operation{
				OperationIdentifier: &types.OperationIdentifier{
					Index: 0,
				},
				Type:   opTransfer,
				Status: &OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: eTx.Sender,
				},
				Amount: &types.Amount{
					Value:    "-" + eTx.Value,
					Currency: tp.config.Currency,
				},
			})

			operationIndex++
		}

		if !isToMetachain {
			operations = append(operations, &types.Operation{
				OperationIdentifier: &types.OperationIdentifier{
					Index: operationIndex,
				},
				RelatedOperations: []*types.OperationIdentifier{
					{Index: 0},
				},
				Type:   opTransfer,
				Status: &OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: eTx.Receiver,
				},
				Amount: &types.Amount{
					Value:    eTx.Value,
					Currency: tp.config.Currency,
				},
			})

			operationIndex++
		}
	}

	// check if transaction has fee and transaction is not in pool
	// TODO / QUESTION for review: can it <not have fee>? can gas limit be 0?
	// TODO: also, why not declare fee as well if it's in pool?
	if eTx.GasLimit != 0 && !isInPool {
		operations = append(operations, &types.Operation{
			OperationIdentifier: &types.OperationIdentifier{
				Index: operationIndex,
			},
			Type:   opFee,
			Status: &OpStatusSuccess,
			Account: &types.AccountIdentifier{
				Address: eTx.Sender,
			},
			Amount: &types.Amount{
				Value:    "-" + eTx.InitiallyPaidFee,
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
				// TODO: how to handle this? Also specify types in NetworkOptionsResponse.
				Type:   opInvalid,
				Status: &OpStatusSuccess,
				Account: &types.AccountIdentifier{
					Address: eTx.Sender,
				},
				Amount: &types.Amount{
					Value:    "-" + eTx.InitiallyPaidFee,
					Currency: tp.config.Currency,
				},
			},
		},
	}
}

func populateRelatedTransactions(rosettaTx *types.Transaction, nodeTx *data.FullTransaction) {
	if nodeTx.OriginalTransactionHash != "" {
		rosettaTx.RelatedTransactions = append(rosettaTx.RelatedTransactions, &types.RelatedTransaction{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: nodeTx.OriginalTransactionHash,
			},
			Direction: types.Backward,
		})
	}

	if nodeTx.PreviousTransactionHash != "" && nodeTx.PreviousTransactionHash != nodeTx.OriginalTransactionHash {
		rosettaTx.RelatedTransactions = append(rosettaTx.RelatedTransactions, &types.RelatedTransaction{
			TransactionIdentifier: &types.TransactionIdentifier{
				Hash: nodeTx.PreviousTransactionHash,
			},
			Direction: types.Backward,
		})
	}
}
