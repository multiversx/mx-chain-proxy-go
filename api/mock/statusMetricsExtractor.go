package mock

import (
	"time"
)

// StatusMetricsExporterStub -
type StatusMetricsExporterStub struct {
	AddRequestDataCalled func(path string, withError bool, duration time.Duration)
}

// AddRequestData -
func (s *StatusMetricsExporterStub) AddRequestData(path string, withError bool, duration time.Duration) {
	if s.AddRequestDataCalled != nil {
		s.AddRequestDataCalled(path, withError, duration)
	}
}

// IsInterfaceNil -
func (s *StatusMetricsExporterStub) IsInterfaceNil() bool {
	return s == nil
}
