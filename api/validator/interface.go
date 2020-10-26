package validator

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// FacadeHandler interface defines methods that can be used from facade context variable
type FacadeHandler interface {
	ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error)
}
