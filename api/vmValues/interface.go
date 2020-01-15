package vmValues

import (
	"github.com/ElrondNetwork/elrond-go/process"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// FacadeHandler interface defines methods that can be used from `elrondFacade` context variable
type FacadeHandler interface {
	ExecuteSCQuery(*process.SCQuery) (*vmcommon.VMOutput, error)
}
