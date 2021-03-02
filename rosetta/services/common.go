package services

// SupportedOperationTypes is a list of the supported operations.
var SupportedOperationTypes = []string{
	opTransfer, opFee, opReward, opScResult, opInvalid,
}

type objectsMap map[string]interface{}
