package process_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
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

func TestNewNodeGroupProcessor_NilProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewNodeGroupProcessor(nil, &mock.HeartbeatCacherMock{}, time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewNodeGroupProcessor_NilCacherShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{}, nil, time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrNilHeartbeatCacher, err)
}

func TestNewNodeGroupProcessor_InvalidCacheValidityDurationShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{}, &mock.HeartbeatCacherMock{}, -time.Second)

	assert.Nil(t, hp)
	assert.Equal(t, process.ErrInvalidCacheValidityDuration, err)
}

func TestNewNodeGroupProcessor_WithOkProcessorShouldErr(t *testing.T) {
	t.Parallel()

	hbp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{}, &mock.HeartbeatCacherMock{}, time.Second)

	assert.NotNil(t, hbp)
	assert.Nil(t, err)
}

func TestNodeGroupProcessor_GetHeartbeatDataWrongValuesShouldErr(t *testing.T) {
	t.Parallel()

	hp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{}, &mock.HeartbeatCacherMock{}, time.Second)
	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestNodeGroupProcessor_GetHeartbeatDataOkValuesShouldPass(t *testing.T) {
	t.Parallel()

	providedAddressShard0 := "addr_0"
	providedHeartbeatsShard0 := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node0-1",
				PublicKey:       "pk0-1",
				ComputedShardID: 0,
			},
			{
				NodeDisplayName: "node0-2",
				PublicKey:       "pk0-2",
				ComputedShardID: 0,
			},
		},
	}
	providedAddressShard1 := "addr_1"
	providedHeartbeatsShard1 := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node1-1",
				PublicKey:       "pk1-1",
				ComputedShardID: 1,
			},
		},
	}
	providedAddressShard2 := "addr_2"
	providedHeartbeatsShard2 := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				// duplicate from shard 0
				NodeDisplayName: "node0-1",
				PublicKey:       "pk0-1",
				ComputedShardID: 0,
			},
			{
				// duplicate from shard 1
				NodeDisplayName: "node1-1",
				PublicKey:       "pk1-1",
				ComputedShardID: 1,
			},
			{
				NodeDisplayName: "node2-1",
				PublicKey:       "pk2-1",
				ComputedShardID: 2,
			},
		},
	}

	providedShardIDs := []uint32{0, 1, 2}
	hp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return providedShardIDs
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			assert.Contains(t, providedShardIDs, shardId)

			var obs []*data.NodeData
			address := fmt.Sprintf("addr_%d", shardId)
			obs = append(obs, &data.NodeData{
				ShardId: shardId,
				Address: address,
			})

			return obs, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResponse := value.(*data.HeartbeatApiResponse)
			switch address {
			case providedAddressShard0:
				valResponse.Data = providedHeartbeatsShard0
			case providedAddressShard1:
				valResponse.Data = providedHeartbeatsShard1
			case providedAddressShard2:
				valResponse.Data = providedHeartbeatsShard2
			}
			return 0, nil
		},
	},
		&mock.HeartbeatCacherMock{},
		time.Second,
	)

	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()
	assert.NotNil(t, res)
	assert.Nil(t, err)

	expectedSortedHeartbeats := []data.PubKeyHeartbeat{
		{
			NodeDisplayName: "node0-1",
			PublicKey:       "pk0-1",
			ComputedShardID: 0,
		},
		{
			NodeDisplayName: "node0-2",
			PublicKey:       "pk0-2",
			ComputedShardID: 0,
		},
		{
			NodeDisplayName: "node1-1",
			PublicKey:       "pk1-1",
			ComputedShardID: 1,
		},
		{
			NodeDisplayName: "node2-1",
			PublicKey:       "pk2-1",
			ComputedShardID: 2,
		},
	}

	assert.Equal(t, len(expectedSortedHeartbeats), len(res.Heartbeats))
	for idx := range res.Heartbeats {
		assert.Equal(t, expectedSortedHeartbeats[idx], res.Heartbeats[idx])
	}
}

func TestNodeGroupProcessor_GetHeartbeatDataShouldReturnDataFromApiBecauseCacheDataIsNil(t *testing.T) {
	t.Parallel()

	providedHeartbeats := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node1",
				PublicKey:       "pk1",
			},
			{
				NodeDisplayName: "node2",
				PublicKey:       "pk2",
			},
		},
	}

	providedShardID := uint32(0)
	providedAddress := "addr"

	httpWasCalled := false
	// set nil hbts response in cache
	cacher := &mock.HeartbeatCacherMock{Data: nil}
	hp, err := process.NewNodeGroupProcessor(
		&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{providedShardID}
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				assert.Equal(t, providedShardID, shardId)
				var obs []*data.NodeData
				obs = append(obs, &data.NodeData{
					ShardId: providedShardID,
					Address: providedAddress,
				})
				return obs, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				assert.Equal(t, providedAddress, address)
				valResponse := value.(*data.HeartbeatApiResponse)
				valResponse.Data = providedHeartbeats
				httpWasCalled = true
				return 0, nil
			},
		},
		cacher,
		time.Second,
	)
	assert.Nil(t, err)

	_, err = hp.GetHeartbeatData()
	assert.Nil(t, err)
	assert.True(t, httpWasCalled)
}

func TestNodeGroupProcessor_GetHeartbeatDataShouldReturnDataFromApiBecauseCacheDataIsNil_MultipleMessagesForSamePK(t *testing.T) {
	t.Parallel()

	providedHeartbeatsFirstCall := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				PublicKey:       "pk1",
				IsActive:        false,
				ComputedShardID: 0,
				ReceivedShardID: 0,
			},
			{
				PublicKey:       "pk2",
				IsActive:        true,
				ComputedShardID: 0,
				ReceivedShardID: 0,
			},
			{
				PublicKey:       "pk4", // node after shuffle out, moved from shard 0 to 1
				IsActive:        false,
				ComputedShardID: 0,
				ReceivedShardID: 0,
			},
		},
	}
	providedHeartbeatsSecondCall := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				PublicKey:       "pk1", // same as on first call
				IsActive:        true,
				ComputedShardID: 1,
				ReceivedShardID: 1,
			},
			{
				PublicKey:       "pk3",
				IsActive:        true,
				ComputedShardID: 1,
				ReceivedShardID: 1,
			},
			{
				PublicKey:       "pk4",
				IsActive:        true,
				ComputedShardID: 1,
				ReceivedShardID: 0,
			},
		},
	}
	expectedHeartbeats := &data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				PublicKey:       "pk1",
				IsActive:        true,
				ComputedShardID: 1,
				ReceivedShardID: 1,
			},
			{
				PublicKey:       "pk2",
				IsActive:        true,
				ComputedShardID: 0,
				ReceivedShardID: 0,
			},
			{
				PublicKey:       "pk3",
				IsActive:        true,
				ComputedShardID: 1,
				ReceivedShardID: 1,
			},
			{
				PublicKey:       "pk4",
				IsActive:        true,
				ComputedShardID: 1,
				ReceivedShardID: 0,
			},
		},
	}

	providedShardID0, providedShardID1 := uint32(0), uint32(1)

	providedAddress0, providedAddress1 := "addr0", "addr1"

	httpWasCalled := false
	// set nil hbts response in cache
	cacher := &mock.HeartbeatCacherMock{Data: nil}
	counter := 0
	hp, err := process.NewNodeGroupProcessor(
		&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{providedShardID0, providedShardID1}
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				var obs []*data.NodeData
				switch counter {
				case 0:
					obs = append(obs, &data.NodeData{
						ShardId: providedShardID0,
						Address: providedAddress0,
					})
				case 1:
					obs = append(obs, &data.NodeData{
						ShardId: providedShardID1,
						Address: providedAddress1,
					})
				default:
					assert.Fail(t, "only 2 shard provided")
				}

				return obs, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				valResponse := value.(*data.HeartbeatApiResponse)
				switch counter {
				case 0:
					valResponse.Data = providedHeartbeatsFirstCall
				case 1:
					valResponse.Data = providedHeartbeatsSecondCall
				default:
					assert.Fail(t, "only 2 shard provided")
				}
				httpWasCalled = true
				counter++
				return 0, nil
			},
		},
		cacher,
		time.Second,
	)
	assert.Nil(t, err)

	heartbeats, err := hp.GetHeartbeatData()
	assert.Nil(t, err)
	assert.True(t, httpWasCalled)
	assert.Equal(t, expectedHeartbeats, heartbeats)
}

func TestNodeGroupProcessor_GetHeartbeatDataShouldReturnDataFromCacher(t *testing.T) {
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
	hp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{}, cacher, time.Millisecond)
	assert.Nil(t, err)

	res, err := hp.GetHeartbeatData()

	assert.Nil(t, err)
	assert.Equal(t, *res, hbtsResp)
}

func TestNodeGroupProcessor_CacheShouldUpdate(t *testing.T) {
	t.Parallel()

	providedShardID := uint32(0)
	providedAddress := "addr"
	numOfTimesHttpWasCalled := int32(0)
	cacher := &mock.HeartbeatCacherMock{}
	hp, err := process.NewNodeGroupProcessor(&mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{providedShardID}
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			assert.Equal(t, providedShardID, shardId)
			var obs []*data.NodeData
			obs = append(obs, &data.NodeData{
				ShardId: providedShardID,
				Address: providedAddress,
			})
			return obs, nil
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

func TestNodeGroupProcessor_NoDataForAShardShouldNotUpdateCache(t *testing.T) {
	t.Parallel()

	providedHeartbeatsShard0 := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node01",
				PublicKey:       "pk01",
				ReceivedShardID: 0,
			},
			{
				NodeDisplayName: "node02",
				PublicKey:       "pk02",
				ReceivedShardID: 1,
			},
		},
	}

	providedShardID0, providedShardID1 := uint32(0), uint32(1)
	providedAddressShard0, providedAddress1Shard1, providedAddress2Shard1 := "addr0_1", "addr1_1", "addr1_2"

	expectedErr := errors.New("expected error")
	// set nil hbts response in cache
	cacher := &mock.HeartbeatCacherMock{Data: nil}
	hp, err := process.NewNodeGroupProcessor(
		&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{providedShardID0, providedShardID1}
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				var obs []*data.NodeData
				if shardId == providedShardID0 {
					obs = append(obs, &data.NodeData{
						ShardId: providedShardID0,
						Address: providedAddressShard0,
					})
					return obs, nil
				}

				obs = append(obs, &data.NodeData{
					ShardId: providedShardID1,
					Address: providedAddress1Shard1,
				}, &data.NodeData{
					ShardId: providedShardID1,
					Address: providedAddress2Shard1,
				})
				return obs, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				// Shard 1 observer 1 returns error
				if address == providedAddress1Shard1 {
					return 0, expectedErr
				}

				// Shard 1 observer 2 returns empty messages
				if address == providedAddress2Shard1 {
					return 0, nil
				}

				// Shard 0 returns valid data
				valResponse := value.(*data.HeartbeatApiResponse)
				valResponse.Data = providedHeartbeatsShard0
				return 0, nil
			},
		},
		cacher,
		time.Second,
	)
	assert.Nil(t, err)

	messages, err := hp.GetHeartbeatData()
	assert.Equal(t, process.ErrHeartbeatNotAvailable, err)
	assert.Nil(t, messages)
}

func TestNodeGroupProcessor_IsOldStorageForToken(t *testing.T) {
	t.Parallel()

	t.Run("all observers fail, should return error", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{Address: "addr0", ShardId: 0},
						{Address: "addr1", ShardId: 1},
					}, nil
				},
				CallGetRestEndPointCalled: func(_ string, _ string, _ interface{}) (int, error) {
					return 0, errors.New("error")
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		_, err := proc.IsOldStorageForToken("token", 37)
		require.True(t, errors.Is(err, process.ErrSendingRequest))
	})

	t.Run("some observers fail, should return error", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{Address: "addr0", ShardId: 0},
						{Address: "addr1", ShardId: 1},
					}, nil
				},
				CallGetRestEndPointCalled: func(address string, _ string, _ interface{}) (int, error) {
					if address == "addr1" {
						return 0, nil
					}
					return 0, errors.New("error")
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		_, err := proc.IsOldStorageForToken("token", 37)
		require.True(t, errors.Is(err, process.ErrSendingRequest))
	})

	t.Run("should work and return false", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{Address: "addr0", ShardId: 0},
						{Address: "addr1", ShardId: 1},
					}, nil
				},
				CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
					valResponse := value.(*data.AccountKeyValueResponse)
					valResponse.Data.Value = "test"
					return 0, nil
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		isOldStorage, err := proc.IsOldStorageForToken("token", 37)
		require.False(t, isOldStorage)
		require.NoError(t, err)
	})

	t.Run("should work and return true", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{Address: "addr0", ShardId: 0},
						{Address: "addr1", ShardId: 1},
					}, nil
				},
				CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
					valResponse := value.(*data.AccountKeyValueResponse)
					valResponse.Data.Value = ""
					return 0, nil
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		isOldStorage, err := proc.IsOldStorageForToken("token", 37)
		require.True(t, isOldStorage)
		require.NoError(t, err)
	})
}

func TestComputeTokenStorageKey(t *testing.T) {
	t.Parallel()

	require.Equal(t, "454c524f4e4465736474746f6b656e25", process.ComputeTokenStorageKey("token", 37))
	require.Equal(t, "454c524f4e4465736474455254574f2d3364313934340284", process.ComputeTokenStorageKey("ERTWO-3d1944", 644))

	testTokenID, testNonce := "TESTTKN", uint64(89)
	expectedKey := append(append([]byte(core.ProtectedKeyPrefix+"esdt"), []byte(testTokenID)...), big.NewInt(int64(testNonce)).Bytes()...)
	require.Equal(t, hex.EncodeToString(expectedKey), process.ComputeTokenStorageKey(testTokenID, testNonce))
}

func TestNodeGroupProcessor_GetWaitingEpochsLeftForPublicKey(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("empty pub key should error", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{},
			&mock.HeartbeatCacherMock{},
			10,
		)

		response, err := proc.GetWaitingEpochsLeftForPublicKey("")
		require.Nil(t, response)
		require.Equal(t, process.ErrEmptyPubKey, err)
	})
	t.Run("GetAllObservers returns error should error", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return nil, expectedErr
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		response, err := proc.GetWaitingEpochsLeftForPublicKey("key")
		require.Nil(t, response)
		require.Equal(t, expectedErr, err)
	})
	t.Run("all observers return error should error", func(t *testing.T) {
		t.Parallel()

		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{Address: "addr0", ShardId: 0},
						{Address: "addr1", ShardId: 1},
					}, nil
				},
				CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
					return 0, expectedErr
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		response, err := proc.GetWaitingEpochsLeftForPublicKey("key")
		require.Nil(t, response)
		require.True(t, errors.Is(err, process.ErrSendingRequest))
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		providedEpochsLeft := uint32(10)
		proc, _ := process.NewNodeGroupProcessor(
			&mock.ProcessorStub{
				GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{Address: "addr0", ShardId: 0},
						{Address: "addr1", ShardId: 1},
					}, nil
				},
				CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
					valResponse := value.(*data.WaitingEpochsLeftApiResponse)
					valResponse.Data.EpochsLeft = providedEpochsLeft
					return 0, nil
				},
			},
			&mock.HeartbeatCacherMock{},
			10,
		)

		response, err := proc.GetWaitingEpochsLeftForPublicKey("key")
		require.NoError(t, err)
		require.Equal(t, providedEpochsLeft, response.Data.EpochsLeft)
	})
}
