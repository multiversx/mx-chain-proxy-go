package metrics

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestNewStatusMetrics(t *testing.T) {
	t.Parallel()

	sm := NewStatusMetrics()
	require.False(t, check.IfNil(sm))
}

func TestStatusMetrics_AddRequestData(t *testing.T) {
	t.Parallel()

	t.Run("test when only a metric exists for an endpoint", testFirstMetric)
	t.Run("test when multiple entries exist for an endpoint", testWhenMultipleMetrics)
	t.Run("test when multiple entries for multiple endpoints", testWhenMultipleEntriesForMultipleEndpoints)
}

func TestStatusMetrics_GetMetricsForPrometheus(t *testing.T) {
	t.Parallel()

	t.Run("test fetching metrics for prometheus", testMetricsForPrometheus)
}

func testFirstMetric(t *testing.T) {
	t.Parallel()

	sm := NewStatusMetrics()

	testEndpoint, testDuration := "/network/config", 1*time.Second
	sm.AddRequestData(testEndpoint, false, testDuration)

	res := sm.GetAll()
	require.Equal(t, res[testEndpoint], &data.EndpointMetrics{
		NumRequests:         1,
		NumErrors:           0,
		TotalResponseTime:   testDuration,
		LowestResponseTime:  testDuration,
		HighestResponseTime: testDuration,
	})
}

func testWhenMultipleMetrics(t *testing.T) {
	t.Parallel()

	sm := NewStatusMetrics()

	testEndpoint := "/network/config"
	testDuration0, testDuration1, testDuration2 := 4*time.Millisecond, 20*time.Millisecond, 2*time.Millisecond
	sm.AddRequestData(testEndpoint, false, testDuration0)
	sm.AddRequestData(testEndpoint, true, testDuration1)
	sm.AddRequestData(testEndpoint, false, testDuration2)

	res := sm.GetAll()
	require.Equal(t, res[testEndpoint], &data.EndpointMetrics{
		NumRequests:         3,
		NumErrors:           1,
		TotalResponseTime:   testDuration0 + testDuration1 + testDuration2,
		LowestResponseTime:  testDuration2,
		HighestResponseTime: testDuration1,
	})
}

func testWhenMultipleEntriesForMultipleEndpoints(t *testing.T) {
	t.Parallel()

	sm := NewStatusMetrics()

	testEndpoint0, testEndpoint1 := "/network/config", "/network/esdts"

	testDuration0End0, testDuration1End0 := time.Second, 5*time.Second
	testDuration0End1, testDuration1End1 := time.Hour, 4*time.Hour

	sm.AddRequestData(testEndpoint0, true, testDuration0End0)
	sm.AddRequestData(testEndpoint0, false, testDuration1End0)

	sm.AddRequestData(testEndpoint1, true, testDuration0End1)
	sm.AddRequestData(testEndpoint1, true, testDuration1End1)

	res := sm.GetAll()

	require.Len(t, res, 2)
	require.Equal(t, res[testEndpoint0], &data.EndpointMetrics{
		NumRequests:         2,
		NumErrors:           1,
		TotalResponseTime:   testDuration0End0 + testDuration1End0,
		LowestResponseTime:  testDuration0End0,
		HighestResponseTime: testDuration1End0,
	})
	require.Equal(t, res[testEndpoint1], &data.EndpointMetrics{
		NumRequests:         2,
		NumErrors:           2,
		TotalResponseTime:   testDuration0End1 + testDuration1End1,
		LowestResponseTime:  testDuration0End1,
		HighestResponseTime: testDuration1End1,
	})
}

func testMetricsForPrometheus(t *testing.T) {
	t.Parallel()

	sm := NewStatusMetrics()

	testEndpoint := "/network/config"
	testDuration0, testDuration1, testDuration2 := 4*time.Millisecond, 20*time.Millisecond, 2*time.Millisecond
	sm.AddRequestData(testEndpoint, false, testDuration0)
	sm.AddRequestData(testEndpoint, true, testDuration1)
	sm.AddRequestData(testEndpoint, false, testDuration2)

	res := sm.GetMetricsForPrometheus()

	expectedString := `num_requests{endpoint="/network/config"} 3
num_errors{endpoint="/network/config"} 1
total_response_time_ns{endpoint="/network/config"} 26000000
highest_response_time_ns{endpoint="/network/config"} 20000000
lowest_response_time_ns{endpoint="/network/config"} 2000000
`

	require.Equal(t, expectedString, res)
}

func TestStatusMetrics_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()
	sm := NewStatusMetrics()

	numIterations := 500
	wg := sync.WaitGroup{}
	wg.Add(numIterations)

	for i := 0; i < numIterations; i++ {
		go func(index int) {
			switch index % 4 {
			case 0:
				sm.AddRequestData(fmt.Sprintf("endpoint_%d", index%5), false, time.Hour*time.Duration(index))
			case 1:
				res := sm.GetAll()
				delete(res, "endpoint_0")
			case 2:
				_ = sm.GetMetricsForPrometheus()
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}
