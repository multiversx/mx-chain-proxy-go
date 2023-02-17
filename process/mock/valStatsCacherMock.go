package mock

import (
	"errors"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

// ValStatsCacherMock --
type ValStatsCacherMock struct {
	Data map[string]*data.ValidatorApiResponse
}

// LoadValStats --
func (vscm *ValStatsCacherMock) LoadValStats() (map[string]*data.ValidatorApiResponse, error) {
	if vscm.Data == nil {
		return nil, errors.New("nil Data")
	}

	return vscm.Data, nil
}

// StoreValStats --
func (vscm *ValStatsCacherMock) StoreValStats(valStats map[string]*data.ValidatorApiResponse) error {
	vscm.Data = valStats
	return nil
}

// IsInterfaceNil --
func (vscm *ValStatsCacherMock) IsInterfaceNil() bool {
	return vscm == nil
}
