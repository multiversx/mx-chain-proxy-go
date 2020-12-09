package database

import (
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	numTopTransactions           = 20
	numTransactionFromAMiniblock = 100
)

type elasticSearchConnector struct {
	client *elasticsearch.Client
}

// NewElasticSearchConnector create a new elastic search database reader object
func NewElasticSearchConnector(url, username, password string) (*elasticSearchConnector, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create database reader %w", err)
	}

	return &elasticSearchConnector{
		client: client,
	}, nil
}

// GetTransactionsByAddress gets transactions TO or FROM the specified address
func (esc *elasticSearchConnector) GetTransactionsByAddress(address string) ([]data.DatabaseTransaction, error) {
	decodedBody, err := esc.doSearchRequestTx(address, "transactions", numTopTransactions)
	if err != nil {
		return nil, err
	}

	return convertObjectToTransactions(decodedBody)
}

// GetAtlasBlockByShardIDAndNonce gets from database a block with the specified shardID and nonce
func (esc *elasticSearchConnector) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	query := blockByNonceAndShardIDQuery(nonce, shardID)
	decodedBody, err := esc.doSearchRequest(query, "blocks", 1)
	if err != nil {
		return data.AtlasBlock{}, err
	}

	metaBlock, metaBlockHash, err := convertObjectToBlock(decodedBody)
	if err != nil {
		return data.AtlasBlock{}, err
	}

	txs, err := esc.getTxsByMiniblockHashes(metaBlock.MiniBlocksHashes)
	if err != nil {
		return data.AtlasBlock{}, err
	}

	transactions, err := esc.getTxsByNotarizedBlockHashes(metaBlock.NotarizedBlocksHashes)
	if err != nil {
		return data.AtlasBlock{}, err
	}

	txs = append(txs, transactions...)

	return data.AtlasBlock{
		Nonce:        metaBlock.Nonce,
		Hash:         metaBlockHash,
		Transactions: txs,
	}, nil
}

func (esc *elasticSearchConnector) getTxsByNotarizedBlockHashes(hashes []string) ([]data.DatabaseTransaction, error) {
	txs := make([]data.DatabaseTransaction, 0)
	for _, hash := range hashes {
		query := blockByHashQuery(hash)
		decodedBody, err := esc.doSearchRequest(query, "blocks", 1)
		if err != nil {
			return nil, err
		}

		shardBlock, _, err := convertObjectToBlock(decodedBody)
		if err != nil {
			return nil, err
		}

		transactions, err := esc.getTxsByMiniblockHashes(shardBlock.MiniBlocksHashes)
		if err != nil {
			return nil, err
		}

		txs = append(txs, transactions...)
	}
	return txs, nil
}

func (esc *elasticSearchConnector) getTxsByMiniblockHashes(hashes []string) ([]data.DatabaseTransaction, error) {
	txs := make([]data.DatabaseTransaction, 0)
	for _, hash := range hashes {
		query := txsByMiniblockHashQuery(hash)
		decodedBody, err := esc.doSearchRequest(query, "transactions", numTransactionFromAMiniblock)
		if err != nil {
			return nil, err
		}

		transactions, err := convertObjectToTransactions(decodedBody)
		if err != nil {
			return nil, err
		}

		txs = append(txs, transactions...)
	}
	return txs, nil
}

func (esc *elasticSearchConnector) doSearchRequest(query object, index string, size int) (object, error) {
	buff, err := encodeQuery(query)
	if err != nil {
		return nil, err
	}

	res, err := esc.client.Search(
		esc.client.Search.WithIndex(index),
		esc.client.Search.WithSize(size),
		esc.client.Search.WithBody(&buff),
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

func (esc *elasticSearchConnector) doSearchRequestTx(address string, index string, size int) (object, error) {
	res, err := esc.client.Search(
		esc.client.Search.WithIndex(index),
		esc.client.Search.WithSize(size),
		esc.client.Search.WithQuery("sender%"+address+"+receiver%"+address),
		esc.client.Search.WithSort("timestamp:desc"),
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

// IsInterfaceNil returns true if there is no value under the interface
func (esc *elasticSearchConnector) IsInterfaceNil() bool {
	return esc == nil
}
