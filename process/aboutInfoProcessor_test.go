package process_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAboutInfoProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil base processor", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor(nil, "", "commitID")
		require.Nil(t, ap)
		require.Equal(t, process.ErrNilCoreProcessor, err)
	})

	t.Run("empty app version", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor(&mock.ProcessorStub{}, "", "commitID")
		require.Nil(t, ap)
		require.Equal(t, process.ErrEmptyAppVersionString, err)
	})

	t.Run("empty commit id", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor(&mock.ProcessorStub{}, "app version", "")
		require.Nil(t, ap)
		require.Equal(t, process.ErrEmptyCommitString, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor(&mock.ProcessorStub{}, "app version", "commitID")
		require.NotNil(t, ap)
		require.Nil(t, err)
	})
}

func TestAboutInfoProcessor_GetAboutInfo(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		appVersion := "appVersion"
		commit := "1221e3037839739dc0e119cc4c29c9f4d4101e57"

		ap, err := process.NewAboutProcessor(&mock.ProcessorStub{}, appVersion, commit)
		require.Nil(t, err)

		aboutInfo := &data.AboutInfo{
			AppVersion: appVersion,
			CommitID:   commit[:process.GetShortHashSize()],
		}

		expectedResp := &data.GenericAPIResponse{
			Data:  aboutInfo,
			Error: "",
			Code:  data.ReturnCodeSuccess,
		}

		resp := ap.GetAboutInfo()
		require.Equal(t, expectedResp, resp)
	})
}

func TestAboutProcessor_GetNodesVersions(t *testing.T) {
	t.Parallel()

	t.Run("one of the nodes responds with non-200, should error", func(t *testing.T) {
		t.Parallel()

		proc := &mock.ProcessorStub{
			GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{
						Address: "addr0",
					},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				return 37, nil
			},
		}

		ap, err := process.NewAboutProcessor(proc, "app", "hash")
		require.Nil(t, err)

		res, err := ap.GetNodesVersions()
		require.Empty(t, res)
		require.Equal(t, "invalid return code 37", err.Error())
	})

	t.Run("one of the nodes responds with error, should error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("request error")

		proc := &mock.ProcessorStub{
			GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{
						Address: "addr0",
					},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				return 200, expectedErr
			},
		}

		ap, err := process.NewAboutProcessor(proc, "app", "hash")
		require.Nil(t, err)

		res, err := ap.GetNodesVersions()
		require.Empty(t, res)
		require.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		proc := &mock.ProcessorStub{
			GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{
						Address: "addrSh0",
						ShardId: 0,
					},
					{
						Address: "addr0Sh1",
						ShardId: 1,
					},
					{
						Address: "addr1Sh1",
						ShardId: 1,
					},
					{
						Address: "addr0ShM",
						ShardId: core.MetachainShardId,
					},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				resp := &data.NodeVersionAPIResponse{}
				if strings.Contains(address, "Sh1") {
					resp.Data.Metrics.Version = "v1.37"
				} else {
					resp.Data.Metrics.Version = "v1.38"
				}

				respBytes, _ := json.Marshal(resp)
				return 200, json.Unmarshal(respBytes, value)
			},
		}

		ap, err := process.NewAboutProcessor(proc, "app", "hash")
		require.Nil(t, err)

		res, err := ap.GetNodesVersions()
		require.NoError(t, err)

		expectedResponse := &data.GenericAPIResponse{
			Data: data.NodesVersionProxyResponseData{
				Versions: map[uint32][]string{
					0:                     {"v1.38"},
					1:                     {"v1.37", "v1.37"},
					core.MetachainShardId: {"v1.38"},
				},
			},
			Code: data.ReturnCodeSuccess,
		}
		require.EqualValues(t, expectedResponse, res)
	})
}
