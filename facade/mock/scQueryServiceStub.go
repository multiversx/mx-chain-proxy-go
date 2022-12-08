package mock

import (
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// SCQueryServiceStub -
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*data.SCQuery) (*vm.VMOutputApi, error)
}

// ExecuteQuery -
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, error) {
	return serviceStub.ExecuteQueryCalled(query)
}
