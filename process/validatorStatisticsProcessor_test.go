package process_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewValidatorStatisticsProcessor_NilProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewValidatorStatisticsProcessor(nil, &mock.ValStatsCacherMock{}, time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewValidatorStatisticsProcessor_NilCacherShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{}, nil, time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilValidatorStatisticsCacher, err)
}

func TestNewValidatorStatisticsProcessor_InvalidCacheValidityDurationShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{}, &mock.ValStatsCacherMock{}, -time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrInvalidCacheValidityDuration, err)
}

func TestNewValidatorStatisticsProcessor_WithOkProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hbp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{}, &mock.ValStatsCacherMock{}, time.Second)

	assert.NotNil(t, hbp)
	assert.Nil(t, err)
}

func TestValidatorStatisticsProcessor_GetValidatorStatisticsDataWrongValuesShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{}, &mock.ValStatsCacherMock{}, time.Second)
	assert.Nil(t, err)

	res, err := hp.GetValidatorStatistics()

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestValidatorStatisticsProcessor_GetValidatorStatisticsDataOkValuesShouldPass(t *testing.T) {
	t.Parallel()

	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			var obs []*data.NodeData
			obs = append(obs, &data.NodeData{
				ShardId: core.MetachainShardId,
				Address: "addr",
			})
			return obs, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 0, nil
		},
	},
		&mock.ValStatsCacherMock{},
		time.Second,
	)

	assert.Nil(t, err)

	res, err := hp.GetValidatorStatistics()
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestValidatorStatisticsProcessor_GetValidatorStatisticsNoMetaObserverShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{
		GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			var obs []*data.NodeData
			obs = append(obs, &data.NodeData{
				ShardId: 1,
				Address: "addr",
			})
			return obs, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 0, nil
		},
	},
		&mock.ValStatsCacherMock{},
		time.Second,
	)

	assert.Nil(t, err)

	res, err := hp.GetValidatorStatistics()
	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestValidatorStatisticsProcessor_GetValidatorStatisticsShouldReturnDataFromApiBecauseCacheDataIsNil(t *testing.T) {
	t.Parallel()

	httpWasCalled := false
	// set nil hbts response in cache
	cacher := &mock.ValStatsCacherMock{Data: nil}
	hp, err := process.NewValidatorStatisticsProcessor(
		&mock.ProcessorStub{
			GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{{Address: "obs1", ShardId: core.MetachainShardId}}, nil
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

	_, err = hp.GetValidatorStatistics()
	assert.Nil(t, err)
	assert.True(t, httpWasCalled)
}

func TestValidatorStatisticsProcessor_GetValidatorStatisticsShouldReturnDataFromCacher(t *testing.T) {
	t.Parallel()

	valStatsMap := map[string]*data.ValidatorApiResponse{
		"key0": {TempRating: 50.7},
	}
	cacher := &mock.ValStatsCacherMock{Data: valStatsMap}
	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{}, cacher, time.Millisecond)
	assert.Nil(t, err)

	res, err := hp.GetValidatorStatistics()

	assert.Nil(t, err)
	assert.Equal(t, res.Statistics, valStatsMap)
}

func TestValidatorStatisticsProcessor_CacheShouldUpdate(t *testing.T) {
	t.Parallel()

	numOfTimesHttpWasCalled := int32(0)
	cacher := &mock.ValStatsCacherMock{}
	hp, err := process.NewValidatorStatisticsProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{Address: "obs1", ShardId: core.MetachainShardId}}, nil
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
