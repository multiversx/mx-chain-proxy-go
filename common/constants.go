package common

// OutportFormat represents the format type returned by api
type OutportFormat uint8

const (
	// Internal outport format returns struct directly, will be serialized into JSON by gin
	Internal OutportFormat = 0

	// Proto outport format returns the bytes of the proto object
	Proto OutportFormat = 1
)
