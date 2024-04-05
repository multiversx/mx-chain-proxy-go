package process

import (
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewStatusProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil base processor - should error", func(t *testing.T) {
		t.Parallel()

		sp, err := NewStatusProcessor(nil, &mock.StatusMetricsProviderStub{})
		require.Nil(t, sp)
		require.Equal(t, ErrNilCoreProcessor, err)
	})

	t.Run("nil status metric provider - should error", func(t *testing.T) {
		t.Parallel()

		sp, err := NewStatusProcessor(&mock.ProcessorStub{}, nil)
		require.Nil(t, sp)
		require.Equal(t, ErrNilStatusMetricsProvider, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		sp, err := NewStatusProcessor(&mock.ProcessorStub{}, &mock.StatusMetricsProviderStub{})
		require.NoError(t, err)
		require.NotNil(t, sp)
	})
}

func TestStatusProcessor_GetMetrics(t *testing.T) {
	t.Parallel()

	expectedMetrics := map[string]*data.EndpointMetrics{
		"endpoint0": {NumErrors: 5},
		"endpoint1": {NumErrors: 37},
	}
	statusProvider := &mock.StatusMetricsProviderStub{
		GetAllCalled: func() map[string]*data.EndpointMetrics {
			return expectedMetrics
		},
	}
	sp, err := NewStatusProcessor(&mock.ProcessorStub{}, statusProvider)
	require.NoError(t, err)
	require.NotNil(t, sp)

	metrics := sp.GetMetrics()
	require.NoError(t, err)
	require.Equal(t, expectedMetrics, metrics)
}

func TestStatusProcessor_GetMetricsForPrometheus(t *testing.T) {
	t.Parallel()

	expectedOutput := "metrics"
	statusProvider := &mock.StatusMetricsProviderStub{
		GetMetricsForPrometheusCalled: func() string {
			return expectedOutput
		},
	}
	sp, err := NewStatusProcessor(&mock.ProcessorStub{}, statusProvider)
	require.NoError(t, err)
	require.NotNil(t, sp)

	metrics := sp.GetMetricsForPrometheus()
	require.NoError(t, err)
	require.Equal(t, expectedOutput, metrics)
}
