package resultsParser

type (
	// WellKnownTopics is an enum that contains the most common topics encapsulated in a smart contract.
	WellKnownTopics = string

	// WellKnownEvents is an enum that contains the most common events encapsulated in a smart contract.
	WellKnownEvents = string
)

const (
	// OnTransactionCompleted is the event when a transaction is completed.
	OnTransactionCompleted WellKnownEvents = "completedTxEvent"

	// OnSignalError is the event when a smart contract encounters an error.
	OnSignalError WellKnownEvents = "signalError"

	// OnWriteLog is the event where a smart contract will write to the Logs.
	OnWriteLog WellKnownEvents = "writeLog"

	// TooMuchGas is the topic when too much gas is provided for processing a smart contract.
	TooMuchGas WellKnownTopics = "@too much gas provided for processing"
)
