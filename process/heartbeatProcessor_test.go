package process_test

import (
	"errors"
	"fmt"
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

	providedAddressShard0 := "addr_0"
	providedHeartbeatsShard0 := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node0-1",
				PublicKey:       "pk0-1",
				ReceivedShardID: 0,
			},
			{
				NodeDisplayName: "node0-2",
				PublicKey:       "pk0-2",
				ReceivedShardID: 0,
			},
		},
	}
	providedAddressShard1 := "addr_1"
	providedHeartbeatsShard1 := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node1-1",
				PublicKey:       "pk1-1",
				ReceivedShardID: 1,
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
				ReceivedShardID: 0,
			},
			{
				// duplicate from shard 1
				NodeDisplayName: "node1-1",
				PublicKey:       "pk1-1",
				ReceivedShardID: 1,
			},
			{
				NodeDisplayName: "node2-1",
				PublicKey:       "pk2-1",
				ReceivedShardID: 2,
			},
		},
	}
	providedAddressShard3 := "addr_3"
	providedHeartbeatsShard3 := data.HeartbeatResponse{}

	providedShardIDs := []uint32{0, 1, 2, 3, 4}
	expectedErr := errors.New("expected error")
	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return providedShardIDs
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			assert.Contains(t, providedShardIDs, shardId)

			if shardId == 4 { // return no observers for this shard
				return nil, expectedErr
			}

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
			case providedAddressShard3:
				valResponse.Data = providedHeartbeatsShard3
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
			ReceivedShardID: 0,
		},
		{
			NodeDisplayName: "node0-2",
			PublicKey:       "pk0-2",
			ReceivedShardID: 0,
		},
		{
			NodeDisplayName: "node1-1",
			PublicKey:       "pk1-1",
			ReceivedShardID: 1,
		},
		{
			NodeDisplayName: "node2-1",
			PublicKey:       "pk2-1",
			ReceivedShardID: 2,
		},
	}

	assert.Equal(t, len(expectedSortedHeartbeats), len(res.Heartbeats))
	for idx := range res.Heartbeats {
		assert.Equal(t, expectedSortedHeartbeats[idx], res.Heartbeats[idx])
	}
}

func TestHeartbeatProcessor_GetHeartbeatDataShouldReturnDataFromApiBecauseCacheDataIsNil(t *testing.T) {
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
	hp, err := process.NewHeartbeatProcessor(
		&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{providedShardID}
			},
			GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
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

func TestHeartbeatProcessor_CacheShouldUpdate(t *testing.T) {
	t.Parallel()

	providedShardID := uint32(0)
	providedAddress := "addr"
	numOfTimesHttpWasCalled := int32(0)
	cacher := &mock.HeartbeatCacherMock{}
	hp, err := process.NewHeartbeatProcessor(&mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{providedShardID}
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
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

func TestHeartbeatProcessor_NoDataForAShardShouldNotUpdateCache(t *testing.T) {
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
	hp, err := process.NewHeartbeatProcessor(
		&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{providedShardID0, providedShardID1}
			},
			GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
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
	assert.Nil(t, err)
	assert.Nil(t, messages)
}
