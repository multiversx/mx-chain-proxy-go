package metrics

import (
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
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
	// TODO: analyze possible way of improving this function - launch on goroutines (with a max num of goroutines),
	// or implement a queue mechanism for writing new data. Currently, the addition is done sequentially.

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

func (sm *statusMetrics) GetAll() map[string]*data.EndpointMetrics {
	sm.mutEndpointsOperations.RLock()
	defer sm.mutEndpointsOperations.RUnlock()

	return sm.endpointMetrics
}

// IsInterfaceNil returns true if there is no value under the interface
func (sm *statusMetrics) IsInterfaceNil() bool {
	return sm == nil
}
