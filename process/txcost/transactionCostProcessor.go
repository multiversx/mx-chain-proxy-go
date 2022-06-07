package txcost

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

// TransactionCostPath defines the transaction's cost path of the node
const TransactionCostPath = "/transaction/cost"

var log = logger.GetOrCreate("process/txcost")

type transactionCostProcessor struct {
	proc                     process.Processor
	pubKeyConverter          core.PubkeyConverter
	responses                []*data.ResponseTxCost
	maxGasLimitPerBlockShard uint64
	maxGasLimitPerBlockMeta  uint64
	txsFromSCR               []*data.Transaction
}

// NewTransactionCostProcessor will create a new instance of the transactionCostProcessor
func NewTransactionCostProcessor(
	proc process.Processor,
	pubKeyConverter core.PubkeyConverter,
	maxGasLimitPerBlockShardStr string,
	maxGasLimitPerBlockMetaStr string,
) (*transactionCostProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	maxGasLimitPerBlockShard, err := strconv.ParseUint(maxGasLimitPerBlockShardStr, 10, 64)
	if err != nil {
		return nil, err
	}
	maxGasLimitPerBlockMeta, err := strconv.ParseUint(maxGasLimitPerBlockMetaStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &transactionCostProcessor{
		proc:                     proc,
		pubKeyConverter:          pubKeyConverter,
		maxGasLimitPerBlockShard: maxGasLimitPerBlockShard,
		maxGasLimitPerBlockMeta:  maxGasLimitPerBlockMeta,
		responses:                make([]*data.ResponseTxCost, 0),
		txsFromSCR:               make([]*data.Transaction, 0),
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

	tcp.prepareGasUsed(senderShardID, receiverShardID, res)

	return res, nil
}

func (tcp *transactionCostProcessor) doCostRequests(senderShardID, receiverShardID uint32, tx *data.Transaction) (*data.TxCostResponseData, error) {
	shouldExecuteOnSource := senderShardID != receiverShardID && len(tcp.responses) == 0
	if shouldExecuteOnSource {
		observers, errGet := tcp.proc.GetObservers(senderShardID)
		if errGet != nil {
			return nil, errGet
		}

		res, errExe := tcp.executeRequest(senderShardID, receiverShardID, observers, tx)
		if errExe != nil {
			return nil, errExe
		}

		if res.RetMessage != "" {
			return res, nil
		}
	}

	observers, err := tcp.proc.GetObservers(receiverShardID)
	if err != nil {
		return nil, err
	}

	return tcp.executeRequest(senderShardID, receiverShardID, observers, tx)
}

func (tcp *transactionCostProcessor) executeRequest(
	senderShardID uint32,
	receiverShardID uint32,
	observers []*data.NodeData,
	tx *data.Transaction,
) (*data.TxCostResponseData, error) {
	for _, observer := range observers {
		txCostResponse := &data.ResponseTxCost{}
		respCode, errCall := tcp.proc.CallPostRestEndPoint(observer.Address, TransactionCostPath, tx, txCostResponse)
		if respCode == http.StatusOK && errCall == nil {
			return tcp.processResponse(senderShardID, receiverShardID, txCostResponse, tx)
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(errCall)
			continue
		}

		// if the request was bad, return the error message
		return nil, errCall

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
	scr *data.ExtendedApiSmartContractResult,
	originalTx *data.Transaction,
) (*data.TxCostResponseData, error) {
	scrSenderShardID, scrReceiverShardID, err := tcp.computeSenderAndReceiverShardID(scr.SndAddr, scr.RcvAddr)
	if err != nil {
		return nil, err
	}

	// TODO check if this condition is enough
	shouldIgnoreSCR := receiverShardID == scrReceiverShardID
	shouldIgnoreSCR = shouldIgnoreSCR || (scrReceiverShardID == senderShardID && scr.CallType == vm.DirectCall)
	shouldIgnoreSCR = shouldIgnoreSCR || scrSenderShardID == core.MetachainShardId
	if shouldIgnoreSCR {
		return nil, nil
	}

	txFromScr := convertSCRInTransaction(scr, originalTx)
	tcp.txsFromSCR = append(tcp.txsFromSCR, txFromScr)

	observers, err := tcp.proc.GetObservers(scrReceiverShardID)
	if err != nil {
		return nil, err
	}

	res, err := tcp.executeRequest(scrSenderShardID, scrReceiverShardID, observers, txFromScr)
	if err != nil {
		return nil, err
	}

	return res, nil
}
