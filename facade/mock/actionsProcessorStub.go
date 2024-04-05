package mock

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// ActionsProcessorStub -
type ActionsProcessorStub struct {
	ReloadObserversCalled            func() data.NodesReloadResponse
	ReloadFullHistoryObserversCalled func() data.NodesReloadResponse
}

// ReloadObservers -
func (a *ActionsProcessorStub) ReloadObservers() data.NodesReloadResponse {
	if a.ReloadObserversCalled != nil {
		return a.ReloadObserversCalled()
	}

	return data.NodesReloadResponse{}
}

// ReloadFullHistoryObservers -
func (a *ActionsProcessorStub) ReloadFullHistoryObservers() data.NodesReloadResponse {
	if a.ReloadFullHistoryObserversCalled != nil {
		return a.ReloadFullHistoryObserversCalled()
	}

	return data.NodesReloadResponse{}
}
