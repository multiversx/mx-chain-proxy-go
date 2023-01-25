package mock

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// StatusProcessorStub -
type StatusProcessorStub struct {
	GetMetricsCalled              func() map[string]*data.EndpointMetrics
	GetMetricsForPrometheusCalled func() string
}

// GetMetricsForPrometheus -
func (s *StatusProcessorStub) GetMetricsForPrometheus() string {
	if s.GetMetricsForPrometheusCalled != nil {
		return s.GetMetricsForPrometheusCalled()
	}

	return ""
}

// GetMetrics -
func (s *StatusProcessorStub) GetMetrics() map[string]*data.EndpointMetrics {
	if s.GetMetricsCalled != nil {
		return s.GetMetricsCalled()
	}

	return nil
}
