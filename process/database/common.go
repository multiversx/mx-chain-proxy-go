package database

import (
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func convertObjectToBlock(d object) (*indexer.Block, string, error) {
	h1 := d["hits"].(object)["hits"].([]interface{})
	if len(h1) == 0 {
		return nil, "", fmt.Errorf("cannot find blocks in database")
	}
	h2 := h1[0].(object)["_source"]

	h3 := h1[0].(object)["_id"]
	blockHash := fmt.Sprint(h3)

	marshalizedBlock, _ := json.Marshal(h2)
	var block indexer.Block
	err := json.Unmarshal(marshalizedBlock, &block)
	if err != nil {
		return nil, "", fmt.Errorf("cannot unmarshal block")
	}

	return &block, blockHash, nil
}

func convertObjectToTransactions(d object) ([]data.DatabaseTransaction, error) {
	hits, ok := d["hits"].(object)
	if !ok {
		return nil, fmt.Errorf("cannot get transactions from decoded body")
	}

	txs := make([]data.DatabaseTransaction, 0)
	for _, h1 := range hits["hits"].([]interface{}) {
		h2 := h1.(object)["_source"]

		var tx indexer.Transaction
		marshalizedBlock, _ := json.Marshal(h2)
		err := json.Unmarshal(marshalizedBlock, &tx)
		if err != nil {
			continue
		}

		h3 := h1.(object)["_id"]
		txHash := fmt.Sprint(h3)
		tx.Hash = txHash

		txs = append(txs, convertToDatabaseTransaction(tx))
	}
	return txs, nil
}

func convertToDatabaseTransaction(srcTx indexer.Transaction) data.DatabaseTransaction {
	return data.DatabaseTransaction{
		Hash:          srcTx.Hash,
		MBHash:        srcTx.MBHash,
		Nonce:         srcTx.Nonce,
		Round:         srcTx.Round,
		Value:         srcTx.Value,
		Receiver:      srcTx.Receiver,
		Sender:        srcTx.Sender,
		ReceiverShard: srcTx.ReceiverShard,
		SenderShard:   srcTx.SenderShard,
		GasPrice:      srcTx.GasPrice,
		GasLimit:      srcTx.GasLimit,
		Data:          srcTx.Data,
		Signature:     srcTx.Signature,
		Timestamp:     srcTx.Timestamp,
		Status:        srcTx.Status,
		Fee:           srcTx.GasUsed,
	}
}
