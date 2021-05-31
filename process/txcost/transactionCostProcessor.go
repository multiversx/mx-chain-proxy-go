package txcost

import (
	"math"
	"net/http"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

// TransactionCostPath defines the transaction's cost path of the node
const TransactionCostPath = "/transaction/cost"

var log = logger.GetOrCreate("process/txcost")

type transactionCostProcessor struct {
	proc            process.Processor
	pubKeyConverter core.PubkeyConverter
	responses       []*data.ResponseTxCost
}

// NewTransactionCostProcessor will create a new instance of the transactionCostProcessor
func NewTransactionCostProcessor(proc process.Processor, pubKeyConverter core.PubkeyConverter) (*transactionCostProcessor, error) {
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
	}, nil
}

// RezolveCostRequest will resolve the transaction cost request
func (tcp *transactionCostProcessor) RezolveCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	initialGasLimit := tx.GasLimit
	if tx.GasLimit == 0 {
		// TODO modify here if the max gas limit for simulate will be max gas limit per block
		initialGasLimit = math.MaxUint64
	}

	res, err := tcp.doCostRequests(tx)
	if err != nil {
		return nil, err
	}

	if len(tcp.responses) < 2 {
		return res, nil
	}

	for _, currentRes := range tcp.responses {
		if currentRes.Data.RetMessage == "" {
			continue
		}

		res.RetMessage = currentRes.Data.RetMessage
		res.TxCost = 0
		return res, nil
	}

	numRes := len(tcp.responses)
	totalGas := tcp.responses[numRes-1].Data.TxCost + initialGasLimit - tcp.responses[numRes-2].Data.TxCost
	res.TxCost = totalGas

	return res, nil
}

func (tcp *transactionCostProcessor) doCostRequests(tx *data.Transaction) (*data.TxCostResponseData, error) {
	senderShardID, receiverShardID, err := tcp.computeSenderAndReceiverShardID(tx.Sender, tx.Receiver)
	if err != nil {
		return nil, err
	}

	observers, err := tcp.proc.GetObservers(receiverShardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		txCostResponse := &data.ResponseTxCost{}
		respCode, errCall := tcp.proc.CallPostRestEndPoint(observer.Address, TransactionCostPath, tx, txCostResponse)
		if respCode == http.StatusOK && errCall == nil {
			return tcp.processResponse(senderShardID, receiverShardID, txCostResponse, tx)
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return nil, err

	}

	return nil, ErrSendingRequest
}

func (tcp *transactionCostProcessor) processResponse(
	senderShardID uint32,
	receiverShardID uint32,
	response *data.ResponseTxCost,
	originalTx *data.Transaction,
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
		res, err := tcp.processScResult(senderShardID, receiverShardID, scr, originalTx)
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
	scr *data.ApiSmartContractResultExtended,
	originalTx *data.Transaction,
) (*data.TxCostResponseData, error) {
	scrSenderShardID, scrReceiverShardID, err := tcp.computeSenderAndReceiverShardID(scr.SndAddr, scr.RcvAddr)
	if err != nil {
		return nil, err
	}

	// TODO check if this condition is enough
	shouldIgnoreSCR := receiverShardID == scrReceiverShardID
	shouldIgnoreSCR = shouldIgnoreSCR || (scrReceiverShardID == senderShardID && scr.CallType == 0)
	shouldIgnoreSCR = shouldIgnoreSCR || scrSenderShardID == core.MetachainShardId
	if shouldIgnoreSCR {
		return nil, nil
	}

	txFromScr := convertSCRInTransaction(scr, originalTx)

	res, errReq := tcp.doCostRequests(txFromScr)
	if errReq != nil {
		return nil, errReq
	}

	return res, errReq
}
