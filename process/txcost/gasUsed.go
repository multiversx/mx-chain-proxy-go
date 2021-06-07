package txcost

import (
	"runtime/debug"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func (tcp *transactionCostProcessor) prepareGasUsed(senderShardID, receiverShardID uint32, res *data.TxCostResponseData) {
	numResponses := len(tcp.responses)
	extra := 0
	if senderShardID != receiverShardID {
		extra = 1
	}

	syncChan := make(chan struct{})
	go func(c chan struct{}) {
		defer func() {
			if r := recover(); r != nil {
				log.Warn("transactionCostProcessor.prepareGasUsed()", "stack", string(debug.Stack()))

				res.RetMessage = "something went wrong"
				res.TxCost = 0
			}

			c <- struct{}{}
		}()

		gasUsed := uint64(0)
		to := len(tcp.responses) - 1 - extra
		for idx := 0; idx < to; idx++ {
			gasUsed += tcp.responses[idx+extra].Data.TxCost - tcp.txsFromSCR[idx].GasLimit
		}

		gasUsed += tcp.responses[numResponses-1].Data.TxCost
		res.TxCost = gasUsed
	}(syncChan)

	<-syncChan
}
