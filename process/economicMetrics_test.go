package process_test

import (
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeStatusProcessor_GetEconomicsDataMetricsShouldReturnDataFromCacher(t *testing.T) {
	t.Parallel()

	respInCache := &data.GenericAPIResponse{
		Data:  "test data",
		Error: "test error",
	}

	cacher := &mock.GenericApiResponseCacherMock{Data: respInCache}
	hp, err := process.NewNodeStatusProcessor(&mock.ProcessorStub{}, cacher, time.Millisecond)
	assert.Nil(t, err)

	res, err := hp.GetEconomicsDataMetrics()

	assert.Nil(t, err)
	assert.Equal(t, res, respInCache)
}

func TestNodeStatusProcessor_CacheShouldUpdate(t *testing.T) {
	t.Parallel()

	numOfTimesHttpWasCalled := int32(0)
	cacher := &mock.GenericApiResponseCacherMock{}
	hp, err := process.NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{Address: "obs1"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			atomic.AddInt32(&numOfTimesHttpWasCalled, 1)
			return 0, nil
		},
	},
		cacher,
		25*time.Millisecond)

	assert.Nil(t, err)
	hp.StartCacheUpdate()

	// cache will become invalid after 25 ms so check if it renews its data

	// >25 => update
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, int32(2), atomic.LoadInt32(&numOfTimesHttpWasCalled))

	// > 25 => update
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, int32(3), atomic.LoadInt32(&numOfTimesHttpWasCalled))

	// < 25 => don't update
	time.Sleep(5 * time.Millisecond)
	assert.Equal(t, int32(3), atomic.LoadInt32(&numOfTimesHttpWasCalled))
}

func TestNodeStatusProcessor_GetEconomicsDataMetricsShouldWork(t *testing.T) {
	t.Parallel()

	addressMeta := "address_meta"
	expectedResponse := &data.GenericAPIResponse{
		Data: map[string]interface{}{
			"erd_total_supply": "12345",
		},
	}

	nodeStatusProc, _ := process.NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, err error) {
			return []*data.NodeData{
				{Address: addressMeta, ShardId: core.MetachainShardId},
			}, nil
		},
		CallGetRestEndPointCalled: func(_ string, _ string, value interface{}) (int, error) {
			expectedResponseBytes, _ := json.Marshal(expectedResponse)
			return 200, json.Unmarshal(expectedResponseBytes, value)
		},
	},
		&mock.GenericApiResponseCacherMock{
			Data: &data.GenericAPIResponse{Data: "default response"},
		},
		time.Millisecond,
	)

	time.Sleep(2 * time.Millisecond)

	nodeStatusProc.StartCacheUpdate()

	time.Sleep(10 * time.Millisecond)

	actualResponse, err := nodeStatusProc.GetEconomicsDataMetrics()
	require.NoError(t, err)
	require.Equal(t, *expectedResponse, *actualResponse)
}
