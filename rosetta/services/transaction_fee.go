package services

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func computeSuggestedFeeAndGas(txType string, options objectsMap, networkConfig *provider.NetworkConfig) (*big.Int, uint64, uint64, *types.Error) {
	var gasLimit, gasPrice uint64

	if gasLimitI, ok := options["gasLimit"]; ok {
		gasLimit = getUint64Value(gasLimitI)

		err := checkProvidedGasLimit(gasLimit, txType, options, networkConfig)
		if err != nil {
			return nil, 0, 0, err
		}

	} else {
		// if gas limit is not provided, we estimate it
		estimatedGasLimit, err := estimateGasLimit(txType, networkConfig, options)
		if err != nil {
			return nil, 0, 0, err
		}

		gasLimit = estimatedGasLimit
	}

	if gasPriceI, ok := options["gasPrice"]; ok {
		gasPrice = getUint64Value(gasPriceI)

		if gasPrice < networkConfig.MinGasPrice {
			return nil, 0, 0, ErrGasPriceTooLow
		}

	} else {
		// if gas price is not provided, we set it to minGasPrice
		gasPrice = networkConfig.MinGasPrice
	}

	suggestedFee := big.NewInt(0).Mul(
		big.NewInt(0).SetUint64(gasPrice),
		big.NewInt(0).SetUint64(gasLimit),
	)

	suggestedFee, gasPrice = adjustTxFeeWithFeeMultiplier(suggestedFee, gasPrice, options)

	return suggestedFee, gasPrice, gasLimit, nil
}

func adjustTxFeeWithFeeMultiplier(txFee *big.Int, gasPrice uint64, options objectsMap) (*big.Int, uint64) {
	feeMultiplierI, ok := options["feeMultiplier"]
	if !ok {
		return txFee, gasPrice
	}

	feeMultiplier, ok := feeMultiplierI.(float64)
	if !ok {
		return txFee, gasPrice
	}

	feeMultiplierBig := big.NewFloat(feeMultiplier)
	bigVal, ok := big.NewFloat(0).SetString(txFee.String())
	if !ok {
		return txFee, gasPrice
	}

	bigVal.Mul(bigVal, feeMultiplierBig)

	result := new(big.Int)
	bigVal.Int(result)

	gasPrice = uint64(feeMultiplier * float64(gasPrice))

	return result, gasPrice
}

func estimateGasLimit(operationType string, networkConfig *provider.NetworkConfig, options objectsMap) (uint64, *types.Error) {
	gasForDataField := uint64(0)
	if dataFieldI, ok := options["data"]; ok {
		dataField := fmt.Sprintf("%v", dataFieldI)
		gasForDataField = networkConfig.GasPerDataByte * uint64(len(dataField))
	}

	switch operationType {
	case opTransfer:
		return networkConfig.MinGasLimit + gasForDataField, nil
	default:
		//  we do not support this yet other operation types, but we might support it in the future
		return 0, ErrNotImplemented
	}
}

func checkProvidedGasLimit(providedGasLimit uint64, txType string, options objectsMap, networkConfig *provider.NetworkConfig) *types.Error {
	estimatedGasLimit, err := estimateGasLimit(txType, networkConfig, options)
	if err != nil {
		return err
	}

	if providedGasLimit < estimatedGasLimit {
		return ErrInsufficientGasLimit
	}

	return nil
}

func getUint64Value(obj interface{}) uint64 {
	if value, ok := obj.(uint64); ok {
		return value
	}
	if value, ok := obj.(float64); ok {
		return uint64(value)
	}

	return 0
}
