package process

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewESDTSupplyProcessor(t *testing.T) {
	t.Parallel()

	_, err := NewESDTSupplyProcessor(nil, &mock.SCQueryServiceStub{})
	require.Equal(t, ErrNilCoreProcessor, err)

	_, err = NewESDTSupplyProcessor(&mock.ProcessorStub{}, nil)
	require.Equal(t, ErrNilSCQueryService, err)
}

func TestEsdtSupplyProcessor_GetESDTSupplyFungible(t *testing.T) {
	t.Parallel()

	baseProc := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1, core.MetachainShardId}
		},
		GetObserversCalled: func(shardID uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			switch address {
			case "shard-0":
				valResp := value.(*data.ESDTSupplyResponse)
				valResp.Data.Supply = "1000"
				return 200, nil
			case "shard-1":
				valResp := value.(*data.ESDTSupplyResponse)
				valResp.Data.Supply = "3000"
				return 200, nil
			}
			return 0, nil
		},
	}
	scQueryProc := &mock.SCQueryServiceStub{
		ExecuteQueryCalled: func(query *data.SCQuery) (*vm.VMOutputApi, error) {
			return &vm.VMOutputApi{
				ReturnData: [][]byte{nil, nil, nil, []byte("500")},
			}, nil
		},
	}
	esdtProc, err := NewESDTSupplyProcessor(baseProc, scQueryProc)
	require.Nil(t, err)

	supplyRes, err := esdtProc.GetESDTSupply("TOKEN-ABCD")
	require.Nil(t, err)
	require.Equal(t, "4500", supplyRes.Data.Supply)
}

func TestEsdtSupplyProcessor_GetESDTSupplyNonFungible(t *testing.T) {
	t.Parallel()

	called := false
	baseProc := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0, 1, core.MetachainShardId}
		},
		GetObserversCalled: func(shardID uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
				{
					ShardId: shardID,
					Address: fmt.Sprintf("shard-%d", shardID),
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			switch address {
			case "shard-0":
				if !called {
					called = true
					return 400, errors.New("local err")
				}
				valResp := value.(*data.ESDTSupplyResponse)
				valResp.Data.Supply = "-1000"
				return 200, nil
			case "shard-1":
				valResp := value.(*data.ESDTSupplyResponse)
				valResp.Data.Supply = "3000"
				return 200, nil
			}
			return 0, nil
		},
	}
	scQueryProc := &mock.SCQueryServiceStub{}
	esdtProc, err := NewESDTSupplyProcessor(baseProc, scQueryProc)
	require.Nil(t, err)

	supplyRes, err := esdtProc.GetESDTSupply("SEMI-ABCD-0A")
	require.Nil(t, err)
	require.Equal(t, "2000", supplyRes.Data.Supply)
}
