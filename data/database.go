package data

import "github.com/ElrondNetwork/elrond-go/core/indexer"

// DatabaseTransaction extends indexer.Transaction with the 'hash' field that is not ignored in json schema
type DatabaseTransaction struct {
	Hash string `json:"hash"`
	indexer.Transaction
}
