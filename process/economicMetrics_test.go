package process_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNodeStatusProcessor_GetEconomicsDataMetricsShouldReturnDataFromApiBecauseCacheDataIsNil(t *testing.T) {
	t.Parallel()

	httpWasCalled := false
	// set nil response in cache
	cacher := &mock.GenericApiResponseCacherMock{Data: nil}
	np, err := process.NewNodeStatusProcessor(
		&mock.ProcessorStub{
			GetObserversCalled: func(_ uint32) ([]*data.NodeData, error) {
				return []*data.NodeData{{Address: "obs1"}}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				httpWasCalled = true
				return 0, nil
			},
		},
		cacher,
		time.Second,
	)
	assert.Nil(t, err)

	_, err = np.GetEconomicsDataMetrics()
	assert.Nil(t, err)
	assert.True(t, httpWasCalled)
}

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
		GetObserversCalled: func(_ uint32) ([]*data.NodeData, error) {
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
