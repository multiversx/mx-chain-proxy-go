package database

import (
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
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

	return formatTxs(decodedBody)
}

func (r *reader) GetLatestBlockHeight() (uint64, error) {
	query := latestBlockQuery()
	decodedBody, err := r.doSearchRequest(query, "blocks", 1)
	if err != nil {
		return 0, err
	}

	block, _, err := formatBlock(decodedBody)
	if err != nil {
		return 0, err
	}

	return block.Nonce, nil
}

// GetBlockByNonce -
func (r *reader) GetBlockByNonce(nonce uint64) (data.ApiBlock, error) {
	query := blockByNonceAndShardIDQuery(nonce, core.MetachainShardId)
	decodedBody, err := r.doSearchRequest(query, "blocks", 1)
	if err != nil {
		return data.ApiBlock{}, err
	}

	block, blockHash, err := formatBlock(decodedBody)
	if err != nil {
		return data.ApiBlock{}, err
	}

	txs, err := r.getTxsByMiniblockHashes(block.MiniBlocksHashes)
	if err != nil {
		return data.ApiBlock{}, err
	}

	transactions, err := r.getTxsByNotarizedBlockHashes(block.NotarizedBlocksHashes)
	if err != nil {
		return data.ApiBlock{}, err
	}

	txs = append(txs, transactions...)

	return data.ApiBlock{
		Nonce:        block.Nonce,
		Hash:         blockHash,
		Transactions: txs,
	}, nil
}

func (r *reader) getTxsByNotarizedBlockHashes(hashes []string) ([]data.ApiTransaction, error) {
	txs := make([]data.ApiTransaction, 0)
	for _, notarizedBlocKHash := range hashes {
		query := blockByHashQuery(notarizedBlocKHash)
		decodedBody, err := r.doSearchRequest(query, "blocks", 1)
		if err != nil {
			return nil, err
		}

		shardBlock, _, err := formatBlock(decodedBody)
		if err != nil {
			return nil, err
		}

		transactions, err := r.getTxsByMiniblockHashes(shardBlock.MiniBlocksHashes)
		if err != nil {
			return nil, err
		}

		txs = append(txs, transactions...)
	}
	return txs, nil
}

func (r *reader) getTxsByMiniblockHashes(hashes []string) ([]data.ApiTransaction, error) {
	txs := make([]data.ApiTransaction, 0)
	for _, hash := range hashes {
		query := txsByMiniblockHashQuery(hash)
		decodedBody, err := r.doSearchRequest(query, "transactions", 100)
		if err != nil {
			return nil, err
		}

		transactions, err := formatTxs(decodedBody)
		if err != nil {
			return nil, err
		}

		txs = append(txs, transactions...)
	}
	return txs, nil
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

func (r *reader) IsInterfaceNil() bool {
	return r == nil
}
