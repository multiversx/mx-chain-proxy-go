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

func TestNodeStatusProcessor_GetConfigMetricsGetRestEndPointError(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetAllObserversCalled: func() []*data.Observer {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
			}
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return localErr
		},
	})

	status, err := nodeStatusProc.GetNetworkConfigMetrics()
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetConfigMetrics(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetAllObserversCalled: func() []*data.Observer {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
			}
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			localMap := map[string]interface{}{
				"key": 1,
			}
			genericResp := &data.GenericAPIResponse{Data: localMap}
			genRespBytes, _ := json.Marshal(genericResp)

			return json.Unmarshal(genRespBytes, value)
		},
	})

	genericResponse, err := nodeStatusProc.GetNetworkConfigMetrics()
	require.Nil(t, err)
	require.NotNil(t, genericResponse)

	map1, ok := genericResponse.Data.(map[string]interface{})
	require.True(t, ok)

	valueFromMap, ok := map1["key"]
	require.True(t, ok)
	require.Equal(t, 1, int(valueFromMap.(float64)))

}

func TestNodeStatusProcessor_GetNetworkMetricsGetObserversFailedShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, err error) {
			return nil, localErr
		},
	})

	status, err := nodeStatusProc.GetNetworkStatusMetrics(0)
	require.Equal(t, localErr, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetNetworkMetricsGetRestEndPointError(t *testing.T) {
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

	status, err := nodeStatusProc.GetNetworkStatusMetrics(0)
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetNetworkMetrics(t *testing.T) {
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
			genericResp := &data.GenericAPIResponse{Data: localMap}
			genRespBytes, _ := json.Marshal(genericResp)

			return json.Unmarshal(genRespBytes, value)
		},
	})

	genericResponse, err := nodeStatusProc.GetNetworkStatusMetrics(0)
	require.Nil(t, err)
	require.NotNil(t, genericResponse)

	map1, ok := genericResponse.Data.(map[string]interface{})
	require.True(t, ok)

	valueFromMap, ok := map1["key"]
	require.True(t, ok)
	require.Equal(t, 1, int(valueFromMap.(float64)))

}
