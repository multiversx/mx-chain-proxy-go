package process

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// TransactionPath defines the transaction group path of the node
const TransactionPath = "/transaction/"

// TransactionsPoolPath defines the transactions pool path of the node
const TransactionsPoolPath = "/transaction/pool"

// TransactionSendPath defines the single transaction send path of the node
const TransactionSendPath = "/transaction/send"

// TransactionSimulatePath defines single transaction simulate path of the node
const TransactionSimulatePath = "/transaction/simulate"

// MultipleTransactionsPath defines the multiple transactions send path of the node
const MultipleTransactionsPath = "/transaction/send-multiple"

// SCRsByTxHash defines smart contract results by transaction hash path of the node
const SCRsByTxHash = "/transaction/scrs-by-tx-hash/"

const (
	withResultsParam                = "?withResults=true"
	scrHashParam                    = "?scrHash=%s"
	checkSignatureFalse             = "?checkSignature=false"
	bySenderParam                   = "&by-sender="
	fieldsParam                     = "?fields="
	lastNonceParam                  = "?last-nonce=true"
	nonceGapsParam                  = "?nonce-gaps=true"
	internalVMErrorsEventIdentifier = "internalVMErrors" // TODO export this in mx-chain-core-go, remove unexported definitions from mx-chain-vm's
	moveBalanceDescriptor           = "MoveBalance"
	relayedV1TransactionDescriptor  = "RelayedTx"
	relayedV2TransactionDescriptor  = "RelayedTxV2"
	relayedV3TransactionDescriptor  = "RelayedTxV3"
	emptyDataStr                    = ""
)

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

type tupleHashWasFetched struct {
	hash    string
	fetched bool
}

// TransactionProcessor is able to process transaction requests
type TransactionProcessor struct {
	proc                         Processor
	pubKeyConverter              core.PubkeyConverter
	hasher                       hashing.Hasher
	marshalizer                  marshal.Marshalizer
	relayedTxsMarshaller         marshal.Marshalizer
	newTxCostProcessor           func() (TransactionCostHandler, error)
	mergeLogsHandler             LogsMergerHandler
	shouldAllowEntireTxPoolFetch bool
}

// NewTransactionProcessor creates a new instance of TransactionProcessor
func NewTransactionProcessor(
	proc Processor,
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
	newTxCostProcessor func() (TransactionCostHandler, error),
	logsMerger LogsMergerHandler,
	allowEntireTxPoolFetch bool,
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
	if newTxCostProcessor == nil {
		return nil, ErrNilNewTxCostHandlerFunc
	}
	if check.IfNil(logsMerger) {
		return nil, ErrNilLogsMerger
	}

	// no reason to get this from configs. If we are going to change the marshaller for the relayed transaction v1,
	// we will need also an enable epoch handler
	relayedTxsMarshaller := &marshal.JsonMarshalizer{}
	return &TransactionProcessor{
		proc:                         proc,
		pubKeyConverter:              pubKeyConverter,
		hasher:                       hasher,
		marshalizer:                  marshalizer,
		newTxCostProcessor:           newTxCostProcessor,
		mergeLogsHandler:             logsMerger,
		shouldAllowEntireTxPoolFetch: allowEntireTxPoolFetch,
		relayedTxsMarshaller:         relayedTxsMarshaller,
	}, nil
}

// SendTransaction relays the post request by sending the request to the right observer and replies back the answer
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

	observers, err := tp.proc.GetObservers(shardID, data.AvailabilityRecent)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	txResponse := data.ResponseTransaction{}
	for _, observer := range observers {

		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, TransactionSendPath, tx, &txResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("Transaction sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				shardID,
				txResponse.Data.TxHash,
			))
			return respCode, txResponse.Data.TxHash, nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return respCode, "", err
	}

	return http.StatusInternalServerError, "", WrapObserversError(txResponse.Error)
}

// SimulateTransaction relays the post request by sending the request to the right observer and replies back the answer
func (tp *TransactionProcessor) SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return nil, err
	}

	senderBuff, err := tp.pubKeyConverter.Decode(tx.Sender)
	if err != nil {
		return nil, err
	}

	senderShardID, err := tp.proc.ComputeShardId(senderBuff)
	if err != nil {
		return nil, err
	}

	observers, err := tp.proc.GetObservers(senderShardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	response, err := tp.simulateTransaction(observers, tx, checkSignature)
	if err != nil {
		return nil, fmt.Errorf("%w while trying to simulate on sender shard (shard %d)", err, senderShardID)
	}

	receiverBuff, err := tp.pubKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return nil, err
	}

	receiverShardID, err := tp.proc.ComputeShardId(receiverBuff)
	if err != nil {
		return nil, err
	}

	if senderShardID == receiverShardID {
		return &data.GenericAPIResponse{
			Data:  response.Data,
			Error: response.Error,
			Code:  response.Code,
		}, nil
	}

	observersForReceiverShard, err := tp.proc.GetObservers(receiverShardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseFromReceiverShard, err := tp.simulateTransaction(observersForReceiverShard, tx, checkSignature)
	if err != nil {
		return nil, fmt.Errorf("%w while trying to simulate on receiver shard (shard %d)", err, receiverShardID)
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
	}, nil
}

func (tp *TransactionProcessor) simulateTransaction(
	observers []*data.NodeData,
	tx *data.Transaction,
	checkSignature bool,
) (*data.ResponseTransactionSimulation, error) {
	txSimulatePath := TransactionSimulatePath
	if !checkSignature {
		txSimulatePath += checkSignatureFalse
	}

	txResponse := data.ResponseTransactionSimulation{}
	for _, observer := range observers {

		respCode, err := tp.proc.CallPostRestEndPoint(observer.Address, txSimulatePath, tx, &txResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("Transaction simulation sent successfully to observer %v from shard %v, received tx hash %s",
				observer.Address,
				observer.ShardId,
				txResponse.Data.Result.Hash,
			))
			return &txResponse, nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return nil, err
	}

	return nil, WrapObserversError(txResponse.Error)
}

// SendMultipleTransactions relays the post request by sending the request to the first available observer and replies back the answer
func (tp *TransactionProcessor) SendMultipleTransactions(txs []*data.Transaction) (
	data.MultipleTransactionsResponseData, error,
) {
	// TODO: Analyze and improve the robustness of this function. Currently, an error within `GetObservers`
	// breaks the function and returns nothing (but an error) even if some transactions were actually sent, successfully.

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
		return data.MultipleTransactionsResponseData{}, ErrNoValidTransactionToSend
	}

	txsHashes := make(map[int]string)
	txsByShardID := tp.groupTxsByShard(txsToSend)
	for shardID, groupOfTxs := range txsByShardID {
		observersInShard, err := tp.proc.GetObservers(shardID, data.AvailabilityRecent)
		if err != nil {
			return data.MultipleTransactionsResponseData{}, ErrMissingObserver
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
	}, nil
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (tp *TransactionProcessor) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	err := tp.checkTransactionFields(tx)
	if err != nil {
		return nil, err
	}

	newTxCostProcessor, err := tp.newTxCostProcessor()
	if err != nil {
		return nil, err
	}

	return newTxCostProcessor.ResolveCostRequest(tx)
}

// GetTransaction should return a transaction from observer
func (tp *TransactionProcessor) GetTransaction(txHash string, withResults bool, relayedTxHash string) (*transaction.ApiTransactionResult, error) {
	txHashToGetFromObservers := txHash

	// if the relayed tx hash was provided, this one should be requested from observers
	innerTxRequested := len(relayedTxHash) == len(txHash)
	if innerTxRequested {
		txHashToGetFromObservers = relayedTxHash
	}

	tx, err := tp.getTxFromObservers(txHashToGetFromObservers, requestTypeFullHistoryNodes, withResults)
	if err != nil {
		return nil, err
	}

	tx.HyperblockNonce = tx.NotarizedAtDestinationInMetaNonce
	tx.HyperblockHash = tx.NotarizedAtDestinationInMetaHash

	if len(tx.InnerTransactions) == 0 {
		if innerTxRequested {
			return nil, fmt.Errorf("%w, requested hash %s with relayedTxHash %s, but the relayedTxHash has no inner transaction",
				ErrInvalidHash, txHash, relayedTxHash)
		}

		return tx, nil
	}

	convertRelayedTxV3ToNetworkTx(tx)

	// requested a relayed transaction, returning it after scrs were moved
	if !innerTxRequested {
		return tx, nil
	}

	for _, innerTx := range tx.InnerTransactions {
		if innerTx.Hash == txHash {
			return innerTx, nil
		}
	}

	return nil, fmt.Errorf("%w, but the relayedTx %s has no inner transaction with hash %s",
		ErrInvalidHash, relayedTxHash, txHash)
}

func convertRelayedTxV3ToNetworkTx(tx *transaction.ApiTransactionResult) {
	movedSCRs := make(map[string]struct{}, 0)
	for _, innerTx := range tx.InnerTransactions {
		for _, scr := range tx.SmartContractResults {
			if isResultOfInnerTx(tx.SmartContractResults, scr, innerTx.Hash) {
				innerTx.SmartContractResults = append(innerTx.SmartContractResults, scr)
				movedSCRs[scr.Hash] = struct{}{}
			}
		}
	}

	if len(movedSCRs) == len(tx.SmartContractResults) {
		// all scrs were generated by inner txs
		tx.SmartContractResults = make([]*transaction.ApiSmartContractResult, 0)
		return
	}

	numSCRsLeftForRelayer := len(tx.SmartContractResults) - len(movedSCRs)
	scrsForRelayer := make([]*transaction.ApiSmartContractResult, 0, numSCRsLeftForRelayer)
	for _, scr := range tx.SmartContractResults {
		_, wasMoved := movedSCRs[scr.Hash]
		if !wasMoved {
			scrsForRelayer = append(scrsForRelayer, scr)
		}
	}

	tx.SmartContractResults = scrsForRelayer
}

func isResultOfInnerTx(allScrs []*transaction.ApiSmartContractResult, currentScr *transaction.ApiSmartContractResult, innerTxHash string) bool {
	if currentScr.PrevTxHash == innerTxHash {
		return true
	}

	parentScr := findSCRByHash(allScrs, currentScr.PrevTxHash)
	if check.IfNilReflect(parentScr) {
		return false
	}

	return isResultOfInnerTx(allScrs, parentScr, innerTxHash)
}

func findSCRByHash(allScrs []*transaction.ApiSmartContractResult, hash string) *transaction.ApiSmartContractResult {
	for _, scr := range allScrs {
		if scr.Hash == hash {
			return scr
		}
	}

	return nil
}

// GetTransactionByHashAndSenderAddress returns a transaction
func (tp *TransactionProcessor) GetTransactionByHashAndSenderAddress(
	txHash string,
	sndAddr string,
	withResults bool,
) (*transaction.ApiTransactionResult, int, error) {
	tx, err := tp.getTxWithSenderAddr(txHash, sndAddr, withResults)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return tx, http.StatusOK, nil
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
func (tp *TransactionProcessor) GetTransactionStatus(txHash string, sender string) (string, error) {
	tx, err := tp.getTransaction(txHash, sender, false)
	if err != nil {
		return string(data.TxStatusUnknown), err
	}

	return string(tx.Status), nil
}

func (tp *TransactionProcessor) getTransaction(txHash string, sender string, withResults bool) (*transaction.ApiTransactionResult, error) {
	if sender != "" {
		return tp.getTxWithSenderAddr(txHash, sender, withResults)
	}

	// get status of transaction from random observers
	return tp.getTxFromObservers(txHash, requestTypeObservers, withResults)
}

// GetProcessedTransactionStatus returns the status of a transaction after local processing
func (tp *TransactionProcessor) GetProcessedTransactionStatus(txHash string) (*data.ProcessStatusResponse, error) {
	const withResults = true
	tx, err := tp.getTxFromObservers(txHash, requestTypeObservers, withResults)
	if err != nil {
		return &data.ProcessStatusResponse{
			Status: string(data.TxStatusUnknown),
		}, err
	}

	return tp.computeTransactionStatus(tx, withResults), nil
}

func (tp *TransactionProcessor) computeTransactionStatus(tx *transaction.ApiTransactionResult, withResults bool) *data.ProcessStatusResponse {
	if !withResults {
		return &data.ProcessStatusResponse{
			Status: string(data.TxStatusUnknown),
		}
	}

	if tx.Status == transaction.TxStatusInvalid {
		return &data.ProcessStatusResponse{
			Status: string(transaction.TxStatusFail),
		}
	}
	if tx.Status != transaction.TxStatusSuccess {
		return &data.ProcessStatusResponse{
			Status: string(tx.Status),
		}
	}

	if checkIfMoveBalanceNotarized(tx) {
		return &data.ProcessStatusResponse{
			Status: string(tx.Status),
		}
	}

	allLogs, allScrs, err := tp.gatherAllLogsAndScrs(tx)
	if err != nil {
		log.Warn("error in TransactionProcessor.computeTransactionStatus", "error", err)
		return &data.ProcessStatusResponse{
			Status: string(data.TxStatusUnknown),
		}
	}

	if hasPendingSCR(allScrs) {
		return &data.ProcessStatusResponse{
			Status: string(transaction.TxStatusPending),
		}
	}

	txLogsOnFirstLevel := []*transaction.ApiLogs{tx.Logs}
	failed, reason := checkIfFailed(txLogsOnFirstLevel)
	if failed {
		return &data.ProcessStatusResponse{
			Status: string(transaction.TxStatusFail),
			Reason: reason,
		}
	}

	allLogs, err = tp.addMissingLogsOnProcessingExceptions(tx, allLogs, allScrs)
	if err != nil {
		log.Warn("error in TransactionProcessor.computeTransactionStatus on addMissingLogsOnProcessingExceptions call", "error", err)
		return &data.ProcessStatusResponse{
			Status: string(data.TxStatusUnknown),
		}
	}

	isRelayedV3, status := checkIfRelayedV3Completed(tx)
	if isRelayedV3 {
		return &data.ProcessStatusResponse{
			Status: status,
		}
	}

	failed, reason = checkIfFailed(allLogs)
	if failed {
		return &data.ProcessStatusResponse{
			Status: string(transaction.TxStatusFail),
			Reason: reason,
		}
	}

	if checkIfFailedOnReturnMessage(allScrs, tx) {
		return &data.ProcessStatusResponse{
			Status: string(transaction.TxStatusFail),
		}
	}

	isUnsigned := string(transaction.TxTypeUnsigned) == tx.Type
	if checkIfCompleted(allLogs) || isUnsigned {
		return &data.ProcessStatusResponse{
			Status: string(transaction.TxStatusSuccess),
		}
	}

	return &data.ProcessStatusResponse{
		Status: string(transaction.TxStatusPending),
	}
}

func hasPendingSCR(scrs []*transaction.ApiTransactionResult) bool {
	for _, scr := range scrs {
		if scr.Status == transaction.TxStatusPending {
			return true
		}
	}

	return false
}

func checkIfFailedOnReturnMessage(allScrs []*transaction.ApiTransactionResult, tx *transaction.ApiTransactionResult) bool {
	hasReturnMessageWithZeroValue := len(tx.ReturnMessage) > 0 && isZeroValue(tx.Value)
	if hasReturnMessageWithZeroValue && !isRefundScr(tx.ReturnMessage) {
		return true
	}

	for _, scr := range allScrs {
		if isRefundScr(scr.ReturnMessage) {
			continue
		}

		if len(scr.ReturnMessage) > 0 && isZeroValue(scr.Value) {
			return true
		}
	}

	return false
}

func isRefundScr(returnMessage string) bool {
	return returnMessage == core.GasRefundForRelayerMessage
}

func isZeroValue(value string) bool {
	if len(value) == 0 {
		return true
	}
	return value == "0"
}

func checkIfFailed(logs []*transaction.ApiLogs) (bool, string) {
	found, reason := findIdentifierInLogs(logs, internalVMErrorsEventIdentifier)
	if found {
		return true, reason
	}

	found, reason = findIdentifierInLogs(logs, core.SignalErrorOperation)
	if found {
		return true, reason
	}

	return false, emptyDataStr
}

func checkIfCompleted(logs []*transaction.ApiLogs) bool {
	found, _ := findIdentifierInLogs(logs, core.CompletedTxEventIdentifier)
	if found {
		return true
	}

	found, _ = findIdentifierInLogs(logs, core.SCDeployIdentifier)
	return found
}

func checkIfRelayedV3Completed(tx *transaction.ApiTransactionResult) (bool, string) {
	if len(tx.InnerTransactions) == 0 {
		return false, string(transaction.TxStatusPending)
	}

	return true, string(transaction.TxStatusSuccess)
}

func checkIfMoveBalanceNotarized(tx *transaction.ApiTransactionResult) bool {
	isNotarized := tx.NotarizedAtSourceInMetaNonce > 0 && tx.NotarizedAtDestinationInMetaNonce > 0
	if !isNotarized {
		return false
	}
	isMoveBalance := tx.ProcessingTypeOnSource == moveBalanceDescriptor && tx.ProcessingTypeOnDestination == moveBalanceDescriptor

	return isMoveBalance
}

func (tp *TransactionProcessor) addMissingLogsOnProcessingExceptions(
	tx *transaction.ApiTransactionResult,
	allLogs []*transaction.ApiLogs,
	allScrs []*transaction.ApiTransactionResult,
) ([]*transaction.ApiLogs, error) {
	newLogs, err := tp.handleIntraShardRelayedMoveBalanceTransactions(tx, allScrs)
	if err != nil {
		return nil, err
	}

	allLogs = append(allLogs, newLogs...)

	return allLogs, nil
}

func (tp *TransactionProcessor) handleIntraShardRelayedMoveBalanceTransactions(
	tx *transaction.ApiTransactionResult,
	allScrs []*transaction.ApiTransactionResult,
) ([]*transaction.ApiLogs, error) {
	var newLogs []*transaction.ApiLogs
	isRelayedMoveBalanceTransaction, err := tp.isRelayedMoveBalanceTransaction(tx, allScrs)
	if err != nil {
		return newLogs, err
	}

	if isRelayedMoveBalanceTransaction {
		newLogs = append(newLogs, &transaction.ApiLogs{
			Address: tx.Sender,
			Events: []*transaction.Events{
				{
					Address:    tx.Sender,
					Identifier: core.CompletedTxEventIdentifier,
				},
			},
		})
	}

	return newLogs, nil
}

func (tp *TransactionProcessor) isRelayedMoveBalanceTransaction(
	tx *transaction.ApiTransactionResult,
	allScrs []*transaction.ApiTransactionResult,
) (bool, error) {
	isNotarized := tx.NotarizedAtSourceInMetaNonce > 0 && tx.NotarizedAtDestinationInMetaNonce > 0
	if !isNotarized {
		return false, nil
	}

	isRelayedV1 := tx.ProcessingTypeOnSource == relayedV1TransactionDescriptor && tx.ProcessingTypeOnDestination == relayedV1TransactionDescriptor
	isRelayedV2 := tx.ProcessingTypeOnSource == relayedV2TransactionDescriptor && tx.ProcessingTypeOnDestination == relayedV2TransactionDescriptor

	isRelayedTransaction := isRelayedV1 || isRelayedV2
	if !isRelayedTransaction {
		return false, nil
	}

	if len(allScrs) == 0 {
		return false, fmt.Errorf("no smart contracts results for the given relayed transaction v1")
	}

	firstScr := allScrs[0]
	innerIsMoveBalance := firstScr.ProcessingTypeOnSource == moveBalanceDescriptor && firstScr.ProcessingTypeOnDestination == moveBalanceDescriptor

	return innerIsMoveBalance, nil
}

func findIdentifierInLogs(logs []*transaction.ApiLogs, identifier string) (bool, string) {
	if len(logs) == 0 {
		return false, emptyDataStr
	}

	for _, logInstance := range logs {
		if logInstance == nil {
			continue
		}

		found, reason := findIdentifierInSingleLog(logInstance, identifier)
		if found {
			return true, string(reason)
		}
	}

	return false, emptyDataStr
}

func findIdentifierInSingleLog(log *transaction.ApiLogs, identifier string) (bool, []byte) {
	for _, event := range log.Events {
		if event.Identifier == identifier {
			return true, event.Data
		}
	}

	return false, []byte(emptyDataStr)
}

func (tp *TransactionProcessor) gatherAllLogsAndScrs(tx *transaction.ApiTransactionResult) ([]*transaction.ApiLogs, []*transaction.ApiTransactionResult, error) {
	const withResults = true
	allLogs := make([]*transaction.ApiLogs, 0)
	allScrs := make([]*transaction.ApiTransactionResult, 0)

	if tx.Logs != nil {
		allLogs = append(allLogs, tx.Logs)
	}

	for _, scrFromTx := range tx.SmartContractResults {
		scr, err := tp.GetTransaction(scrFromTx.Hash, withResults, "")
		if err != nil {
			return nil, nil, fmt.Errorf("%w for scr hash %s", err, scrFromTx.Hash)
		}

		if scr == nil {
			continue
		}
		allScrs = append(allScrs, scr)

		if scr.Logs == nil {
			continue
		}
		allLogs = append(allLogs, scr.Logs)
	}

	return allLogs, allScrs, nil
}

func (tp *TransactionProcessor) getTxFromObservers(txHash string, reqType requestType, withResults bool) (*transaction.ApiTransactionResult, error) {
	observersShardIDs := tp.proc.GetShardIDs()
	shardIDWasFetch := make(map[uint32]*tupleHashWasFetched)
	for _, observerShardID := range observersShardIDs {
		nodesInShard, err := tp.getNodesInShard(observerShardID, reqType)
		if err != nil {
			return nil, err
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
		shardIDWasFetch[sndShardID] = &tupleHashWasFetched{
			hash:    getTxResponse.Data.Transaction.Hash,
			fetched: false,
		}

		rcvShardID, err := tp.getShardByAddress(getTxResponse.Data.Transaction.Receiver)
		if err != nil {
			log.Warn("cannot compute shard ID from receiver address",
				"receiver address", getTxResponse.Data.Transaction.Receiver,
				"error", err.Error())
		}
		shardIDWasFetch[rcvShardID] = &tupleHashWasFetched{
			hash:    getTxResponse.Data.Transaction.Hash,
			fetched: false,
		}

		isIntraShard := sndShardID == rcvShardID
		observerIsInDestShard := rcvShardID == observerShardID

		if isIntraShard {
			shardIDWasFetch[sndShardID].fetched = true
			if len(getTxResponse.Data.Transaction.SmartContractResults) == 0 {
				return &getTxResponse.Data.Transaction, nil
			}

			tp.extraShardFromSCRs(getTxResponse.Data.Transaction.SmartContractResults, shardIDWasFetch)
		}

		if observerIsInDestShard {
			// need to get transaction from source shard and merge scResults
			// if withEvents is true
			txFromSource := tp.alterTxWithScResultsFromSourceIfNeeded(txHash, &getTxResponse.Data.Transaction, withResults, shardIDWasFetch)

			tp.extraShardFromSCRs(txFromSource.SmartContractResults, shardIDWasFetch)

			err = tp.fetchSCRSBasedOnShardMap(txFromSource, shardIDWasFetch)
			if err != nil {
				return nil, err
			}

			return txFromSource, nil
		}

		// get transaction from observer that is in destination shard
		txFromDstShard, ok := tp.getTxFromDestShard(txHash, rcvShardID, withResults)
		if ok {
			tp.extraShardFromSCRs(txFromDstShard.SmartContractResults, shardIDWasFetch)

			alteredTxFromDest := tp.mergeScResultsFromSourceAndDestIfNeeded(&getTxResponse.Data.Transaction, txFromDstShard, withResults)

			err = tp.fetchSCRSBasedOnShardMap(alteredTxFromDest, shardIDWasFetch)
			if err != nil {
				return nil, err
			}

			return alteredTxFromDest, nil
		}

		// return transaction from observer from source shard
		// if did not get ok responses from observers from destination shard

		err = tp.fetchSCRSBasedOnShardMap(&getTxResponse.Data.Transaction, shardIDWasFetch)
		if err != nil {
			return nil, err
		}

		return &getTxResponse.Data.Transaction, nil
	}

	return nil, errors.ErrTransactionNotFound
}

func (tp *TransactionProcessor) fetchSCRSBasedOnShardMap(tx *transaction.ApiTransactionResult, shardIDWasFetch map[uint32]*tupleHashWasFetched) error {
	for shardID, info := range shardIDWasFetch {
		scrs, err := tp.fetchSCRs(tx.Hash, info.hash, shardID)
		if err != nil {
			return err
		}

		scResults := append(tx.SmartContractResults, scrs...)
		scResultsNew := tp.getScResultsUnion(scResults)

		tx.SmartContractResults = scResultsNew
		info.fetched = true
	}

	return nil
}

func (tp *TransactionProcessor) fetchSCRs(txHash, scrHash string, shardID uint32) ([]*transaction.ApiSmartContractResult, error) {
	observers, err := tp.getNodesInShard(shardID, requestTypeFullHistoryNodes)
	if err != nil {
		return nil, err
	}

	apiPath := SCRsByTxHash + txHash + fmt.Sprintf(scrHashParam, scrHash)
	for _, observer := range observers {
		getTxResponseDst := &data.GetSCRsResponse{}
		respCode, errG := tp.proc.CallGetRestEndPoint(observer.Address, apiPath, getTxResponseDst)
		if errG != nil {
			log.Trace("cannot get smart contract results", "address", observer.Address, "error", errG)
			continue
		}

		if respCode != http.StatusOK {
			continue
		}

		return getTxResponseDst.Data.SCRs, nil
	}

	return []*transaction.ApiSmartContractResult{}, nil

}

func (tp *TransactionProcessor) extraShardFromSCRs(scrs []*transaction.ApiSmartContractResult, shardIDWasFetch map[uint32]*tupleHashWasFetched) {
	for _, scr := range scrs {
		sndShardID, err := tp.getShardByAddress(scr.SndAddr)
		if err != nil {
			log.Warn("cannot compute shard ID from sender address",
				"sender address", scr.SndAddr,
				"error", err.Error())
			continue
		}

		_, found := shardIDWasFetch[sndShardID]
		if !found {
			shardIDWasFetch[sndShardID] = &tupleHashWasFetched{
				hash:    scr.Hash,
				fetched: false,
			}
		}

		rcvShardID, err := tp.getShardByAddress(scr.RcvAddr)
		if err != nil {
			log.Warn("cannot compute shard ID from receiver address",
				"receiver address", scr.RcvAddr,
				"error", err.Error())
			continue
		}

		_, found = shardIDWasFetch[rcvShardID]
		if !found {
			shardIDWasFetch[rcvShardID] = &tupleHashWasFetched{
				hash:    scr.Hash,
				fetched: false,
			}
		}
	}
}

func (tp *TransactionProcessor) alterTxWithScResultsFromSourceIfNeeded(txHash string, tx *transaction.ApiTransactionResult, withResults bool, shardIDWasFetch map[uint32]*tupleHashWasFetched) *transaction.ApiTransactionResult {
	if !withResults || len(tx.SmartContractResults) == 0 {
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

		alteredTxFromDest := tp.mergeScResultsFromSourceAndDestIfNeeded(&getTxResponse.Data.Transaction, tx, withResults)

		shardIDWasFetch[tx.SourceShard] = &tupleHashWasFetched{
			hash:    getTxResponse.Data.Transaction.Hash,
			fetched: true,
		}

		return alteredTxFromDest
	}

	return tx
}

func (tp *TransactionProcessor) getTxWithSenderAddr(txHash, sender string, withResults bool) (*transaction.ApiTransactionResult, error) {
	observers, sndShardID, err := tp.getShardObserversForSender(sender, requestTypeFullHistoryNodes)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		getTxResponse, ok, _ := tp.getTxFromObserver(observer, txHash, withResults)
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
			return &getTxResponse.Data.Transaction, nil
		}

		txFromDstShard, ok := tp.getTxFromDestShard(txHash, rcvShardID, withResults)
		if ok {
			alteredTxFromDest := tp.mergeScResultsFromSourceAndDestIfNeeded(&getTxResponse.Data.Transaction, txFromDstShard, withResults)
			return alteredTxFromDest, nil
		}

		return &getTxResponse.Data.Transaction, nil
	}

	return nil, errors.ErrTransactionNotFound
}

func (tp *TransactionProcessor) mergeScResultsFromSourceAndDestIfNeeded(
	sourceTx *transaction.ApiTransactionResult,
	destTx *transaction.ApiTransactionResult,
	withEvents bool,
) *transaction.ApiTransactionResult {
	if !withEvents {
		return destTx
	}

	scResults := append(sourceTx.SmartContractResults, destTx.SmartContractResults...)
	scResultsNew := tp.getScResultsUnion(scResults)

	destTx.SmartContractResults = scResultsNew

	return destTx
}

func (tp *TransactionProcessor) getScResultsUnion(scResults []*transaction.ApiSmartContractResult) []*transaction.ApiSmartContractResult {
	scResultsHash := make(map[string]*transaction.ApiSmartContractResult)
	for _, scResult := range scResults {
		scResultFromMap, found := scResultsHash[scResult.Hash]
		if !found {
			scResultsHash[scResult.Hash] = scResult
			continue
		}

		mergedLog := tp.mergeLogsHandler.MergeLogEvents(scResultFromMap.Logs, scResult.Logs)
		scResultsHash[scResult.Hash] = scResult
		scResultsHash[scResult.Hash].Logs = mergedLog
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

		if respCode == http.StatusTooManyRequests {
			log.Warn("too many requests while getting tx from observer", "address", observer.Address, "tx hash", txHash)
		}

		return nil, false, true
	}

	if respCode != http.StatusOK {
		return nil, false, false
	}

	return getTxResponse, true, false
}

func (tp *TransactionProcessor) getTxFromDestShard(txHash string, dstShardID uint32, withEvents bool) (*transaction.ApiTransactionResult, bool) {
	// cross shard transaction
	destinationShardObservers, err := tp.proc.GetObservers(dstShardID, data.AvailabilityAll)
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

	if len(tx.GuardianSignature) > 0 {
		_, err = hex.DecodeString(tx.GuardianSignature)
		if err != nil {
			return &errors.ErrInvalidTxFields{
				Message: errors.ErrInvalidGuardianSignatureHex.Error(),
				Reason:  err.Error(),
			}
		}
	}
	if len(tx.GuardianAddr) > 0 {
		_, err = tp.pubKeyConverter.Decode(tx.GuardianAddr)
		if err != nil {
			return &errors.ErrInvalidTxFields{
				Message: errors.ErrInvalidGuardianAddress.Error(),
				Reason:  err.Error(),
			}
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

	if len(tx.GuardianAddr) > 0 {
		regularTx.GuardianAddr, err = tp.pubKeyConverter.Decode(tx.GuardianAddr)
		if err != nil {
			return "", errors.ErrInvalidGuardianAddress
		}
	}

	if len(tx.GuardianSignature) > 0 {
		regularTx.GuardianSignature, err = hex.DecodeString(tx.GuardianSignature)
		if err != nil {
			return "", errors.ErrInvalidGuardianSignatureHex
		}
	}

	txHash, err := core.CalculateHash(tp.marshalizer, tp.hasher, regularTx)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(txHash), nil
}

func (tp *TransactionProcessor) getNodesInShard(shardID uint32, reqType requestType) ([]*data.NodeData, error) {
	if reqType == requestTypeFullHistoryNodes {
		fullHistoryNodes, err := tp.proc.GetFullHistoryNodes(shardID, data.AvailabilityAll)
		if err == nil && len(fullHistoryNodes) > 0 {
			return fullHistoryNodes, nil
		}
	}

	observers, err := tp.proc.GetObservers(shardID, data.AvailabilityAll)

	return observers, err
}

// GetTransactionsPool should return all transactions from all shards pool
func (tp *TransactionProcessor) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	if !tp.shouldAllowEntireTxPoolFetch {
		return nil, errors.ErrOperationNotAllowed
	}

	txPool, err := tp.getTxPool(fields)
	if err != nil {
		return nil, err
	}

	return txPool, nil
}

// GetTransactionsPoolForShard should return transactions pool from one observer from shard
func (tp *TransactionProcessor) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	if !tp.shouldAllowEntireTxPoolFetch {
		return nil, errors.ErrOperationNotAllowed
	}

	txPool, err := tp.getTxPoolForShard(shardID, fields)
	if err != nil {
		return nil, err
	}

	return txPool, nil
}

// GetTransactionsPoolForSender should return transactions for sender from observer's pool
func (tp *TransactionProcessor) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	txPool, err := tp.getTxPoolForSender(sender, fields)
	if err != nil {
		return nil, err
	}

	return txPool, nil
}

// GetLastPoolNonceForSender should return last nonce for sender from observer's pool
func (tp *TransactionProcessor) GetLastPoolNonceForSender(sender string) (uint64, error) {
	return tp.getLastTxPoolNonceForSender(sender)
}

// GetTransactionsPoolNonceGapsForSender should return nonce gaps for sender from observer's pool
func (tp *TransactionProcessor) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	return tp.getTxPoolNonceGapsForSender(sender)
}

func (tp *TransactionProcessor) getShardObserversForSender(sender string, observersType requestType) ([]*data.NodeData, uint32, error) {
	sndShardID, err := tp.getShardByAddress(sender)
	if err != nil {
		return nil, 0, errors.ErrInvalidSenderAddress
	}

	observers, err := tp.getNodesInShard(sndShardID, observersType)
	if err != nil {
		return nil, 0, err
	}

	return observers, sndShardID, nil
}

func (tp *TransactionProcessor) getTxPool(fields string) (*data.TransactionsPool, error) {
	shardIDs := tp.proc.GetShardIDs()
	txs := &data.TransactionsPool{
		RegularTransactions:  make([]data.WrappedTransaction, 0),
		SmartContractResults: make([]data.WrappedTransaction, 0),
		Rewards:              make([]data.WrappedTransaction, 0),
	}
	for _, shard := range shardIDs {
		intraShardTxs, err := tp.getTxPoolForShard(shard, fields)
		if err != nil {
			continue
		}

		txs.RegularTransactions = append(txs.RegularTransactions, intraShardTxs.RegularTransactions...)
		txs.Rewards = append(txs.Rewards, intraShardTxs.Rewards...)
		txs.SmartContractResults = append(txs.SmartContractResults, intraShardTxs.SmartContractResults...)
	}

	return txs, nil
}

func (tp *TransactionProcessor) getTxPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	observers, err := tp.getNodesInShard(shardID, requestTypeObservers)
	if err != nil {
		log.Trace("cannot get observers for shard", "shard", shardID, "error", err)
		return nil, err
	}

	for _, observer := range observers {
		txs, ok := tp.getTxPoolFromObserver(observer, fields)
		if !ok {
			continue
		}

		return txs, nil
	}

	log.Trace("cannot get tx pool for shard", "shard", shardID, "error", errors.ErrTransactionsNotFoundInPool.Error())
	return nil, errors.ErrTransactionsNotFoundInPool
}

func (tp *TransactionProcessor) getTxPoolFromObserver(
	observer *data.NodeData,
	fields string,
) (*data.TransactionsPool, bool) {
	txsPoolResponse := &data.TransactionsPoolApiResponse{}
	apiPath := TransactionsPoolPath + fieldsParam + fields

	respCode, err := tp.proc.CallGetRestEndPoint(observer.Address, apiPath, txsPoolResponse)
	if err != nil {
		log.Trace("cannot get tx pool", "address", observer.Address, "error", err)

		if respCode == http.StatusTooManyRequests {
			log.Warn("too many requests while getting tx pool", "address", observer.Address)
		}

		return nil, false
	}

	if respCode != http.StatusOK {
		return nil, false
	}

	return &txsPoolResponse.Data.Transactions, true
}

func (tp *TransactionProcessor) getTxPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	observers, _, err := tp.getShardObserversForSender(sender, requestTypeObservers)
	if err != nil {
		return nil, err
	}

	txsInPool := &data.TransactionsPoolForSender{
		Transactions: []data.WrappedTransaction{},
	}
	var ok bool
	for _, observer := range observers {
		txsInPool, ok = tp.getTxPoolForSenderFromObserver(observer, sender, fields)
		if ok {
			break
		}
	}

	return txsInPool, nil
}

func (tp *TransactionProcessor) getTxPoolForSenderFromObserver(
	observer *data.NodeData,
	sender string,
	fields string,
) (*data.TransactionsPoolForSender, bool) {
	txsPoolResponse := &data.TransactionsPoolForSenderApiResponse{}
	apiPath := TransactionsPoolPath + fieldsParam + fields + bySenderParam + sender

	respCode, err := tp.proc.CallGetRestEndPoint(observer.Address, apiPath, txsPoolResponse)
	if err != nil {
		log.Trace("cannot get tx pool for sender", "address", observer.Address, "sender", sender, "error", err)

		if respCode == http.StatusTooManyRequests {
			log.Warn("too many requests while getting tx pool for sender", "address", observer.Address, "sender", sender)
		}

		return nil, false
	}

	if respCode != http.StatusOK {
		return nil, false
	}

	return &txsPoolResponse.Data.TxPool, true
}

func (tp *TransactionProcessor) getLastTxPoolNonceForSender(sender string) (uint64, error) {
	observers, _, err := tp.getShardObserversForSender(sender, requestTypeObservers)
	if err != nil {
		return 0, err
	}

	for _, observer := range observers {
		nonce, ok := tp.getLastTxPoolNonceFromObserver(observer, sender)
		if !ok {
			continue
		}

		return nonce, nil
	}

	return 0, errors.ErrTransactionsNotFoundInPool
}

func (tp *TransactionProcessor) getLastTxPoolNonceFromObserver(
	observer *data.NodeData,
	sender string,
) (uint64, bool) {
	lastNonceResponse := &data.TransactionsPoolLastNonceForSenderApiResponse{}
	apiPath := TransactionsPoolPath + lastNonceParam + bySenderParam + sender

	respCode, err := tp.proc.CallGetRestEndPoint(observer.Address, apiPath, lastNonceResponse)
	if err != nil {
		log.Trace("cannot get last nonce from tx pool", "address", observer.Address, "sender", sender, "error", err)

		if respCode == http.StatusTooManyRequests {
			log.Warn("too many requests while getting last nonce from tx pool", "address", observer.Address, "sender", sender)
		}

		return 0, false
	}

	if respCode != http.StatusOK {
		return 0, false
	}

	return lastNonceResponse.Data.Nonce, true
}

func (tp *TransactionProcessor) getTxPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	observers, _, err := tp.getShardObserversForSender(sender, requestTypeObservers)
	if err != nil {
		return nil, err
	}

	nonceGaps := &data.TransactionsPoolNonceGaps{
		Gaps: []data.NonceGap{},
	}
	var ok bool
	for _, observer := range observers {
		nonceGaps, ok = tp.getTxPoolNonceGapsFromObserver(observer, sender)
		if ok {
			break
		}
	}

	return nonceGaps, nil
}

func (tp *TransactionProcessor) getTxPoolNonceGapsFromObserver(
	observer *data.NodeData,
	sender string,
) (*data.TransactionsPoolNonceGaps, bool) {
	nonceGapsResponse := &data.TransactionsPoolNonceGapsForSenderApiResponse{}
	apiPath := TransactionsPoolPath + nonceGapsParam + bySenderParam + sender

	respCode, err := tp.proc.CallGetRestEndPoint(observer.Address, apiPath, nonceGapsResponse)
	if err != nil {
		log.Warn("cannot get nonce gaps from tx pool", "address", observer.Address, "sender", sender, "error", err)

		if respCode == http.StatusTooManyRequests {
			log.Warn("too many requests while getting nonce gaps from tx pool", "address", observer.Address, "sender", sender)
		}

		return nil, false
	}

	if respCode != http.StatusOK {
		return nil, false
	}

	return &nonceGapsResponse.Data.NonceGaps, true
}
