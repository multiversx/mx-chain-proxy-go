package services

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// SupportedOperationTypes is a list of the supported operations.
var SupportedOperationTypes = []string{
	opTransfer, opFee, opReward, opScResult,
}

// ElrondCurrency is the currency used on the Elrond blockchain.
var ElrondCurrency = &types.Currency{
	Symbol:   "eGLD",
	Decimals: 18,
}

type objectsMap map[string]interface{}

func estimateGasLimit(operationType string, networkConfig *client.NetworkConfig, options objectsMap) (uint64, *types.Error) {
	gasForDataField := uint64(0)
	if dataFieldI, ok := options["data"]; ok {
		dataField := fmt.Sprintf("%v", dataFieldI)
		gasForDataField = networkConfig.GasPerDataByte * uint64(len(dataField))
	}

	switch operationType {
	case opTransfer:
		return networkConfig.MinGasLimit + gasForDataField, nil
	default:
		return 0, ErrNotImplemented
	}
}

func checkProvidedGasLimit(providedGasLimit uint64, txType string, options objectsMap, networkConfig *client.NetworkConfig) *types.Error {
	estimatedGasLimit, err := estimateGasLimit(txType, networkConfig, options)
	if err != nil {
		return err
	}

	if providedGasLimit < estimatedGasLimit {
		return ErrInsufficientGasLimit
	}

	return nil
}
