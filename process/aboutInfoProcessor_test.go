package process_test

import (
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/stretchr/testify/require"
)

func TestNewAboutInfoProcessor(t *testing.T) {
	t.Parallel()

	t.Run("empty app version", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor("", "commitID")
		require.Nil(t, ap)
		require.Equal(t, process.ErrEmptyAppVersionString, err)
	})

	t.Run("empty commit id", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor("app version", "")
		require.Nil(t, ap)
		require.Equal(t, process.ErrEmptyCommitString, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ap, err := process.NewAboutProcessor("app version", "commitID")
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

		ap, err := process.NewAboutProcessor(appVersion, commit)
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
