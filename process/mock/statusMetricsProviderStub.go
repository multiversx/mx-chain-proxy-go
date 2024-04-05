package mock

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// StatusMetricsProviderStub -
type StatusMetricsProviderStub struct {
	GetAllCalled                  func() map[string]*data.EndpointMetrics
	GetMetricsForPrometheusCalled func() string
}

// GetMetricsForPrometheus -
func (s *StatusMetricsProviderStub) GetMetricsForPrometheus() string {
	if s.GetMetricsForPrometheusCalled != nil {
		return s.GetMetricsForPrometheusCalled()
	}

	return ""
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
