package metrics

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

// statusMetrics will handle displaying at /status/metrics all collected metrics
type statusMetrics struct {
	endpointMetrics        map[string]*data.EndpointMetrics
	mutEndpointsOperations sync.RWMutex
}

// NewStatusMetrics will return an instance of the struct
func NewStatusMetrics() *statusMetrics {
	return &statusMetrics{
		endpointMetrics: make(map[string]*data.EndpointMetrics),
	}
}

// AddRequestData will add the received data to the metrics map
func (sm *statusMetrics) AddRequestData(path string, withError bool, duration time.Duration) {
	// TODO: refactor this by using a buffered channel that receives new request data and stores them into the map
	// from time to time

	sm.mutEndpointsOperations.Lock()
	defer sm.mutEndpointsOperations.Unlock()

	currentData := sm.endpointMetrics[path]
	withErrorIncrementalStep := uint64(0)
	if withError {
		withErrorIncrementalStep = 1
	}
	if currentData == nil {
		sm.endpointMetrics[path] = &data.EndpointMetrics{
			NumRequests:         1,
			NumErrors:           withErrorIncrementalStep,
			TotalResponseTime:   duration,
			LowestResponseTime:  duration,
			HighestResponseTime: duration,
		}

		return
	}

	currentData.NumRequests++
	currentData.NumErrors += withErrorIncrementalStep
	if duration < currentData.LowestResponseTime {
		currentData.LowestResponseTime = duration
	}
	if duration > currentData.HighestResponseTime {
		currentData.HighestResponseTime = duration
	}
	currentData.TotalResponseTime += duration
}

// GetAll returns the metrics map
func (sm *statusMetrics) GetAll() map[string]*data.EndpointMetrics {
	sm.mutEndpointsOperations.RLock()
	defer sm.mutEndpointsOperations.RUnlock()

	newMap := make(map[string]*data.EndpointMetrics)
	for key, value := range sm.endpointMetrics {
		newMap[key] = value
	}

	return newMap
}

// GetMetricsForPrometheus returns the metrics in a prometheus format
func (sm *statusMetrics) GetMetricsForPrometheus() string {
	metricsMap := sm.GetAll()

	stringBuilder := strings.Builder{}

	for endpointPath, endpointData := range metricsMap {
		stringBuilder.WriteString(fmt.Sprintf("num_requests{endpoint=\"%s\"} %d\n", endpointPath, endpointData.NumRequests))
		stringBuilder.WriteString(fmt.Sprintf("num_errors{endpoint=\"%s\"} %d\n", endpointPath, endpointData.NumErrors))
		stringBuilder.WriteString(fmt.Sprintf("total_response_time_ns{endpoint=\"%s\"} %d\n", endpointPath, endpointData.TotalResponseTime))
		stringBuilder.WriteString(fmt.Sprintf("highest_response_time_ns{endpoint=\"%s\"} %d\n", endpointPath, endpointData.HighestResponseTime))
		stringBuilder.WriteString(fmt.Sprintf("lowest_response_time_ns{endpoint=\"%s\"} %d\n", endpointPath, endpointData.LowestResponseTime))
	}

	return stringBuilder.String()
}

// IsInterfaceNil returns true if there is no value under the interface
func (sm *statusMetrics) IsInterfaceNil() bool {
	return sm == nil
}
