package process

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the address path at which the nodes answer
const TransactionPath = "/transaction/send"

// MultipleTransactionsPath defines the address path at which the nodes answer
const MultipleTransactionsPath = "/transaction/send-multiple"

type erdTransaction struct {
	Nonce     uint64   `capid:"0" json:"nonce"`
	Value     *big.Int `capid:"1" json:"value"`
	RcvAddr   []byte   `capid:"2" json:"receiver"`
	SndAddr   []byte   `capid:"3" json:"sender"`
	GasPrice  uint64   `capid:"4" json:"gasPrice,omitempty"`
	GasLimit  uint64   `capid:"5" json:"gasLimit,omitempty"`
	Data      string   `capid:"6" json:"data,omitempty"`
	Signature []byte   `capid:"7" json:"signature,omitempty"`
	Challenge []byte   `capid:"8" json:"challenge,omitempty"`
}

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc    Processor
	keyGen  crypto.KeyGenerator
	signer  crypto.SingleSigner
	skSlice [][]byte
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(
	proc Processor,
	keyGen crypto.KeyGenerator,
	signer crypto.SingleSigner,
	skSlice [][]byte,
) (*TransactionProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}
	if keyGen == nil {
		return nil, ErrNilKeyGen
	}
	if signer == nil {
		return nil, ErrNilSingleSigner
	}
	if skSlice == nil {
		return nil, errors.New("nil private keys slice")
	}

	return &TransactionProcessor{
		proc:    proc,
		keyGen:  keyGen,
		signer:  signer,
		skSlice: skSlice,
	}, nil
}

// SignAndSendTransaction relay the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SignAndSendTransaction(tx *data.Transaction, sk []byte) (string, error) {
	tx, err := tp.getSignedTx(tx, sk)
	if err != nil {
		return "", err
	}

	return tp.SendTransaction(tx)
}

func (tp *TransactionProcessor) getSignedTx(tx *data.Transaction, sk []byte) (*data.Transaction, error) {
	marshalizedTxBeforeSigning := tp.marshalTxForSigning(tx)
	privKey, err := tp.keyGen.PrivateKeyFromByteArray(sk)
	if err != nil {
		return nil, err
	}

	signature, err := tp.signer.Sign(privKey, marshalizedTxBeforeSigning)
	if err != nil {
		return nil, err
	}

	signHex := hex.EncodeToString(signature)
	tx.Signature = signHex

	return tx, nil
}

func (tp *TransactionProcessor) marshalTxForSigning(tx *data.Transaction) []byte {
	snrB, _ := hex.DecodeString(tx.Sender)
	rcB, _ := hex.DecodeString(tx.Receiver)
	erdTx := erdTransaction{
		Nonce:    tx.Nonce,
		Value:    tx.Value,
		RcvAddr:  rcB,
		SndAddr:  snrB,
		GasPrice: tx.GasPrice,
		GasLimit: tx.GasLimit,
		Data:     tx.Data,
	}

	mtx, _ := json.Marshal(erdTx)
	return mtx
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
	//receiverBuff, err := hex.DecodeString(receiver)
	//if err != nil {
	//	return err
	//}
	//
	//shardId, err := tp.proc.ComputeShardId(receiverBuff)
	//if err != nil {
	//	return err
	//}
	//
	//observers, err := tp.proc.GetObservers(shardId)
	//if err != nil {
	//	return err
	//}
	//
	//fundsBody := &data.FundsRequest{
	//	Receiver: receiver,
	//	Value:    value,
	//	TxCount:  1,
	//}
	//fundsResponse := &data.ResponseFunds{}

	privKeyBytes := tp.skSlice[0]
	privKey, err := tp.keyGen.PrivateKeyFromByteArray(privKeyBytes)
	senderPubKeyHex, err := tp.hexPubKeyFromPrivKey(privKey)
	if err != nil {
		return err
	}

	genTx := data.Transaction{
		Nonce:     0,
		Value:     value,
		Receiver:  receiver,
		Sender:    senderPubKeyHex,
		GasPrice:  1,
		GasLimit:  5,
		Data:      "",
		Signature: "",
		Challenge: "",
	}

	_, err = tp.SignAndSendTransaction(&genTx, privKeyBytes)

	return err
}

func (tp *TransactionProcessor) hexPubKeyFromPrivKey(sk crypto.PrivateKey) (string, error) {
	pk := sk.GeneratePublic()
	pkBytes, err := pk.ToByteArray()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(pkBytes), nil
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
