package process

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the address path at which the nodes answer
const TransactionPath = "/transaction/send"

// MultipleTransactionsPath defines the address path at which the nodes answer
const MultipleTransactionsPath = "/transaction/send-multiple"

// TransactionCostPath define the address path at which the observer node answer
const TransactionCostPath = "/transaction/cost"

type erdTransaction struct {
	Nonce     uint64 `capid:"0" json:"nonce"`
	Value     string `capid:"1" json:"value"`
	RcvAddr   string `capid:"2" json:"receiver"`
	SndAddr   string `capid:"3" json:"sender"`
	GasPrice  uint64 `capid:"4" json:"gasPrice,omitempty"`
	GasLimit  uint64 `capid:"5" json:"gasLimit,omitempty"`
	Data      []byte `capid:"6" json:"data,omitempty"`
	Signature string `capid:"7" json:"signature,omitempty"`
}

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc            Processor
	pubKeyConverter state.PubkeyConverter
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(
	proc Processor,
	pubKeyConverter state.PubkeyConverter,
) (*TransactionProcessor, error) {
	if proc == nil {
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
func (tp *TransactionProcessor) SendTransaction(apiTx *data.ApiTransaction) (int, string, error) {
	err := tp.checkTransactionFields(apiTx)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	tx := convertToInnerStruct(apiTx)
	senderBuff, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return http.StatusBadRequest, "", err
	}

	shardId, err := tp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	observers, err := tp.proc.GetObservers(shardId)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	for _, observer := range observers {
		txResponse := &data.ResponseTransaction{}

		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, TransactionPath, tx, txResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("Transaction sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				shardId,
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
func (tp *TransactionProcessor) SendMultipleTransactions(apiTxs []*data.ApiTransaction) (uint64, error) {
	totalTxsSent := uint64(0)
	txs := make([]*data.Transaction, len(apiTxs))
	for i := 0; i < len(apiTxs); i++ {
		currentTx := apiTxs[i]
		err := tp.checkTransactionFields(currentTx)
		if err != nil {
			log.Warn("invalid tx received",
				"sender", currentTx.Sender,
				"receiver", currentTx.Receiver,
				"error", err)
			continue
		}
		txs[i] = convertToInnerStruct(currentTx)
	}
	if len(txs) == 0 {
		return 0, ErrNoValidTransactionToSend
	}

	txsByShardId := tp.getTxsByShardId(txs)
	for shardId, txsInShard := range txsByShardId {
		observersInShard, err := tp.proc.GetObservers(shardId)
		if err != nil {
			return 0, ErrMissingObserver
		}

		for _, observer := range observersInShard {
			txResponse := &data.ResponseMultiTransactions{}
			respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, MultipleTransactionsPath, txsInShard, txResponse)
			if respCode == http.StatusOK && err == nil {
				log.Info("transactions sent",
					"observer", observer.Address,
					"shard id", shardId,
					"total processed", txResponse.NumOfTxs,
				)
				totalTxsSent += txResponse.NumOfTxs
				break
			}

			log.LogIfError(err)
		}
	}

	return totalTxsSent, nil
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (tp *TransactionProcessor) TransactionCostRequest(tx *data.ApiTransaction) (string, error) {
	observers := tp.proc.GetAllObservers()

	err := tp.checkTransactionFields(tx)
	if err != nil {
		return "", err
	}

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

func (tp *TransactionProcessor) getTxsByShardId(txs []*data.Transaction) map[uint32][]*data.Transaction {
	txsMap := make(map[uint32][]*data.Transaction, 0)
	for _, tx := range txs {
		senderBytes, err := tp.pubKeyConverter.Decode(tx.Sender)
		if err != nil {
			continue
		}

		senderShardId, err := tp.proc.ComputeShardId(senderBytes)
		if err != nil {
			continue
		}

		txsMap[senderShardId] = append(txsMap[senderShardId], tx)
	}

	return txsMap
}

func convertToInnerStruct(tx *data.ApiTransaction) *data.Transaction {
	return &data.Transaction{
		Nonce:     tx.Nonce,
		Value:     tx.Value,
		Receiver:  tx.Receiver,
		Sender:    tx.Sender,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.GasLimit,
		Data:      []byte(tx.Data),
		Signature: tx.Signature,
	}
}

func convertToAPIStruct(tx *data.Transaction) *data.ApiTransaction {
	return &data.ApiTransaction{
		Nonce:     tx.Nonce,
		Value:     tx.Value,
		Receiver:  tx.Receiver,
		Sender:    tx.Sender,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.GasLimit,
		Data:      string(tx.Data),
		Signature: tx.Signature,
	}
}

func (tp *TransactionProcessor) checkTransactionFields(tx *data.ApiTransaction) error {
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
