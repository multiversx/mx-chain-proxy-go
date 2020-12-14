package mock

import (
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// SCQueryServiceStub is a stub
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*data.SCQuery) (*vm.VMOutputApi, int, error)
}

// ExecuteQuery is a stub
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, int, error) {
	return serviceStub.ExecuteQueryCalled(query)
}
