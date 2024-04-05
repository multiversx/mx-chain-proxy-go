package process

import (
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestValidatorStatisticsProcessor_GetAuctionList(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		node := &data.NodeData{
			Address: "addr",
			ShardId: core.MetachainShardId,
		}
		expectedResp := &data.AuctionListAPIResponse{
			Data: data.AuctionListResponse{
				AuctionListValidators: []*data.AuctionListValidatorAPIResponse{
					{
						Owner:          "owner",
						NumStakedNodes: 4,
						TotalTopUp:     "100",
						TopUpPerNode:   "50",
						QualifiedTopUp: "50",
					},
				},
			},
			Error: "",
			Code:  "ok",
		}
		expectedRespMarshalled, err := json.Marshal(expectedResp)
		require.Nil(t, err)

		processor := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, core.MetachainShardId, shardId)

				return []*data.NodeData{node}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				require.Equal(t, node.Address, address)
				require.Equal(t, auctionListPath, path)

				err = json.Unmarshal(expectedRespMarshalled, value)
				require.Nil(t, err)
				return 0, nil
			},
		}
		vsp, _ := NewValidatorStatisticsProcessor(processor, &mock.ValStatsCacherMock{}, time.Second)
		resp, err := vsp.GetAuctionList()
		require.Nil(t, err)
		require.Equal(t, expectedResp.Data, *resp)
	})

	t.Run("get observers failed, should return error", func(t *testing.T) {
		t.Parallel()

		errGetObservers := errors.New("err getting observers")
		callGetRestEndPointCalledCt := int32(0)

		processor := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, core.MetachainShardId, shardId)
				return nil, errGetObservers
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				atomic.AddInt32(&callGetRestEndPointCalledCt, 1)

				return 0, nil
			},
		}
		vsp, _ := NewValidatorStatisticsProcessor(processor, &mock.ValStatsCacherMock{}, time.Second)

		resp, err := vsp.GetAuctionList()
		require.Equal(t, errGetObservers, err)
		require.Nil(t, resp)
		require.Equal(t, int32(0), callGetRestEndPointCalledCt)
	})

	t.Run("could not get auction list from observer, should return error", func(t *testing.T) {
		t.Parallel()

		node := &data.NodeData{
			Address: "addr",
			ShardId: core.MetachainShardId,
		}

		errCallEndpoint := errors.New("error call endpoint")
		processor := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, core.MetachainShardId, shardId)

				return []*data.NodeData{node}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				require.Equal(t, node.Address, address)
				require.Equal(t, auctionListPath, path)

				return 0, errCallEndpoint
			},
		}
		vsp, _ := NewValidatorStatisticsProcessor(processor, &mock.ValStatsCacherMock{}, time.Second)

		resp, err := vsp.GetAuctionList()
		require.Equal(t, ErrAuctionListNotAvailable, err)
		require.Nil(t, resp)
	})
}
