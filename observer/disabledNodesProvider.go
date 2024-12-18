package observer

import (
	"errors"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

type disabledNodesProvider struct {
	returnMessage string
}

func NewDisabledNodesProvider(returnMessage string) *disabledNodesProvider {
	returnMessageToUse := "not implemented"
	if returnMessage != "" {
		returnMessageToUse = returnMessage
	}
	return &disabledNodesProvider{
		returnMessage: returnMessageToUse,
	}
}

// UpdateNodesBasedOnSyncState won't do anything as this is a disabled component
func (d *disabledNodesProvider) UpdateNodesBasedOnSyncState(_ []*data.NodeData) {
}

// GetAllNodesWithSyncState returns an empty slice
func (d *disabledNodesProvider) GetAllNodesWithSyncState() []*data.NodeData {
	return make([]*data.NodeData, 0)
}

// GetNodesByShardId returns the desired return message as an error
func (d *disabledNodesProvider) GetNodesByShardId(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	return nil, errors.New(d.returnMessage)
}

// GetAllNodes returns the desired return message as an error
func (d *disabledNodesProvider) GetAllNodes(_ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	return nil, errors.New(d.returnMessage)
}

// ReloadNodes return the desired return message as an error
func (d *disabledNodesProvider) ReloadNodes(_ data.NodeType) data.NodesReloadResponse {
	return data.NodesReloadResponse{Description: "disabled nodes provider", Error: d.returnMessage}
}

// PrintNodesInShards does nothing as it is disabled
func (d *disabledNodesProvider) PrintNodesInShards() {
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *disabledNodesProvider) IsInterfaceNil() bool {
	return d == nil
}
