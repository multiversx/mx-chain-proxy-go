package services

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/stretchr/testify/assert"
)

func TestEstimateGasLimit(t *testing.T) {
	t.Parallel()

	minGasLimit := uint64(1000)
	gasPerDataByte := uint64(100)
	networkConfig := &client.NetworkConfig{
		GasPerDataByte: gasPerDataByte,
		MinGasLimit:    minGasLimit,
	}

	dataField := "transaction-data"
	options := objectsMap{
		"data": dataField,
	}

	expectedGasLimit := minGasLimit + uint64(len(dataField))*gasPerDataByte

	gasLimit, err := estimateGasLimit(opTransfer, networkConfig, options)
	assert.Nil(t, err)
	assert.Equal(t, expectedGasLimit, gasLimit)

	gasLimit, err = estimateGasLimit(opTransfer, networkConfig, nil)
	assert.Nil(t, err)
	assert.Equal(t, minGasLimit, gasLimit)

	// unsupported operation type you cannot estimate gasLimit for a reward operation
	// reward operation can be generated only by the network not by a user
	gasLimit, err = estimateGasLimit(opReward, networkConfig, nil)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Equal(t, uint64(0), gasLimit)
}

func TestProvidedGasLimit(t *testing.T) {
	t.Parallel()

	minGasLimit := uint64(1000)
	gasPerDataByte := uint64(100)
	networkConfig := &client.NetworkConfig{
		GasPerDataByte: gasPerDataByte,
		MinGasLimit:    minGasLimit,
	}

	dataField := "transaction-data"
	options := objectsMap{
		"data": dataField,
	}

	err := checkProvidedGasLimit(uint64(900), opTransfer, options, networkConfig)
	assert.Equal(t, ErrInsufficientGasLimit, err)

	err = checkProvidedGasLimit(uint64(900), opReward, options, networkConfig)
	assert.Equal(t, ErrNotImplemented, err)

	err = checkProvidedGasLimit(uint64(9000), opTransfer, options, networkConfig)
	assert.Nil(t, err)
}

func TestAdjustTxFeeWithFeeMultiplier(t *testing.T) {
	t.Parallel()

	options := objectsMap{
		"feeMultiplier": 1.1,
	}

	expectedFee := "1100"
	suggestedFee := big.NewInt(1000)

	suggestedFeeResult := adjustTxFeeWithFeeMultiplier(suggestedFee, options)
	assert.Equal(t, expectedFee, suggestedFeeResult.String())

	expectedFee = "1000"
	suggestedFeeResult = adjustTxFeeWithFeeMultiplier(suggestedFee, make(objectsMap))
	assert.Equal(t, expectedFee, suggestedFeeResult.String())
}
