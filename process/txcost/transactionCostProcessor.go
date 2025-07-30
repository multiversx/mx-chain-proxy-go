package txcost

import (
	"bytes"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
)

const (
	// TransactionCostPath defines the transaction's cost path of the node
	TransactionCostPath = "/transaction/cost"
	SCRCostPath         = "/transaction/cost-scr"

	tooMuchGasProvidedMessage = "@too much gas provided"
)

var log = logger.GetOrCreate("process/txcost")

type transactionCostProcessor struct {
	proc            process.Processor
	pubKeyConverter core.PubkeyConverter
	responses       []*data.ResponseTxCost
	scrsToExecute   []*smartContractResult.SmartContractResult
	hasExecutedSCR  bool
}

// NewTransactionCostProcessor will create a new instance of the transactionCostProcessor
func NewTransactionCostProcessor(
	proc process.Processor,
	pubKeyConverter core.PubkeyConverter,
) (*transactionCostProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &transactionCostProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
		responses:       make([]*data.ResponseTxCost, 0),
		scrsToExecute:   make([]*smartContractResult.SmartContractResult, 0),
	}, nil
}

// ResolveCostRequest will resolve the transaction cost request
func (tcp *transactionCostProcessor) ResolveCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	senderShardID, receiverShardID, err := tcp.computeSenderAndReceiverShardID(tx.Sender, tx.Receiver)
	if err != nil {
		return nil, err
	}

	res, err := tcp.doCostRequests(senderShardID, receiverShardID, tx)
	if err != nil {
		return nil, err
	}

	shouldReturn := len(tcp.responses) == 1 || (len(tcp.responses) == 2 && senderShardID != receiverShardID)
	if shouldReturn {
		return tcp.extractCorrectResponse(tx.Sender, res), nil
	}

	for _, currentRes := range tcp.responses {
		hasUnsupportedOperations := doEventsContainTopic(&currentRes.Data, tooMuchGasProvidedMessage) || hasSCRWithRefundForSender(tx.Sender, &currentRes.Data)
		shouldReturn = hasUnsupportedOperations && !tcp.hasExecutedSCR
		if shouldReturn {
			return &currentRes.Data, nil
		}

		if currentRes.Data.RetMessage == "" {
			continue
		}

		res.RetMessage = currentRes.Data.RetMessage
		res.TxCost = 0
		return res, nil
	}

	tcp.prepareGasUsed(senderShardID, receiverShardID, res)

	return res, nil
}

func (tcp *transactionCostProcessor) doCostRequests(senderShardID, receiverShardID uint32, tx *data.Transaction) (*data.TxCostResponseData, error) {
	shouldExecuteOnSource := senderShardID != receiverShardID && len(tcp.responses) == 0
	if shouldExecuteOnSource {
		observers, errGet := tcp.proc.GetObservers(senderShardID, data.AvailabilityRecent)
		if errGet != nil {
			return nil, errGet
		}

		res, errExe := tcp.executeRequest(senderShardID, receiverShardID, observers, tx, TransactionCostPath)
		if errExe != nil {
			return nil, errExe
		}

		if res.RetMessage != "" {
			return res, nil
		}
	}

	observers, err := tcp.proc.GetObservers(receiverShardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	return tcp.executeRequest(senderShardID, receiverShardID, observers, tx, TransactionCostPath)
}

func (tcp *transactionCostProcessor) executeRequest(
	senderShardID uint32,
	receiverShardID uint32,
	observers []*data.NodeData,
	scrOrTx interface{},
	endpoint string,
) (*data.TxCostResponseData, error) {
	txCostResponse := data.ResponseTxCost{}
	for _, observer := range observers {
		respCode, errCall := tcp.proc.CallPostRestEndPoint(observer.Address, endpoint, scrOrTx, &txCostResponse)
		if respCode == http.StatusOK && errCall == nil {
			return tcp.processResponse(senderShardID, receiverShardID, &txCostResponse)
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(errCall)
			continue
		}

		// if the request was bad, return the error message
		return nil, errCall

	}

	return nil, process.WrapObserversError(txCostResponse.Error)
}

func (tcp *transactionCostProcessor) processResponse(
	senderShardID uint32,
	receiverShardID uint32,
	response *data.ResponseTxCost,
) (*data.TxCostResponseData, error) {
	tcp.responses = append(tcp.responses, response)

	if len(response.Data.ScResults) == 0 || response.Data.RetMessage != "" {
		return &response.Data, nil
	}

	for scrHash, scr := range response.Data.ScResults {
		if scr.Used {
			continue
		}

		scr.Used = true
		res, err := tcp.processScResult(senderShardID, receiverShardID, scr)
		if err != nil {
			log.Warn("cannot process smart contract result", "hash", scrHash, "error", err)
			continue
		}

		if res == nil {
			continue
		}

		mergeResponses(response, res)
	}

	return &response.Data, nil
}

func mergeResponses(finalRes *data.ResponseTxCost, currentRes *data.TxCostResponseData) {
	for scrHash, scr := range currentRes.ScResults {
		finalRes.Data.ScResults[scrHash] = scr
		finalRes.Data.ScResults[scrHash].Used = true
	}
}

func (tcp *transactionCostProcessor) processScResult(
	senderShardID uint32,
	receiverShardID uint32,
	scr *data.ExtendedApiSmartContractResult,
) (*data.TxCostResponseData, error) {
	scrSenderShardID, scrReceiverShardID, err := tcp.computeSenderAndReceiverShardID(scr.SndAddr, scr.RcvAddr)
	if err != nil {
		return nil, err
	}

	ignoreSCRWithESDTTransferNoSCCall := scr.Function == "" && len(scr.Tokens) > 0

	shouldIgnoreSCR := receiverShardID == scrReceiverShardID
	shouldIgnoreSCR = shouldIgnoreSCR || (scrReceiverShardID == senderShardID && scr.CallType == vm.DirectCall)
	shouldIgnoreSCR = shouldIgnoreSCR || scrSenderShardID == core.MetachainShardId
	shouldIgnoreSCR = shouldIgnoreSCR || ignoreSCRWithESDTTransferNoSCCall
	if shouldIgnoreSCR {
		return nil, nil
	}

	protocolSCR, err := tcp.convertExtendedSCRInProtocolSCR(scr)
	if err != nil {
		return nil, err
	}

	tcp.scrsToExecute = append(tcp.scrsToExecute, protocolSCR)

	observers, err := tcp.proc.GetObservers(scrReceiverShardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	res, err := tcp.executeRequest(scrSenderShardID, scrReceiverShardID, observers, protocolSCR, SCRCostPath)
	if err != nil {
		return nil, err
	}

	tcp.hasExecutedSCR = true

	return res, nil
}

func (tcp *transactionCostProcessor) extractCorrectResponse(txSender string, currentRes *data.TxCostResponseData) *data.TxCostResponseData {
	if len(tcp.responses) == 1 {
		return currentRes
	}

	for _, res := range tcp.responses {
		if doEventsContainTopic(&res.Data, tooMuchGasProvidedMessage) || hasSCRWithRefundForSender(txSender, &res.Data) {
			return &res.Data
		}
	}

	return currentRes
}

func doEventsContainTopic(res *data.TxCostResponseData, providedTopic string) bool {
	if res.Logs == nil {
		return false
	}

	for _, event := range res.Logs.Events {
		if event.Identifier != core.WriteLogIdentifier {
			continue
		}

		for _, topic := range event.Topics {
			if bytes.Contains(topic, []byte(providedTopic)) {
				return true
			}
		}
	}

	return false
}

func hasSCRWithRefundForSender(txSender string, res *data.TxCostResponseData) bool {
	for _, scr := range res.ScResults {
		if scr.IsRefund && scr.RcvAddr == txSender {
			return true
		}
	}

	return false
}
