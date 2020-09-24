package services

import (
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
