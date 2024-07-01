package resultsParser

type (
	WellKnownTopics = string
	WellKnownEvents = string
)

const (
	OnTransactionCompleted WellKnownEvents = "completedTxEvent"
	OnSignalError          WellKnownEvents = "signalError"
	OnWriteLog             WellKnownEvents = "writeLog"

	TooMuchGas WellKnownTopics = "@too much gas provided for processing"
)
