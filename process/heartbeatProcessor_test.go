package process_test

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"testing"
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
		GetFirstAvailableObserverCalled: func() (*data.Observer, error) {
			return &data.Observer{
				ShardId: 0,
				Address: "addr",
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return nil
		},
	})

	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()

	assert.NotNil(t, res)
	assert.Nil(t, err)
}
