package txcost

import (
	"runtime/debug"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
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

	defer func() {
		r := recover()
		if r != nil {
			log.Warn("transactionCostProcessor.computeResponsesGasUsed()", "stack", string(debug.Stack()))

			res.RetMessage = "something went wrong"
			res.TxCost = 0
		}
	}()

	gasUsed := uint64(0)
	to := len(tcp.responses) - 1 - extra
	for idx := 0; idx < to; idx++ {
		gasUsed += tcp.responses[idx+extra].Data.TxCost - tcp.txsFromSCR[idx].GasLimit
	}

	gasUsed += tcp.responses[numResponses-1].Data.TxCost
	res.TxCost = gasUsed
}
