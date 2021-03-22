package mock

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GenericApiResponseCacherMock -
type GenericApiResponseCacherMock struct {
	Data *data.GenericAPIResponse
}

// Load -
func (g *GenericApiResponseCacherMock) Load() (*data.GenericAPIResponse, error) {
	if g.Data == nil {
		return nil, errors.New("nil data")
	}

	return g.Data, nil
}

// Store -
func (g *GenericApiResponseCacherMock) Store(response *data.GenericAPIResponse) {
	g.Data = response
}

// IsInterfaceNil -
func (g *GenericApiResponseCacherMock) IsInterfaceNil() bool {
	return g == nil
}
