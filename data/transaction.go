package data

import "math/big"

// Transaction represents the structure that maps and validates user input for publishing a new transaction
type Transaction struct {
	Nonce     uint64   `form:"nonce" json:"nonce"`
	Value     *big.Int `form:"value" json:"value"`
	Receiver  string   `form:"receiver" json:"receiver"`
	Sender    string   `form:"sender" json:"sender"`
	GasPrice  *big.Int `form:"gasPrice" json:"gasPrice,omitempty"`
	GasLimit  *big.Int `form:"gasLimit" json:"gasLimit,omitempty"`
	Data      string   `form:"data" json:"data,omitempty"`
	Signature string   `form:"signature" json:"signature,omitempty"`
	Challenge string   `form:"challenge" json:"challenge,omitempty"`
}
