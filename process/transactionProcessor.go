package process

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the address path at which the nodes answer
const TransactionPath = "/transaction/send"

// MultipleTransactionsPath defines the address path at which the nodes answer
const MultipleTransactionsPath = "/transaction/send-multiple"

// GenerateMultiplePath defines the path for generating transactions
const GenerateMultiplePath = "/transaction/generate-and-send-multiple"

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc Processor
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(proc Processor) (*TransactionProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &TransactionProcessor{
		proc: proc,
	}, nil
}

// SendTransaction relay the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SendTransaction(tx *data.Transaction) (string, error) {

	senderBuff, err := hex.DecodeString(tx.Sender)
	if err != nil {
		return "", err
	}

	shardId, err := tp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return "", err
	}

	observers, err := tp.proc.GetObservers(shardId)
	if err != nil {
		return "", err
	}

	for _, observer := range observers {
		txResponse := &data.ResponseTransaction{}

		err = tp.proc.CallPostRestEndPoint(observer.Address, TransactionPath, tx, txResponse)
		if err == nil {
			log.Info(fmt.Sprintf("Transaction sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				shardId,
				txResponse.TxHash,
			))
			return txResponse.TxHash, nil
		}

		log.LogIfError(err)
	}

	return "", ErrSendingRequest
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
			err = tp.proc.CallPostRestEndPoint(observer.Address, MultipleTransactionsPath, txsInShard, txResponse)
			if err == nil {
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

// SendUserFunds transmits a request to the right observer to load a provided address with some predefined balance
func (tp *TransactionProcessor) SendUserFunds(receiver string, value *big.Int) error {
	receiverBuff, err := hex.DecodeString(receiver)
	if err != nil {
		return err
	}

	shardId, err := tp.proc.ComputeShardId(receiverBuff)
	if err != nil {
		return err
	}

	observers, err := tp.proc.GetObservers(shardId)
	if err != nil {
		return err
	}

	fundsBody := &data.FundsRequest{
		Receiver: receiver,
		Value:    value,
		TxCount:  1,
	}
	fundsResponse := &data.ResponseFunds{}

	for _, observer := range observers {
		err = tp.proc.CallPostRestEndPoint(observer.Address, GenerateMultiplePath, fundsBody, fundsResponse)
		if err == nil {
			log.Info(fmt.Sprintf("Funds sent successfully from observer %v from shard %v, to address %s",
				observer.Address,
				shardId,
				receiver,
			))
			return nil
		}

		log.LogIfError(err)
	}

	return ErrSendingRequest
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
