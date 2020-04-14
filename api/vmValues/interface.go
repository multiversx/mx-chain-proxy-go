package vmValues

import (
	"github.com/ElrondNetwork/elrond-proxy-go/shared"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// FacadeHandler interface defines methods that can be used from `elrondFacade` context variable
type FacadeHandler interface {
	ExecuteSCQuery(*shared.SCQuery) (*vmcommon.VMOutput, error)
}
