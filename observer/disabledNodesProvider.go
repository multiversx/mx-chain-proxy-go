package observer

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
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

// GetNodesByShardId returns the desired return message as an error
func (d *disabledNodesProvider) GetNodesByShardId(_ uint32) ([]*data.NodeData, error) {
	return nil, errors.New(d.returnMessage)
}

// GetAllNodes returns the desired return message as an error
func (d *disabledNodesProvider) GetAllNodes() ([]*data.NodeData, error) {
	return nil, errors.New(d.returnMessage)
}

// IsInterfaceNil returns true if there is no value under the interface
func (d *disabledNodesProvider) IsInterfaceNil() bool {
	return d == nil
}
