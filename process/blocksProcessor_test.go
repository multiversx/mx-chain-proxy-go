package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewBlocksProcessor_NilProcessor_ExpectError(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlocksProcessor(nil)

	require.Nil(t, bp)
	require.Equal(t, err, process.ErrNilCoreProcessor)
}

func TestBlocksProcessor_GetBlocksByRound_InvalidObservers_ExpectError(t *testing.T) {
	t.Parallel()

	err := errors.New("err observers")
	proc := &mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, err
		},
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1, 2}
		},
	}

	bp, _ := process.NewBlocksProcessor(proc)

	ret, actualErr := bp.GetBlocksByRound(0, common.BlockQueryOptions{})

	require.Equal(t, err, actualErr)
	require.Equal(t, (*data.BlocksApiResponse)(nil), ret)
}

func TestBlocksProcessor_GetBlocksByRound_InvalidCallGetRestEndPoint_ExpectZeroFetchedBlocks(t *testing.T) {
	t.Parallel()

	err := errors.New("err call get")
	proc := &mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			switch shardId {
			case 0:
				return []*data.NodeData{
					{
						ShardId: 0,
						Address: "erd1a",
					},
				}, nil
			case 1:
				return []*data.NodeData{
					{
						ShardId: 1,
						Address: "erd1b",
					},
				}, nil
			default:
				return nil, nil
			}
		},
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1}
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 404, err
		},
	}

	bp, _ := process.NewBlocksProcessor(proc)

	ret, actualErr := bp.GetBlocksByRound(0, common.BlockQueryOptions{})
	expectedRet := &data.BlocksApiResponse{
		Data: data.BlocksApiResponsePayload{
			Blocks: make([]*api.Block, 0, 2),
		},
	}
	require.Equal(t, nil, actualErr)
	require.Equal(t, expectedRet, ret)
}

func TestBlocksProcessor_GetBlocksByRound_TwoBlocks_ThreeObservers_OneObserverGetEndpointInvalid_ExpectTwoFetchedBlocks(t *testing.T) {
	t.Parallel()

	block1 := api.Block{
		Round: 111,
		Nonce: 222,
	}
	block2 := api.Block{
		Round: 333,
		Nonce: 444,
	}

	proc := &mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			switch shardId {
			case 0:
				return []*data.NodeData{
					{
						ShardId: 0,
						Address: "erd1a",
					},
					{
						ShardId: 0,
						Address: "erd1b",
					},
				}, nil
			case 1:
				return []*data.NodeData{
					{
						ShardId: 1,
						Address: "erd1c",
					},
				}, nil
			default:
				return nil, nil
			}
		},
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1}
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			if address == "erd1b" {
				response := value.(*data.BlockApiResponse)
				response.Data = data.BlockApiResponsePayload{Block: block1}

				return 200, nil
			}
			if address == "erd1c" {
				response := value.(*data.BlockApiResponse)
				response.Data = data.BlockApiResponsePayload{Block: block2}

				return 200, nil
			}
			return 404, errors.New("error call get")
		},
	}

	bp, _ := process.NewBlocksProcessor(proc)
	ret, err := bp.GetBlocksByRound(0, common.BlockQueryOptions{WithTransactions: true})

	expectedApiResp := &data.BlocksApiResponse{
		Data: data.BlocksApiResponsePayload{
			Blocks: []*api.Block{&block1, &block2},
		},
	}
	require.Nil(t, err)
	require.Equal(t, expectedApiResp, ret)
}
