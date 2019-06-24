package process

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the address path at which the nodes answer
const TransactionPath = "/transaction/send"
const GenerateMultiplePath = "/transaction/generate-and-send-multiple"

const faucetValue = 10

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
func (ap *TransactionProcessor) SendTransaction(nonce uint64, sender string, receiver string, value *big.Int, code string, signature []byte) (string, error) {
	senderBuff, err := hex.DecodeString(sender)
	if err != nil {
		return "", err
	}

	shardId, err := ap.proc.ComputeShardId(senderBuff)
	if err != nil {
		return "", err
	}

	observers, err := ap.proc.GetObservers(shardId)
	if err != nil {
		return "", err
	}

	for _, observer := range observers {
		tx := &data.Transaction{
			Nonce:     nonce,
			Sender:    sender,
			Receiver:  receiver,
			Value:     value,
			Data:      code,
			Signature: hex.EncodeToString(signature),
		}
		txResponse := &data.ResponseTransaction{}

		err = ap.proc.CallPostRestEndPoint(observer.Address, TransactionPath, tx, txResponse)
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

func (ap *TransactionProcessor) SendUserFunds(receiver string) error {
	receiverBuff, err := hex.DecodeString(receiver)
	if err != nil {
		return err
	}

	shardId, err := ap.proc.ComputeShardId(receiverBuff)
	if err != nil {
		return err
	}

	observers, err := ap.proc.GetObservers(shardId)
	if err != nil {
		return err
	}

	fundsBody := &data.FundsRequest{
		Receiver: receiver,
		Value: big.NewInt(faucetValue),
		TxCount: 1,
	}
	fundsResponse := &data.ResponseFunds{}

	for _, observer := range observers {
		err = ap.proc.CallPostRestEndPoint(observer.Address, GenerateMultiplePath, fundsBody, fundsResponse)
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
