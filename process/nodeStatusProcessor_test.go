package process

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewNodeStatusProcessor_NilBaseProcessor(t *testing.T) {
	t.Parallel()

	nodeStatusProc, err := NewNodeStatusProcessor(nil)

	require.Equal(t, ErrNilCoreProcessor, err)
	require.Nil(t, nodeStatusProc)
}

func TestNodeStatusProcessor_GetShardStatusNoObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return nil, localErr
		},
	})

	status, err := nodeStatusProc.GetShardStatus(0)
	require.Equal(t, localErr, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetShardStatusGetRestEndPointError(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return localErr
		},
	})

	status, err := nodeStatusProc.GetShardStatus(0)
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetShardStatus(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			localMap := map[string]interface{}{
				"key": 1,
			}
			localMapBytes, _ := json.Marshal(localMap)

			return json.Unmarshal(localMapBytes, value)
		},
	})

	statusMap, err := nodeStatusProc.GetShardStatus(0)
	require.Nil(t, err)
	require.NotNil(t, statusMap)

	valueFromMap, ok := statusMap["key"]
	require.True(t, ok)
	require.Equal(t, 1, int(valueFromMap.(float64)))

}

func TestNodeStatusProcessor_GetEpochMetricsGetObserversFailedShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return nil, localErr
		},
	})

	status, err := nodeStatusProc.GetEpochMetrics(0)
	require.Equal(t, localErr, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetEpochMetricsGetRestEndPointError(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return localErr
		},
	})

	status, err := nodeStatusProc.GetEpochMetrics(0)
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetEpochMetrics(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			localMap := map[string]interface{}{
				"key": 1,
			}
			localMapBytes, _ := json.Marshal(localMap)

			return json.Unmarshal(localMapBytes, value)
		},
	})

	statusMap, err := nodeStatusProc.GetEpochMetrics(0)
	require.Nil(t, err)
	require.NotNil(t, statusMap)

	valueFromMap, ok := statusMap["key"]
	require.True(t, ok)
	require.Equal(t, 1, int(valueFromMap.(float64)))

}
