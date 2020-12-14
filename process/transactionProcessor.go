package process

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionPath defines the transaction group path of the node
const TransactionPath = "/transaction/"

// TransactionSendPath defines the single transaction send path of the node
const TransactionSendPath = "/transaction/send"

// TransactionSimulatePath defines single transaction simulate path of the node
const TransactionSimulatePath = "/transaction/simulate"

// MultipleTransactionsPath defines the multiple transactions send path of the node
const MultipleTransactionsPath = "/transaction/send-multiple"

// TransactionCostPath defines the transaction's cost path of the node
const TransactionCostPath = "/transaction/cost"

// UnknownStatusTx defines the response that should be received from an observer when transaction status is unknown
const UnknownStatusTx = "unknown"

const withResultsParam = "?withResults=true"

type requestType int

const (
	requestTypeObservers        requestType = 0
	requestTypeFullHistoryNodes requestType = 1
)

type erdTransaction struct {
	Nonce     uint64 `json:"nonce"`
	Value     string `json:"value"`
	RcvAddr   string `json:"receiver"`
	SndAddr   string `json:"sender"`
	GasPrice  uint64 `json:"gasPrice,omitempty"`
	GasLimit  uint64 `json:"gasLimit,omitempty"`
	Data      []byte `json:"data,omitempty"`
	Signature string `json:"signature,omitempty"`
	ChainID   string `json:"chainID"`
	Version   uint32 `json:"version"`
}

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc            Processor
	pubKeyConverter core.PubkeyConverter
	hasher          hashing.Hasher
	marshalizer     marshal.Marshalizer
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(
	proc Processor,
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
) (*TransactionProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}
	if check.IfNil(hasher) {
		return nil, ErrNilHasher
	}
	if check.IfNil(marshalizer) {
		return nil, ErrNilMarshalizer
	}

	return &TransactionProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
		hasher:          hasher,
		marshalizer:     marshalizer,
	}, nil
}

// SendTransaction relays the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SendTransaction(tx *data.Transaction) (string, int, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	senderBuff, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	shardID, err := tp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	observers, err := tp.proc.GetObservers(shardID)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	for _, observer := range observers {
		txResponse := &data.ResponseTransaction{}

		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, TransactionSendPath, tx, txResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("Transaction sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				shardID,
				txResponse.Data.TxHash,
			))
			return txResponse.Data.TxHash, respCode, nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return "", respCode, err
	}

	return "", http.StatusInternalServerError, ErrSendingRequest
}

// SimulateTransaction relays the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, int, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	senderBuff, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	senderShardID, err := tp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	observers, err := tp.proc.GetObservers(senderShardID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	response, status, err := tp.simulateTransaction(observers, tx)
	if err != nil {
		return nil, status, fmt.Errorf("%w while trying to simulate on sender shard (shard %d)", err, senderShardID)
	}

	receiverBuff, err := tp.pubKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	receiverShardID, err := tp.proc.ComputeShardId(receiverBuff)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if senderShardID == receiverShardID {
		return &data.GenericAPIResponse{
			Data:  response.Data,
			Error: response.Error,
			Code:  response.Code,
		}, http.StatusOK, nil
	}

	observersForReceiverShard, err := tp.proc.GetObservers(receiverShardID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	responseFromReceiverShard, status, err := tp.simulateTransaction(observersForReceiverShard, tx)
	if err != nil {
		return nil, status, fmt.Errorf("%w while trying to simulate on receiver shard (shard %d)", err, receiverShardID)
	}

	simulationResult := data.ResponseTransactionSimulationCrossShard{}
	simulationResult.Data.Result = map[string]data.TransactionSimulationResults{
		"senderShard":   response.Data.Result,
		"receiverShard": responseFromReceiverShard.Data.Result,
	}

	return &data.GenericAPIResponse{
		Data:  simulationResult.Data,
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}, http.StatusOK, nil
}

func (tp *TransactionProcessor) simulateTransaction(observers []*data.NodeData, tx *data.Transaction) (*data.ResponseTransactionSimulation, int, error) {
	gRespCode := http.StatusInternalServerError
	for _, observer := range observers {
		txResponse := &data.ResponseTransactionSimulation{}

		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, TransactionSimulatePath, tx, txResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("Transaction simulation sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				observer.ShardId,
				txResponse.Data.Result.Hash,
			))
			return txResponse, respCode, nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			gRespCode = respCode
			continue
		}

		// if the request was bad, return the error message
		return nil, respCode, err
	}

	return nil, gRespCode, ErrSendingRequest
}

// SendMultipleTransactions relays the post request by sending the request to the first available observer and replies back the answer
func (tp *TransactionProcessor) SendMultipleTransactions(txs []*data.Transaction) (
	data.MultipleTransactionsResponseData, int, error,
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
		return data.MultipleTransactionsResponseData{}, http.StatusBadRequest, ErrNoValidTransactionToSend
	}

	txsHashes := make(map[int]string)
	txsByShardID := tp.groupTxsByShard(txsToSend)
	for shardID, groupOfTxs := range txsByShardID {
		observersInShard, err := tp.proc.GetObservers(shardID)
		if err != nil {
			return data.MultipleTransactionsResponseData{}, http.StatusInternalServerError, ErrMissingObserver
		}

		for _, observer := range observersInShard {
			txResponse := &data.ResponseMultipleTransactions{}
			respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, MultipleTransactionsPath, groupOfTxs, txResponse)
			if respCode == http.StatusOK && err == nil {
				log.Info("transactions sent",
					"observer", observer.Address,
					"shard ID", shardID,
					"total processed", txResponse.Data.NumOfTxs,
				)
				totalTxsSent += txResponse.Data.NumOfTxs

				for key, hash := range txResponse.Data.TxsHashes {
					txsHashes[groupOfTxs[key].Index] = hash
				}

				break
			}

			log.LogIfError(err)
		}
	}

	return data.MultipleTransactionsResponseData{
		NumOfTxs:  totalTxsSent,
		TxsHashes: txsHashes,
	}, http.StatusOK, nil
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (tp *TransactionProcessor) TransactionCostRequest(tx *data.Transaction) (string, int, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	observers, err := tp.proc.GetAllObservers()
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	gRespCode := http.StatusInternalServerError
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
			return strconv.Itoa(int(txCostResponse.Data.TxCost)), http.StatusOK, nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			gRespCode = respCode
			continue
		}

		// if the request was bad, return the error message
		return "", respCode, err

	}

	return "", gRespCode, ErrSendingRequest
}

// GetTransaction should return a transaction from observer
func (tp *TransactionProcessor) GetTransaction(txHash string, withResults bool) (*data.FullTransaction, int, error) {
	tx, status, err := tp.getTxFromObservers(txHash, requestTypeFullHistoryNodes, withResults)
	if err != nil {
		return nil, status, err
	}

	tx.HyperblockNonce = tx.NotarizedAtDestinationInMetaNonce
	tx.HyperblockHash = tx.NotarizedAtDestinationInMetaHash
	return tx, status, nil
}

//GetTransactionByHashAndSenderAddress returns a transaction
func (tp *TransactionProcessor) GetTransactionByHashAndSenderAddress(
	txHash string,
	sndAddr string,
	withEvents bool,
) (*data.FullTransaction, int, error) {
	tx, status, err := tp.getTxWithSenderAddr(txHash, sndAddr, withEvents)
	if err != nil {
		return nil, status, err
	}

	return tx, status, nil
}

func (tp *TransactionProcessor) getShardByAddress(address string) (uint32, error) {
	var shardID uint32
	if metachainIDStr := fmt.Sprintf("%d", core.MetachainShardId); address != metachainIDStr {
		senderBuff, err := tp.pubKeyConverter.Decode(address)
		if err != nil {
			return 0, err
		}

		shardID, err = tp.proc.ComputeShardId(senderBuff)
		if err != nil {
			return 0, err
		}
	} else {
		shardID = core.MetachainShardId
	}

	return shardID, nil
}

// GetTransactionStatus returns the status of a transaction
func (tp *TransactionProcessor) GetTransactionStatus(txHash string, sender string) (string, int, error) {
	if sender != "" {
		tx, status, err := tp.getTxWithSenderAddr(txHash, sender, false)
		if err != nil {
			return UnknownStatusTx, status, err
		}

		return string(tx.Status), status, nil
	}

	// get status of transaction from random observers
	tx, status, err := tp.getTxFromObservers(txHash, requestTypeObservers, false)
	if err != nil {
		return UnknownStatusTx, status, errors.ErrTransactionNotFound
	}

	return string(tx.Status), status, nil
}

func (tp *TransactionProcessor) getTxFromObservers(txHash string, reqType requestType, withResults bool) (*data.FullTransaction, int, error) {
	observersShardIDs := tp.proc.GetShardIDs()
	for _, observerShardID := range observersShardIDs {
		nodesInShard, err := tp.getNodesInShard(observerShardID, reqType)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var getTxResponse *data.GetTransactionResponse
		var withHttpError bool
		var ok bool
		for _, observerInShard := range nodesInShard {
			getTxResponse, ok, withHttpError = tp.getTxFromObserver(observerInShard, txHash, withResults)
			if !withHttpError {
				break
			}
		}

		if !ok || getTxResponse == nil {
			continue
		}

		sndShardID, err := tp.getShardByAddress(getTxResponse.Data.Transaction.Sender)
		if err != nil {
			log.Warn("cannot compute shard ID from sender address",
				"sender address", getTxResponse.Data.Transaction.Sender,
				"error", err.Error())
		}

		rcvShardID, err := tp.getShardByAddress(getTxResponse.Data.Transaction.Receiver)
		if err != nil {
			log.Warn("cannot compute shard ID from receiver address",
				"receiver address", getTxResponse.Data.Transaction.Receiver,
				"error", err.Error())
		}

		isIntraShard := sndShardID == rcvShardID
		observerIsInDestShard := rcvShardID == observerShardID
		if isIntraShard {
			return &getTxResponse.Data.Transaction, http.StatusOK, nil
		}

		if observerIsInDestShard {
			// need to get transaction from source shard and merge scResults
			// if withEvents is true
			return tp.alterTxWithScResultsFromSourceIfNeeded(txHash, &getTxResponse.Data.Transaction, withResults), http.StatusOK, nil
		}

		// get transaction from observer that is in destination shard
		txFromDstShard, ok := tp.getTxFromDestShard(txHash, rcvShardID, withResults)
		if ok {
			alteredTxFromDest := mergeScResultsFromSourceAndDestIfNeeded(&getTxResponse.Data.Transaction, txFromDstShard, withResults)
			return alteredTxFromDest, http.StatusOK, nil
		}

		// return transaction from observer from source shard
		//if did not get ok responses from observers from destination shard
		return &getTxResponse.Data.Transaction, http.StatusOK, nil
	}

	return nil, http.StatusNotFound, errors.ErrTransactionNotFound
}

func (tp *TransactionProcessor) alterTxWithScResultsFromSourceIfNeeded(txHash string, tx *data.FullTransaction, withResults bool) *data.FullTransaction {
	if !withResults || len(tx.ScResults) == 0 {
		return tx
	}

	observers, err := tp.getNodesInShard(tx.SourceShard, requestTypeFullHistoryNodes)
	if err != nil {
		return tx
	}

	for _, observer := range observers {
		getTxResponse, ok, _ := tp.getTxFromObserver(observer, txHash, withResults)
		if !ok {
			continue
		}

		alteredTxFromDest := mergeScResultsFromSourceAndDestIfNeeded(&getTxResponse.Data.Transaction, tx, withResults)
		return alteredTxFromDest
	}

	return tx
}

func (tp *TransactionProcessor) getTxWithSenderAddr(txHash, sender string, withEvents bool) (*data.FullTransaction, int, error) {
	sndShardID, err := tp.getShardByAddress(sender)
	if err != nil {
		return nil, http.StatusBadRequest, errors.ErrInvalidSenderAddress
	}

	observers, err := tp.getNodesInShard(sndShardID, requestTypeFullHistoryNodes)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, observer := range observers {
		getTxResponse, ok, _ := tp.getTxFromObserver(observer, txHash, withEvents)
		if !ok {
			continue
		}

		rcvShardID, err := tp.getShardByAddress(getTxResponse.Data.Transaction.Receiver)
		if err != nil {
			log.Warn("cannot compute shard ID from receiver address",
				"receiver address", getTxResponse.Data.Transaction.Receiver,
				"error", err.Error())
		}

		isIntraShard := rcvShardID == sndShardID
		if isIntraShard {
			return &getTxResponse.Data.Transaction, http.StatusOK, nil
		}

		txFromDstShard, ok := tp.getTxFromDestShard(txHash, rcvShardID, withEvents)
		if ok {
			alteredTxFromDest := mergeScResultsFromSourceAndDestIfNeeded(&getTxResponse.Data.Transaction, txFromDstShard, withEvents)
			return alteredTxFromDest, http.StatusOK, nil
		}

		return &getTxResponse.Data.Transaction, http.StatusOK, nil
	}

	return nil, http.StatusNotFound, errors.ErrTransactionNotFound
}

func mergeScResultsFromSourceAndDestIfNeeded(
	sourceTx *data.FullTransaction,
	destTx *data.FullTransaction,
	withEvents bool,
) *data.FullTransaction {
	if !withEvents {
		return destTx
	}

	scResults := append(sourceTx.ScResults, destTx.ScResults...)
	scResultsNew := getScResultsUnion(scResults)

	destTx.ScResults = scResultsNew

	return destTx
}

func getScResultsUnion(scResults []*transaction.ApiSmartContractResult) []*transaction.ApiSmartContractResult {
	scResultsHash := make(map[string]*transaction.ApiSmartContractResult, 0)
	for _, scResult := range scResults {
		scResultsHash[scResult.Hash] = scResult
	}

	newSlice := make([]*transaction.ApiSmartContractResult, 0)
	for _, scResult := range scResultsHash {
		newSlice = append(newSlice, scResult)
	}

	return newSlice
}

func (tp *TransactionProcessor) getTxFromObserver(
	observer *data.NodeData,
	txHash string,
	withResults bool,
) (*data.GetTransactionResponse, bool, bool) {
	getTxResponse := &data.GetTransactionResponse{}
	apiPath := TransactionPath + txHash
	if withResults {
		apiPath += withResultsParam
	}

	respCode, err := tp.proc.CallGetRestEndPoint(observer.Address, apiPath, getTxResponse)
	if err != nil {
		log.Trace("cannot get transaction", "address", observer.Address, "error", err)

		return nil, false, true
	}

	if respCode != http.StatusOK {
		return nil, false, false
	}

	return getTxResponse, true, false
}

func (tp *TransactionProcessor) getTxFromDestShard(txHash string, dstShardID uint32, withEvents bool) (*data.FullTransaction, bool) {
	// cross shard transaction
	destinationShardObservers, err := tp.proc.GetObservers(dstShardID)
	if err != nil {
		return nil, false
	}

	apiPath := TransactionPath + txHash
	if withEvents {
		apiPath += withResultsParam
	}

	for _, dstObserver := range destinationShardObservers {
		getTxResponseDst := &data.GetTransactionResponse{}
		respCode, err := tp.proc.CallGetRestEndPoint(dstObserver.Address, apiPath, getTxResponseDst)
		if err != nil {
			log.Trace("cannot get transaction", "address", dstObserver.Address, "error", err)
			continue
		}

		if respCode != http.StatusOK {
			continue
		}

		return &getTxResponseDst.Data.Transaction, true
	}

	return nil, false
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

	if tx.ChainID == "" {
		return &errors.ErrInvalidTxFields{
			Message: "transaction must contain chainID",
			Reason:  "no chainID",
		}
	}

	if tx.Version == 0 {
		return &errors.ErrInvalidTxFields{
			Message: "transaction must contain version",
			Reason:  "no version",
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

// ComputeTransactionHash will compute the hash of a given transaction
// TODO move to node
func (tp *TransactionProcessor) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	valueBig, ok := big.NewInt(0).SetString(tx.Value, 10)
	if !ok {
		return "", ErrInvalidTransactionValueField
	}
	receiverAddress, err := tp.pubKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return "", ErrInvalidAddress
	}

	senderAddress, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return "", ErrInvalidAddress
	}

	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return "", ErrInvalidSignatureBytes
	}

	regularTx := &transaction.Transaction{
		Nonce:     tx.Nonce,
		Value:     valueBig,
		RcvAddr:   receiverAddress,
		SndAddr:   senderAddress,
		GasPrice:  tx.GasPrice,
		GasLimit:  tx.GasLimit,
		Data:      tx.Data,
		ChainID:   []byte(tx.ChainID),
		Version:   tx.Version,
		Signature: signatureBytes,
	}

	txHash, err := core.CalculateHash(tp.marshalizer, tp.hasher, regularTx)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(txHash), nil
}

func (tp *TransactionProcessor) getNodesInShard(shardID uint32, reqType requestType) ([]*data.NodeData, error) {
	if reqType == requestTypeFullHistoryNodes {
		fullHistoryNodes, err := tp.proc.GetFullHistoryNodes(shardID)
		if err == nil && len(fullHistoryNodes) > 0 {
			return fullHistoryNodes, nil
		}
	}

	observers, err := tp.proc.GetObservers(shardID)

	return observers, err
}
