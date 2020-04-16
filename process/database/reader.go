package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/elastic/go-elasticsearch/v7"
)

const numTxs = 20

type reader struct {
	client *elasticsearch.Client
}

// NewDatabaseReader create a new elastic search database reader object
func NewDatabaseReader(url, username, password string) (*reader, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create database reader %w", err)
	}

	return &reader{
		client: client,
	}, nil
}

// GetTransactionsByAddress will read from elasticsearch server all transaction that have senderAddress and destinationAddress
// equals with provided address
func (r *reader) GetTransactionsByAddress(address string) ([]data.ApiTransaction, error) {
	query := txsByAddrQuery(address)
	decodedBody, err := r.doSearchRequest(query, "transactions", numTxs)
	if err != nil {
		return nil, err
	}

	hits, ok := decodedBody["hits"].(map[string]interface{})
	if !ok {
		return nil, errors.New("cannot get data from response body")
	}

	txs := formatTxs(hits)
	return txs, nil
}

func (r *reader) GetLatestBlockHeight() (uint64, error) {
	query := latestBlockQuery()
	decodedBody, err := r.doSearchRequest(query, "blocks", 1)
	if err != nil {
		return 0, err
	}

	block, err := formatBlock(decodedBody)
	if err != nil {
		return 0, err
	}

	return block.Nonce, nil
}

func (r *reader) doSearchRequest(query map[string]interface{}, index string, size int) (map[string]interface{}, error) {
	buff, err := encodeQuery(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex(index),
		r.client.Search.WithSize(size),
		r.client.Search.WithBody(&buff),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get data from database: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()
	if res.IsError() {
		return nil, fmt.Errorf("cannot get data from database: %v", res)
	}

	var decodedBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&decodedBody); err != nil {
		return nil, err
	}

	return decodedBody, nil
}

func formatBlock(d map[string]interface{}) (*indexer.Block, error) {
	h1 := d["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(h1) == 0 {
		return nil, fmt.Errorf("cannot find blocks in database")
	}

	h2 := h1[0].(map[string]interface{})["_source"]

	bbb, _ := json.Marshal(h2)
	var block indexer.Block
	err := json.Unmarshal(bbb, &block)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal block")
	}

	return &block, nil
}

func formatTxs(d map[string]interface{}) []data.ApiTransaction {
	var err error

	txs := make([]data.ApiTransaction, 0)
	for _, h1 := range d["hits"].([]interface{}) {
		h2 := h1.(map[string]interface{})["_source"]

		var tx data.ApiTransaction
		bbb, _ := json.Marshal(h2)
		err = json.Unmarshal(bbb, &tx)
		if err != nil {
			continue
		}

		txs = append(txs, tx)
	}
	return txs
}
