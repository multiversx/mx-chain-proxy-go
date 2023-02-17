package mock

import (
	"errors"
	"sync"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

// GenericApiResponseCacherMock -
type GenericApiResponseCacherMock struct {
	Data *data.GenericAPIResponse
	sync.RWMutex
}

// Load -
func (g *GenericApiResponseCacherMock) Load() (*data.GenericAPIResponse, error) {
	g.RLock()
	defer g.RUnlock()

	if g.Data == nil {
		return nil, errors.New("nil data")
	}

	return g.Data, nil
}

// Store -
func (g *GenericApiResponseCacherMock) Store(response *data.GenericAPIResponse) {
	g.Lock()
	g.Data = response
	g.Unlock()
}

// IsInterfaceNil -
func (g *GenericApiResponseCacherMock) IsInterfaceNil() bool {
	return g == nil
}
