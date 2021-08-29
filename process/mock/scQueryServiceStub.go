package mock

import (
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// SCQueryServiceStub is a stub
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*data.SCQuery) (*vm.VMOutputApi, error)
}

// ExecuteQuery is a stub
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, error) {
	return serviceStub.ExecuteQueryCalled(query)
}

// IsInterfaceNil returns true if the value under the interface is nil
func (serviceStub *SCQueryServiceStub) IsInterfaceNil() bool {
	return serviceStub == nil
}
