package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// StatusProcessorStub -
type StatusProcessorStub struct {
	GetMetricsCalled func() (map[string]*data.EndpointMetrics, error)
}

// GetMetrics -
func (s *StatusProcessorStub) GetMetrics() (map[string]*data.EndpointMetrics, error) {
	if s.GetMetricsCalled != nil {
		return s.GetMetricsCalled()
	}

	return nil, nil
}
