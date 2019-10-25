package data

import "math/big"

// Transaction represents the structure that maps and validates user input for publishing a new transaction
type Transaction struct {
	Nonce     uint64 `form:"nonce" json:"nonce"`
	Value     string `form:"value" json:"value"`
	Receiver  string `form:"receiver" json:"receiver"`
	Sender    string `form:"sender" json:"sender"`
	GasPrice  uint64 `form:"gasPrice" json:"gasPrice,omitempty"`
	GasLimit  uint64 `form:"gasLimit" json:"gasLimit,omitempty"`
	Data      string `form:"data" json:"data,omitempty"`
	Signature string `form:"signature" json:"signature,omitempty"`
	Challenge string `form:"challenge" json:"challenge,omitempty"`
}

// ResponseTransaction defines a response tx holding the resulting hash
type ResponseTransaction struct {
	TxHash string `json:"txHash"`
}

// ResponseMultiTransactions defines a response from the node holding the number of transactions sent to the chain
type ResponseMultiTransactions struct {
	NumOfTxs uint64 `json:"txsSent"`
}

// FundsRequest represents the data structure needed as input for sending funds from a node to an address
type FundsRequest struct {
	Receiver string   `form:"receiver" json:"receiver"`
	Value    *big.Int `form:"value" json:"value,omitempty"`
	TxCount  int      `form:"txCount" json:"txCount,omitempty"`
}

// ResponseFunds defines the response structure for the node's generate-and-send-multiple endpoint
type ResponseFunds struct {
	Message string `json:"message"`
}
