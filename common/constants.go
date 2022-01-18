package common

// OutputFormat represents the format type returned by api
type OutputFormat uint8

const (
	// Internal output format returns struct directly, will be serialized into JSON by gin
	Internal OutputFormat = 0

	// Proto output format returns the bytes of the proto object
	Proto OutputFormat = 1
)
