package data

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
)

// Transaction represents the structure that maps and validates user input for publishing a new transaction
type Transaction struct {
	// This field is used to tag transactions for send-multiple route
	Index     int    `json:"-"`
	Nonce     uint64 `form:"nonce" json:"nonce"`
	Value     string `form:"value" json:"value"`
	Receiver  string `form:"receiver" json:"receiver"`
	Sender    string `form:"sender" json:"sender"`
	GasPrice  uint64 `form:"gasPrice" json:"gasPrice,omitempty"`
	GasLimit  uint64 `form:"gasLimit" json:"gasLimit,omitempty"`
	Data      string `form:"data" json:"data,omitempty"`
	Signature string `form:"signature" json:"signature,omitempty"`
}

type GetTransactionResponse struct {
	Transaction transaction.ApiTransactionResult `json:"transaction"`
}

// transactionWrapper is a wrapper over a normal transaction in order to implement the interface needed in elrond-go
// for computing gas cost for a transaction
type transactionWrapper struct {
	transaction     *Transaction
	pubKeyConverter state.PubkeyConverter
}

// NewTransactionWrapper returns a new instance of transactionWrapper
func NewTransactionWrapper(transaction *Transaction, pubKeyConverter state.PubkeyConverter) (*transactionWrapper, error) {
	if transaction == nil {
		return nil, ErrNilTransaction
	}
	if pubKeyConverter == nil {
		return nil, ErrNilPubKeyConverter
	}

	return &transactionWrapper{
		transaction:     transaction,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// GetRcvAddr will return the receiver address in byte slice format
func (tw *transactionWrapper) GetRcvAddr() []byte {
	rcvrBytes, _ := tw.pubKeyConverter.Decode(tw.transaction.Receiver)
	return rcvrBytes
}

// GetGasLimit will return the gas limit of the tx
func (tw *transactionWrapper) GetGasLimit() uint64 {
	return tw.transaction.GasLimit
}

// GetGasPrice will return the gas price of the tx
func (tw *transactionWrapper) GetGasPrice() uint64 {
	return tw.transaction.GasPrice
}

// GetData will return the data of the tx
func (tw *transactionWrapper) GetData() []byte {
	return []byte(tw.transaction.Data)
}

// ResponseTransaction defines a response tx holding the resulting hash
type ResponseTransaction struct {
	TxHash string `json:"txHash"`
}

// ResponseMultipleTransactions defines a response from the node holding the number of transactions sent to the chain
type ResponseMultipleTransactions struct {
	NumOfTxs  uint64         `json:"txsSent"`
	TxsHashes map[int]string `json:"txsHashes"`
}

// ResponseTxCost defines a response from the node holding the transaction cost
type ResponseTxCost struct {
	TxCost uint64 `json:"txGasUnits"`
}

// ResponseTxStatus defines a response from the node holding the transaction status
type ResponseTxStatus struct {
	Status string `json:"status"`
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
