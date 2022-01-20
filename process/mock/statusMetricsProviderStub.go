package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// StatusMetricsProviderStub -
type StatusMetricsProviderStub struct {
	GetAllCalled func() map[string]*data.EndpointMetrics
}

// GetAll -
func (s *StatusMetricsProviderStub) GetAll() map[string]*data.EndpointMetrics {
	if s.GetAllCalled != nil {
		return s.GetAllCalled()
	}

	return make(map[string]*data.EndpointMetrics)
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *StatusMetricsProviderStub) IsInterfaceNil() bool {
	return s == nil
}
