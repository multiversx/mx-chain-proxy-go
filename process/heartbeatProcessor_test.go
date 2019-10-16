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

func TestNewHeartbeatProcessor_NilProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(nil, &mock.HeartbeatCacherMock{}, time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewHeartbeatProcessor_NilCacherShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{}, nil, time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilHeartbeatCacher, err)
}

func TestNewHeartbeatProcessor_InvalidCacheValidityDurationShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{}, &mock.HeartbeatCacherMock{}, -time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrInvalidCacheValidityDuration, err)
}

func TestNewHeartbeatProcessor_WithOkProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hbp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{}, &mock.HeartbeatCacherMock{}, time.Second)

	assert.NotNil(t, hbp)
	assert.Nil(t, err)
}

func TestHeartbeatProcessor_GetHeartbeatDataWrongValuesShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{}, &mock.HeartbeatCacherMock{}, time.Second)
	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestHeartbeatProcessor_GetHeartbeatDataOkValuesShouldPass(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{
		GetAllObserversCalled: func() (observers []*data.Observer, e error) {
			var obs []*data.Observer
			obs = append(obs, &data.Observer{
				ShardId: 1,
				Address: "addr",
			})
			return obs, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return nil
		},
	},
		&mock.HeartbeatCacherMock{},
		time.Second,
	)

	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestHeartbeatProcessor_GetHeartbeatDataShouldReturnDataFromApiBecauseCacheDataIsNil(t *testing.T) {
	t.Parallel()

	numOfTimeHttpWasCalled := int32(0)
	// set nil hbts response in cache
	cacher := &mock.HeartbeatCacherMock{Data: nil}
	hp, err := process.NewHeartbeatProcessor(
		&mock.ProcessorStub{
			GetAllObserversCalled: func() ([]*data.Observer, error) {
				return []*data.Observer{{Address: "obs1"}}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
				atomic.AddInt32(&numOfTimeHttpWasCalled, 1)
				return nil
			},
		},
		cacher,
		time.Second,
	)
	assert.Nil(t, err)

	_, err = hp.GetHeartbeatData()
	assert.Nil(t, err)
	// expect 2 calls to the http: first one when updateCache() will start and second one when cache doesn't have
	// any value so it will load the data from api
	assert.Equal(t, int32(2), atomic.LoadInt32(&numOfTimeHttpWasCalled))
}

func TestHeartbeatProcessor_GetHeartbeatDataShouldReturnDataFromCacher(t *testing.T) {
	t.Parallel()

	hbtsResp := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node1",
			},
			{
				NodeDisplayName: "node2",
			},
		},
	}
	cacher := &mock.HeartbeatCacherMock{Data: &hbtsResp}
	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{}, cacher, time.Millisecond)
	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()

	assert.Nil(t, err)
	assert.Equal(t, *res, hbtsResp)
}
