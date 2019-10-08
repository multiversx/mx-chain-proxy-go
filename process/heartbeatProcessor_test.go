package process_test

import (
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewHeartbeatProcessor_NilProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(nil)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewHeartbeatProcessor_WithOkProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hbp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{})

	assert.NotNil(t, hbp)
	assert.Nil(t, err)
}

func TestHeartbeatProcessor_GetHeartbeatDataWrongValuesShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{})
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
		CallGetRestEndPointWithTimeoutCalled: func(address string, path string, value interface{}, timeout time.Duration) error {
			return nil
		},
	})

	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()
	assert.NotNil(t, res)
	assert.Nil(t, err)
}
