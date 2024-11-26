package factory

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/marshal"

	"github.com/multiversx/mx-chain-proxy-go/facade"
	"github.com/multiversx/mx-chain-proxy-go/factory"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/logsevents"
	"github.com/multiversx/mx-chain-proxy-go/process/txcost"
)

// CreateTransactionProcessor will return the transaction processor needed for current settings
func CreateTransactionProcessor(
	proc process.Processor,
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
	allowEntireTxPoolFetch bool,
	runTypeComponents factory.RunTypeComponentsHolder,
) (facade.TransactionProcessor, error) {
	newTxCostProcessor := func() (process.TransactionCostHandler, error) {
		return txcost.NewTransactionCostProcessor(
			proc,
			pubKeyConverter,
		)
	}

	logsMerger, err := logsevents.NewLogsMerger(hasher, &marshal.JsonMarshalizer{})
	if err != nil {
		return nil, err
	}

	return process.NewTransactionProcessor(
		proc,
		pubKeyConverter,
		hasher,
		marshalizer,
		newTxCostProcessor,
		logsMerger,
		allowEntireTxPoolFetch,
		runTypeComponents.TxNotarizationCheckerHandlerCreator(),
	)
}
