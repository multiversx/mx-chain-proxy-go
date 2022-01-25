package process

import (
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
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
	metrics := sp.statusMetricsProvider.GetAll()

	return metrics
}

// GetMetricsForPrometheus returns the metrics in a prometheus format
func (sp *StatusProcessor) GetMetricsForPrometheus() string {
	return sp.statusMetricsProvider.GetMetricsForPrometheus()
}
