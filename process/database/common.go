package database

import (
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func formatBlock(d map[string]interface{}) (*indexer.Block, string, error) {
	h1 := d["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(h1) == 0 {
		return nil, "", fmt.Errorf("cannot find blocks in database")
	}
	h2 := h1[0].(map[string]interface{})["_source"]

	h3 := h1[0].(map[string]interface{})["_id"]
	blockHash := fmt.Sprint(h3)

	bbb, _ := json.Marshal(h2)
	var block indexer.Block
	err := json.Unmarshal(bbb, &block)
	if err != nil {
		return nil, "", fmt.Errorf("cannot unmarshal block")
	}

	return &block, blockHash, nil
}

func formatTxs(d map[string]interface{}) ([]data.DatabaseTransaction, error) {
	hits, ok := d["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot get transactions from decoded body")
	}

	txs := make([]data.DatabaseTransaction, 0)
	for _, h1 := range hits["hits"].([]interface{}) {
		h2 := h1.(map[string]interface{})["_source"]

		var tx data.DatabaseTransaction
		bbb, _ := json.Marshal(h2)
		err := json.Unmarshal(bbb, &tx)
		if err != nil {
			continue
		}

		h3 := h1.(map[string]interface{})["_id"]
		txHash := fmt.Sprint(h3)
		tx.Hash = txHash

		txs = append(txs, tx)
	}
	return txs, nil
}
