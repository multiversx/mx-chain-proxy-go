package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/shared"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// SCQueryServiceStub is a stub
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*shared.SCQuery) (*vmcommon.VMOutput, error)
}

// ExecuteQuery is a stub
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *shared.SCQuery) (*vmcommon.VMOutput, error) {
	return serviceStub.ExecuteQueryCalled(query)
}
