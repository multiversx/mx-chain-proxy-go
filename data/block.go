package data

import "github.com/ElrondNetwork/elrond-go/core/indexer"

// ApiBlock represents the structure of a block as it is sent to API
type ApiBlock struct {
	Nonce        uint64                `form:"nonce" json:"nonce"`
	Hash         string                `form:"hash" json:"hash"`
	Transactions []indexer.Transaction `form:"transactions" json:"transactions"`
}
