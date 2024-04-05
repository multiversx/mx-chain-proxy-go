package mock

import (
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// SCQueryServiceStub -
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error)
}

// ExecuteQuery -
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	return serviceStub.ExecuteQueryCalled(query)
}
