package data

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

// Transaction represents the structure that maps and validates user input for publishing a new transaction
type Transaction struct {
	// This field is used to tag transactions for send-multiple route
	Index             int    `json:"-"`
	Nonce             uint64 `json:"nonce"`
	Value             string `json:"value"`
	Receiver          string `json:"receiver"`
	Sender            string `json:"sender"`
	SenderUsername    []byte `json:"senderUsername,omitempty"`
	ReceiverUsername  []byte `json:"receiverUsername,omitempty"`
	GasPrice          uint64 `json:"gasPrice"`
	GasLimit          uint64 `json:"gasLimit"`
	Data              []byte `json:"data,omitempty"`
	Signature         string `json:"signature,omitempty"`
	ChainID           string `json:"chainID"`
	Version           uint32 `json:"version"`
	Options           uint32 `json:"options,omitempty"`
	GuardianAddr      string `json:"guardian,omitempty"`
	GuardianSignature string `json:"guardianSignature,omitempty"`
}

// GetTransactionResponseData follows the format of the data field of get transaction response
type GetTransactionResponseData struct {
	Transaction transaction.ApiTransactionResult `json:"transaction"`
}

// GetTransactionResponse defines a response from the node holding the transaction sent from the chain
type GetTransactionResponse struct {
	Data  GetTransactionResponseData `json:"data"`
	Error string                     `json:"error"`
	Code  string                     `json:"code"`
}

// GetSCRsResponseData follows the format of the data field of get smart contract results response
type GetSCRsResponseData struct {
	SCRs []*transaction.ApiSmartContractResult `json:"scrs"`
}

// GetSCRsResponse defines a response from the node holding the smart contract results
type GetSCRsResponse struct {
	Data  GetSCRsResponseData `json:"data"`
	Error string              `json:"error"`
	Code  string              `json:"code"`
}

// transactionWrapper is a wrapper over a normal transaction in order to implement the interface needed in mx-chain-go
// for computing gas cost for a transaction
type transactionWrapper struct {
	transaction     *Transaction
	pubKeyConverter core.PubkeyConverter
}

// NewTransactionWrapper returns a new instance of transactionWrapper
func NewTransactionWrapper(transaction *Transaction, pubKeyConverter core.PubkeyConverter) (*transactionWrapper, error) {
	if transaction == nil {
		return nil, ErrNilTransaction
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &transactionWrapper{
		transaction:     transaction,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// GetValue will return the value of the transaction
func (tw *transactionWrapper) GetValue() *big.Int {
	valueBigInt, ok := big.NewInt(0).SetString(tw.transaction.Value, 10)
	if !ok {
		return big.NewInt(0)
	}

	return valueBigInt
}

// GetRcvAddr will return the receiver address in byte slice format
func (tw *transactionWrapper) GetRcvAddr() []byte {
	rcvrBytes, _ := tw.pubKeyConverter.Decode(tw.transaction.Receiver)
	return rcvrBytes
}

// GetGasLimit will return the gas limit of the transaction
func (tw *transactionWrapper) GetGasLimit() uint64 {
	return tw.transaction.GasLimit
}

// GetGasPrice will return the gas price of the transaction
func (tw *transactionWrapper) GetGasPrice() uint64 {
	return tw.transaction.GasPrice
}

// GetData will return the data of the tx
func (tw *transactionWrapper) GetData() []byte {
	return tw.transaction.Data
}

// TransactionResponseData represents the format of the data field of a transaction response
type TransactionResponseData struct {
	TxHash string `json:"txHash"`
}

// ResponseTransaction defines a response tx holding the resulting hash
type ResponseTransaction struct {
	Data  TransactionResponseData `json:"data"`
	Error string                  `json:"error"`
	Code  string                  `json:"code"`
}

// TransactionSimulationResults holds the results of a transaction's simulation
type TransactionSimulationResults struct {
	Status     transaction.TxStatus                           `json:"status,omitempty"`
	FailReason string                                         `json:"failReason,omitempty"`
	ScResults  map[string]*transaction.ApiSmartContractResult `json:"scResults,omitempty"`
	Receipts   map[string]*transaction.ApiReceipt             `json:"receipts,omitempty"`
	Hash       string                                         `json:"hash,omitempty"`
}

// TransactionSimulationResponseData represents the format of the data field of a transaction simulation response
type TransactionSimulationResponseData struct {
	Result TransactionSimulationResults `json:"result"`
}

// ResponseTransactionSimulation defines a response tx holding the results of simulating a transaction execution
type ResponseTransactionSimulation struct {
	Data  TransactionSimulationResponseData `json:"data"`
	Error string                            `json:"error"`
	Code  ReturnCode                        `json:"code"`
}

// TransactionSimulationResponseDataCrossShard represents the format of the data field of a transaction simulation response in cross shard transactions
type TransactionSimulationResponseDataCrossShard struct {
	Result map[string]TransactionSimulationResults `json:"result"`
}

// ResponseTransactionSimulationCrossShard defines a response tx holding the results of simulating a transaction execution in a cross-shard way
type ResponseTransactionSimulationCrossShard struct {
	Data  TransactionSimulationResponseDataCrossShard `json:"data"`
	Error string                                      `json:"error"`
	Code  ReturnCode                                  `json:"code"`
}

// MultipleTransactionsResponseData holds the data which is returned when sending a bulk of transactions
type MultipleTransactionsResponseData struct {
	NumOfTxs  uint64         `json:"txsSent"`
	TxsHashes map[int]string `json:"txsHashes"`
}

// ResponseMultipleTransactions defines a response from the node holding the number of transactions sent to the chain
type ResponseMultipleTransactions struct {
	Data  MultipleTransactionsResponseData `json:"data"`
	Error string                           `json:"error"`
	Code  string                           `json:"code"`
}

// TxCostResponseData follows the format of the data field of a transaction cost request
type TxCostResponseData struct {
	TxCost     uint64                                     `json:"txGasUnits"`
	RetMessage string                                     `json:"returnMessage"`
	ScResults  map[string]*ExtendedApiSmartContractResult `json:"smartContractResults"`
	Logs       *transaction.ApiLogs                       `json:"logs,omitempty"`
}

// ExtendedApiSmartContractResult extends the structure transaction.ApiSmartContractResult with an extra field
type ExtendedApiSmartContractResult struct {
	*transaction.ApiSmartContractResult
	Used bool `json:"-"`
}

// ResponseTxCost defines a response from the node holding the transaction cost
type ResponseTxCost struct {
	Data  TxCostResponseData `json:"data"`
	Error string             `json:"error"`
	Code  string             `json:"code"`
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

// WrappedTransaction represents a wrapped transaction that is received from tx pool
type WrappedTransaction struct {
	TxFields map[string]interface{} `json:"txFields"`
}

// TransactionsPool represents a structure that holds all wrapped transactions from pool
type TransactionsPool struct {
	RegularTransactions  []WrappedTransaction `json:"regularTransactions"`
	SmartContractResults []WrappedTransaction `json:"smartContractResults"`
	Rewards              []WrappedTransaction `json:"rewards"`
}

// TransactionsPoolResponseData matches the data field of get tx pool response
type TransactionsPoolResponseData struct {
	Transactions TransactionsPool `json:"txPool"`
}

// TransactionsPoolApiResponse matches the output of an observer's tx pool endpoint
type TransactionsPoolApiResponse struct {
	Data  TransactionsPoolResponseData `json:"data"`
	Error string                       `json:"error"`
	Code  string                       `json:"code"`
}

// TransactionsPoolForSender represents a structure that holds wrapped transactions from pool for a sender
type TransactionsPoolForSender struct {
	Transactions []WrappedTransaction `json:"transactions"`
}

// TransactionsPoolForSenderResponseData matches the data field of get tx pool for sender response
type TransactionsPoolForSenderResponseData struct {
	TxPool TransactionsPoolForSender `json:"txPool"`
}

// TransactionsPoolForSenderApiResponse matches the output of an observer's tx pool for sender endpoint
type TransactionsPoolForSenderApiResponse struct {
	Data  TransactionsPoolForSenderResponseData `json:"data"`
	Error string                                `json:"error"`
	Code  string                                `json:"code"`
}

// TransactionsPoolLastNonceForSender matches the data field of get last nonce from pool for sender response
type TransactionsPoolLastNonceForSender struct {
	Nonce uint64 `json:"nonce"`
}

// TransactionsPoolLastNonceForSenderApiResponse matches the output of an observer's last nonce from tx pool for sender endpoint
type TransactionsPoolLastNonceForSenderApiResponse struct {
	Data  TransactionsPoolLastNonceForSender `json:"data"`
	Error string                             `json:"error"`
	Code  string                             `json:"code"`
}

// NonceGap represents a struct that holds a nonce gap from tx pool
// From - first unknown nonce
// To   - last unknown nonce
type NonceGap struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
}

// TransactionsPoolNonceGaps represents a structure that holds nonce gaps
type TransactionsPoolNonceGaps struct {
	Gaps []NonceGap `json:"gaps"`
}

// TransactionsPoolNonceGapsForSenderResponseData matches the data field of get nonce gaps from tx pool for sender response
type TransactionsPoolNonceGapsForSenderResponseData struct {
	NonceGaps TransactionsPoolNonceGaps `json:"nonceGaps"`
}

// TransactionsPoolNonceGapsForSenderApiResponse matches the output of an observer's nonce gaps from tx pool for sender endpoint
type TransactionsPoolNonceGapsForSenderApiResponse struct {
	Data  TransactionsPoolNonceGapsForSenderResponseData `json:"data"`
	Error string                                         `json:"error"`
	Code  string                                         `json:"code"`
}

// ProcessStatusResponse represents a structure that holds the process status of a transaction
type ProcessStatusResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}
