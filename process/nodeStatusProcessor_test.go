package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
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
		GetAllObserversCalled: func() ([]*data.NodeData, error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 0, localErr
		},
	})

	status, err := nodeStatusProc.GetNetworkConfigMetrics()
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetConfigMetrics(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetAllObserversCalled: func() ([]*data.NodeData, error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			localMap := map[string]interface{}{
				"key": 1,
			}
			genericResp := &data.GenericAPIResponse{Data: localMap}
			genRespBytes, _ := json.Marshal(genericResp)

			return 0, json.Unmarshal(genRespBytes, value)
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

func TestNodeStatusProcessor_GetTotalStakedErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: core.MetachainShardId},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 0, localErr
		},
	})

	genericResponse, err := nodeStatusProc.GetTotalStaked()
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, genericResponse)
}

func TestNodeStatusProcessor_GetTotalStakedShouldWork(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: core.MetachainShardId},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			localMap := map[string]interface{}{
				"totalStakedValue": "250000000",
			}
			genericResp := &data.GenericAPIResponse{Data: localMap}
			genRespBytes, _ := json.Marshal(genericResp)

			return 0, json.Unmarshal(genRespBytes, value)
		},
	})

	genericResponse, err := nodeStatusProc.GetTotalStaked()
	require.Nil(t, err)
	require.NotNil(t, genericResponse)

	map1, ok := genericResponse.Data.(map[string]interface{})
	require.True(t, ok)

	valueFromMap, ok := map1["totalStakedValue"]
	require.True(t, ok)
	require.Equal(t, "250000000", fmt.Sprintf("%v", valueFromMap))
}

func TestNodeStatusProcessor_GetNetworkMetricsGetObserversFailedShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
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
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 0, localErr
		},
	})

	status, err := nodeStatusProc.GetNetworkStatusMetrics(0)
	require.Equal(t, ErrSendingRequest, err)
	require.Nil(t, status)
}

func TestNodeStatusProcessor_GetNetworkMetrics(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			localMap := map[string]interface{}{
				"key": 1,
			}
			genericResp := &data.GenericAPIResponse{Data: localMap}
			genRespBytes, _ := json.Marshal(genericResp)

			return 0, json.Unmarshal(genRespBytes, value)
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

func TestNodeStatusProcessor_GetLatestBlockNonce(t *testing.T) {
	t.Parallel()

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetAllObserversCalled: func() (observers []*data.NodeData, err error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: core.MetachainShardId},
			}, nil
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			if shardId == 0 {
				return []*data.NodeData{
					{Address: "address1", ShardId: 0},
				}, nil
			} else {
				return []*data.NodeData{
					{Address: "address2", ShardId: core.MetachainShardId},
				}, nil
			}
		},

		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {

			var localMap map[string]interface{}
			if address == "address1" {
				localMap = map[string]interface{}{
					"metrics": map[string]interface{}{
						core.MetricCrossCheckBlockHeight: "meta 123",
					},
				}
			} else {
				localMap = map[string]interface{}{
					"metrics": map[string]interface{}{
						core.MetricNonce: 122,
					},
				}
			}

			genericResp := &data.GenericAPIResponse{Data: localMap}
			genRespBytes, _ := json.Marshal(genericResp)

			return 0, json.Unmarshal(genRespBytes, value)
		},
	})

	nonce, err := nodeStatusProc.GetLatestFullySynchronizedHyperblockNonce()
	require.NoError(t, err)
	require.Equal(t, uint64(122), nonce)
}

func TestNodeStatusProcessor_GetEconomicsDataMetricsGetRestEndPointErrorOnMetaShouldTryOnShard(t *testing.T) {
	t.Parallel()

	addressMeta := "address_meta"
	addressShard := "address_shard"
	shardNodeWasCalled := false

	localErr := errors.New("local error")
	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
			return []*data.NodeData{
				{Address: addressShard, ShardId: 0},
				{Address: addressMeta, ShardId: core.MetachainShardId},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			if address == addressMeta {
				return 0, localErr
			}
			if address == addressShard {
				shardNodeWasCalled = true
			}
			return 200, nil
		},
	})

	_, err := nodeStatusProc.GetEconomicsDataMetrics()
	require.NoError(t, err)
	require.True(t, shardNodeWasCalled)
}

func TestNodeStatusProcessor_GetEconomicsDataMetricsShouldWork(t *testing.T) {
	t.Parallel()

	addressMeta := "address_meta"
	expectedResponse := &data.GenericAPIResponse{
		Data: map[string]interface{}{
			"erd_total_supply": "12345",
		},
	}

	nodeStatusProc, _ := NewNodeStatusProcessor(&mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
			return []*data.NodeData{
				{Address: addressMeta, ShardId: core.MetachainShardId},
			}, nil
		},
		CallGetRestEndPointCalled: func(_ string, _ string, value interface{}) (int, error) {
			expectedResponseBytes, _ := json.Marshal(expectedResponse)
			return 200, json.Unmarshal(expectedResponseBytes, value)
		},
	})

	actualResponse, err := nodeStatusProc.GetEconomicsDataMetrics()
	require.NoError(t, err)
	require.Equal(t, *expectedResponse, *actualResponse)
}
