package data

import "math/big"

// Transaction represents the structure that maps and validates user input for publishing a new transaction
type Transaction struct {
	Sender    string   `form:"sender" json:"sender"`
	Receiver  string   `form:"receiver" json:"receiver"`
	Value     *big.Int `form:"value" json:"value"`
	Data      string   `form:"data" json:"data"`
	Nonce     uint64   `form:"nonce" json:"nonce"`
	GasPrice  *big.Int `form:"gasPrice" json:"gasPrice"`
	GasLimit  *big.Int `form:"gasLimit" json:"gasLimit"`
	Signature string   `form:"signature" json:"signature"`
	Challenge string   `form:"challenge" json:"challenge"`
}
