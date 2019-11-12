package process

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the address path at which the nodes answer
const TransactionPath = "/transaction/send"

// MultipleTransactionsPath defines the address path at which the nodes answer
const MultipleTransactionsPath = "/transaction/send-multiple"

type erdTransaction struct {
	Nonce     uint64 `capid:"0" json:"nonce"`
	Value     string `capid:"1" json:"value"`
	RcvAddr   []byte `capid:"2" json:"receiver"`
	SndAddr   []byte `capid:"3" json:"sender"`
	GasPrice  uint64 `capid:"4" json:"gasPrice,omitempty"`
	GasLimit  uint64 `capid:"5" json:"gasLimit,omitempty"`
	Data      string `capid:"6" json:"data,omitempty"`
	Signature []byte `capid:"7" json:"signature,omitempty"`
	Challenge []byte `capid:"8" json:"challenge,omitempty"`
}

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc Processor
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(
	proc Processor,
) (*TransactionProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &TransactionProcessor{
		proc: proc,
	}, nil
}

// SendTransaction relay the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SendTransaction(tx *data.Transaction) (int, string, error) {

	senderBuff, err := hex.DecodeString(tx.Sender)
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
func (tp *TransactionProcessor) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	totalTxsSent := uint64(0)
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
				log.Info(fmt.Sprintf("Transactions sent successfully to observer %v from shard %v, total processed: %d",
					observer.Address,
					observer.ShardId,
					txResponse.NumOfTxs,
				))
				totalTxsSent += txResponse.NumOfTxs
				continue
			}

			log.LogIfError(err)
		}
	}

	return totalTxsSent, nil
}

func (tp *TransactionProcessor) getTxsByShardId(txs []*data.Transaction) map[uint32][]*data.Transaction {
	txsMap := make(map[uint32][]*data.Transaction, 0)
	for _, tx := range txs {
		senderBytes, err := hex.DecodeString(tx.Sender)
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
