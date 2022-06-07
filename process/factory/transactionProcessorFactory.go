package factory

import (
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/logsevents"
	"github.com/ElrondNetwork/elrond-proxy-go/process/txcost"
)

// CreateTransactionProcessor will return the transaction processor needed for current settings
func CreateTransactionProcessor(
	proc process.Processor,
	pubKeyConverter core.PubkeyConverter,
	hasher hashing.Hasher,
	marshalizer marshal.Marshalizer,
	maxGasLimitPerBlockShardStr string,
	maxGasLimitPerBlockMetaStr string,
) (facade.TransactionProcessor, error) {
	newTxCostProcessor := func() (process.TransactionCostHandler, error) {
		return txcost.NewTransactionCostProcessor(
			proc,
			pubKeyConverter,
			maxGasLimitPerBlockShardStr,
			maxGasLimitPerBlockMetaStr,
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
	)
}
