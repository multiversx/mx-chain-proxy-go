package mock

import (
	"github.com/ElrondNetwork/elrond-go/process"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// SCQueryServiceStub is a stub
type SCQueryServiceStub struct {
	ExecuteQueryCalled func(*process.SCQuery) (*vmcommon.VMOutput, error)
}

// ExecuteQuery is a stub
func (serviceStub *SCQueryServiceStub) ExecuteQuery(query *process.SCQuery) (*vmcommon.VMOutput, error) {
	return serviceStub.ExecuteQueryCalled(query)
}
