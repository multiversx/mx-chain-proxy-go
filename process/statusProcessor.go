package process

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// StatusProcessor is able to process status requests
type StatusProcessor struct {
	proc                  Processor
	statusMetricsProvider StatusMetricsProvider
}

// NewStatusProcessor creates a new instance of AccountProcessor
func NewStatusProcessor(proc Processor, statusMetricsProvider StatusMetricsProvider) (*StatusProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(statusMetricsProvider) {
		return nil, ErrNilStatusMetricsProvider
	}

	return &StatusProcessor{
		proc:                  proc,
		statusMetricsProvider: statusMetricsProvider,
	}, nil
}

// GetMetrics returns the metrics for all the endpoints
func (sp *StatusProcessor) GetMetrics() map[string]*data.EndpointMetrics {
	return sp.statusMetricsProvider.GetAll()
}

// GetMetricsForPrometheus returns the metrics in a prometheus format
func (sp *StatusProcessor) GetMetricsForPrometheus() string {
	return sp.statusMetricsProvider.GetMetricsForPrometheus()
}
