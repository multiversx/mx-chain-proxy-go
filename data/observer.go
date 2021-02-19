package data

// NodeData holds an observer data
type NodeData struct {
	ShardId uint32
	Address string
}

// NodesReloadResponse is a DTO that holds details about nodes reloading
type NodesReloadResponse struct {
	OkRequest   bool
	Description string
	Error       string
}

// NodeType is a type which identifies the type of a node (observer or full history)
type NodeType string

const (
	// Observer identifies a node which is a regular observer
	Observer NodeType = "observer"

	// FullHistoryNode identifier a node that has full history mode enabled
	FullHistoryNode NodeType = "full history"
)
