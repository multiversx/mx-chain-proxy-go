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
	// OpStatusSuccess is the operation status for successful operations.
	OpStatusSuccess = "Success"
	// OpStatusFailed is the operation status for failed operations.
	OpStatusFailed = "Failed"
)

// TxProcessingType represents a processing type for transactions
type TxProcessingType string

const (
	// TxProcessingTypeMoveBalance refers to a "MoveBalance"
	TxProcessingTypeMoveBalance TxProcessingType = "MoveBalance"
)
