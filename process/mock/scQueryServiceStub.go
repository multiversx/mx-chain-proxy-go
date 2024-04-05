package mock

import (
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// SCQueryServiceStub is a stub
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error)
}

// ExecuteQuery is a stub
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	return serviceStub.ExecuteQueryCalled(query)
}

// IsInterfaceNil returns true if the value under the interface is nil
func (serviceStub *SCQueryServiceStub) IsInterfaceNil() bool {
	return serviceStub == nil
}
