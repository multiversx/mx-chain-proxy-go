package mock

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ValStatsCacherMock --
type ValStatsCacherMock struct {
	Data map[string]*data.ValidatorApiResponse
}

// Load --
func (vscm *ValStatsCacherMock) Load() (map[string]*data.ValidatorApiResponse, error) {
	if vscm.Data == nil {
		return nil, errors.New("nil Data")
	}

	return vscm.Data, nil
}

// Store --
func (vscm *ValStatsCacherMock) Store(valStats map[string]*data.ValidatorApiResponse) error {
	vscm.Data = valStats
	return nil
}

// IsInterfaceNil --
func (vscm *ValStatsCacherMock) IsInterfaceNil() bool {
	return vscm == nil
}
