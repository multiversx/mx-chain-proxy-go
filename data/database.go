package data

import (
	"math/big"

	"github.com/multiversx/mx-chain-es-indexer-go/data"
)

// DatabaseTransaction extends indexer.Transaction with the 'hash' field that is not ignored in json schema
type DatabaseTransaction struct {
	Hash string `json:"hash"`
	Fee  string `json:"fee"`
	data.Transaction
}

// CalculateFee calculates transaction fee using gasPrice and gasUsed
func (dt *DatabaseTransaction) CalculateFee() string {
	gasPrice := big.NewInt(0).SetUint64(dt.GasPrice)
	gasUsed := big.NewInt(0).SetUint64(dt.GasUsed)
	fee := big.NewInt(0).Mul(gasPrice, gasUsed)

	return fee.String()
}
