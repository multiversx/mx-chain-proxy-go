package observer

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type baseNodeProvider struct {
	mutNodes sync.RWMutex
	nodes    map[uint32][]*data.NodeData
	allNodes []*data.NodeData
}

func (bop *baseNodeProvider) initNodesMaps(nodes []*data.NodeData) error {
	if len(nodes) == 0 {
		return ErrEmptyObserversList
	}

	newNodes := make(map[uint32][]*data.NodeData)
	for _, observer := range nodes {
		shardId := observer.ShardId
		newNodes[shardId] = append(newNodes[shardId], observer)
	}

	bop.mutNodes.Lock()
	bop.nodes = newNodes
	bop.allNodes = nodes
	bop.mutNodes.Unlock()

	return nil
}
