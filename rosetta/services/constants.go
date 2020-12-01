package services

const (
	NumBlocksToGet = uint64(200)

	RosettaVersion = "1.4.5"
	NodeVersion    = "1.1.0"

	opTransfer = "Transfer"
	opFee      = "Fee"
	opReward   = "Reward"
	opScResult = "SmartContractResult"
	opInvalid  = "Invalid"
)

var (
	opsStatusSuccess = "Success"
	opStatusFailed   = "Failed"

	// OpStatusOK is the operation status for successful operations.
	OpStatusSuccess = &opsStatusSuccess
	// OpStatusFailed is the operation status for failed operations.
	OpStatusFailed = &opStatusFailed
)
