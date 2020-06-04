package process

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the transaction group path of the node
const TransactionPath = "/transaction/"

// TransactionSendPath defines the single transaction send path of the node
const TransactionSendPath = "/transaction/send"

// MultipleTransactionsPath defines the multiple transactions send path of the node
const MultipleTransactionsPath = "/transaction/send-multiple"

// TransactionCostPath defines the transaction's cost path of the node
const TransactionCostPath = "/transaction/cost"

// UnknownStatusTx defines the response that should be received from an observer when transaction status is unknown
const UnknownStatusTx = "unknown"

type erdTransaction struct {
	Nonce     uint64 `json:"nonce"`
	Value     string `json:"value"`
	RcvAddr   string `json:"receiver"`
	SndAddr   string `json:"sender"`
	GasPrice  uint64 `json:"gasPrice,omitempty"`
	GasLimit  uint64 `json:"gasLimit,omitempty"`
	Data      string `json:"data,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc            Processor
	pubKeyConverter core.PubkeyConverter
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(
	proc Processor,
	pubKeyConverter core.PubkeyConverter,
) (*TransactionProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &TransactionProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// SendTransaction relay the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SendTransaction(tx *data.Transaction) (int, string, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	senderBuff, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	shardID, err := tp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	observers, err := tp.proc.GetObservers(shardID)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	for _, observer := range observers {
		txResponse := &data.ResponseTransaction{}

		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, TransactionSendPath, tx, txResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("Transaction sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				shardID,
				txResponse.TxHash,
			))
			return respCode, txResponse.TxHash, nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return respCode, "", err
	}

	return http.StatusInternalServerError, "", ErrSendingRequest
}

// SendMultipleTransactions relay the post request by sending the request to the first available observer and replies back the answer

func (tp *TransactionProcessor) SendMultipleTransactions(txs []*data.Transaction) (
	data.ResponseMultipleTransactions, error,
) {
	//TODO: Analyze and improve the robustness of this function. Currently, an error within `GetObservers`
	//breaks the function and returns nothing (but an error) even if some transactions were actually sent, successfully.

	totalTxsSent := uint64(0)
	txsToSend := make([]*data.Transaction, 0)
	for i := 0; i < len(txs); i++ {
		currentTx := txs[i]
		err := tp.checkTransactionFields(currentTx)
		if err != nil {
			log.Warn("invalid tx received",
				"sender", currentTx.Sender,
				"receiver", currentTx.Receiver,
				"error", err)
			continue
		}
		txsToSend = append(txsToSend, currentTx)
	}
	if len(txsToSend) == 0 {
		return data.ResponseMultipleTransactions{}, ErrNoValidTransactionToSend
	}

	txsHashes := make(map[int]string)
	txsByShardID := tp.groupTxsByShard(txsToSend)
	for shardID, groupOfTxs := range txsByShardID {
		observersInShard, err := tp.proc.GetObservers(shardID)
		if err != nil {
			return data.ResponseMultipleTransactions{}, ErrMissingObserver
		}

		for _, observer := range observersInShard {
			txResponse := &data.ResponseMultipleTransactions{}
			respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, MultipleTransactionsPath, groupOfTxs, txResponse)
			if respCode == http.StatusOK && err == nil {
				log.Info("transactions sent",
					"observer", observer.Address,
					"shard ID", shardID,
					"total processed", txResponse.NumOfTxs,
				)
				totalTxsSent += txResponse.NumOfTxs

				for key, hash := range txResponse.TxsHashes {
					txsHashes[groupOfTxs[key].Index] = hash
				}

				break
			}

			log.LogIfError(err)
		}
	}

	return data.ResponseMultipleTransactions{
		NumOfTxs:  totalTxsSent,
		TxsHashes: txsHashes,
	}, nil
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (tp *TransactionProcessor) TransactionCostRequest(tx *data.Transaction) (string, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return "", err
	}

	observers := tp.proc.GetAllObservers()
	for _, observer := range observers {
		if observer.ShardId == core.MetachainShardId {
			continue
		}

		txCostResponse := &data.ResponseTxCost{}
		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, TransactionCostPath, tx, txCostResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info("calculate tx cost request was sent successfully",
				"observer ", observer.Address,
				"shard", observer.ShardId,
			)
			return strconv.Itoa(int(txCostResponse.TxCost)), nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return "", err

	}

	return "", ErrSendingRequest
}

// GetTransaction should return a transaction from observer
func (tp *TransactionProcessor) GetTransaction(txHash string) (*transaction.ApiTransactionResult, error) {
	var err error

	observers := tp.proc.GetAllObservers()
	for _, observer := range observers {
		getTxResponse := &data.GetTransactionResponse{}
		err = tp.proc.CallGetRestEndPoint(observer.Address, TransactionPath+txHash, getTxResponse)
		if err != nil {
			log.Trace("cannot get transaction", "error", err)
			continue
		}

		return &getTxResponse.Transaction, nil
	}

	return nil, err
}

// GetTransactionStatus will return the transaction's status
func (tp *TransactionProcessor) GetTransactionStatus(txHash string) (string, error) {
	observersResponses := make(map[uint32][]string)
	observers := tp.proc.GetAllObservers()
	for _, observer := range observers {
		txStatusResponse := &data.ResponseTxStatus{}
		err := tp.proc.CallGetRestEndPoint(observer.Address, TransactionPath+txHash+"/status", txStatusResponse)
		if err != nil {
			continue
		}

		observersResponses[observer.ShardId] = append(observersResponses[observer.ShardId], txStatusResponse.Status)
	}

	return parseTxStatusResponses(observersResponses)
}

func parseTxStatusResponses(allResponses map[uint32][]string) (string, error) {
	okResponses := make(map[uint32][]string)
	for shardID, responses := range allResponses {
		for _, response := range responses {
			if response == UnknownStatusTx {
				continue
			}
			okResponses[shardID] = append(okResponses[shardID], response)
		}
	}

	if len(okResponses) > 1 {
		return "", ErrCannotGetTransactionStatus
	}

	for _, responses := range okResponses {
		firstOkResponse := responses[0]
		for _, response := range responses {
			if firstOkResponse != response {
				return "", ErrCannotGetTransactionStatus
			}
		}
		return firstOkResponse, nil
	}
	return UnknownStatusTx, nil
}

func (tp *TransactionProcessor) groupTxsByShard(txs []*data.Transaction) map[uint32][]*data.Transaction {
	txsMap := make(map[uint32][]*data.Transaction)
	for idx, tx := range txs {
		senderBytes, err := tp.pubKeyConverter.Decode(tx.Sender)
		if err != nil {
			continue
		}

		senderShardID, err := tp.proc.ComputeShardId(senderBytes)
		if err != nil {
			continue
		}

		tx.Index = idx
		txsMap[senderShardID] = append(txsMap[senderShardID], tx)
	}

	return txsMap
}

func (tp *TransactionProcessor) checkTransactionFields(tx *data.Transaction) error {
	_, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return &errors.ErrInvalidTxFields{
			Message: errors.ErrInvalidSenderAddress.Error(),
			Reason:  err.Error(),
		}
	}

	_, err = tp.pubKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return &errors.ErrInvalidTxFields{
			Message: errors.ErrInvalidReceiverAddress.Error(),
			Reason:  err.Error(),
		}
	}

	_, err = hex.DecodeString(tx.Signature)
	if err != nil {
		return &errors.ErrInvalidTxFields{
			Message: errors.ErrInvalidSignatureHex.Error(),
			Reason:  err.Error(),
		}
	}

	return nil
}
