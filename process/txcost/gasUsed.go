package txcost

import (
	"runtime/debug"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

func (tcp *transactionCostProcessor) prepareGasUsed(senderShardID, receiverShardID uint32, res *data.TxCostResponseData) {
	extra := 0
	if senderShardID != receiverShardID {
		extra = 1
	}

	tcp.computeResponsesGasUsed(extra, res)
}

func (tcp *transactionCostProcessor) computeResponsesGasUsed(extra int, res *data.TxCostResponseData) {
	numResponses := len(tcp.responses)

	to := numResponses - 1 - extra
	gasUsed := uint64(0)
	for idx := 0; idx < to; idx++ {
		responseIndex := idx + extra
		if numResponses-1 < responseIndex || len(tcp.scrsToExecute)-1 < idx {
			log.Warn("transactionCostProcessor.computeResponsesGasUsed()", "stack", string(debug.Stack()))

			res.RetMessage = "something went wrong"
			res.TxCost = 0
			return
		}

		diff := tcp.responses[idx+extra].Data.TxCost - tcp.scrsToExecute[idx].GasLimit
		gasUsed += diff
	}

	gasForLastResponse := tcp.responses[numResponses-1].Data.TxCost
	gasUsed += gasForLastResponse
	res.TxCost = gasUsed
}
