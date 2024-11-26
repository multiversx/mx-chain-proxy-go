package runType

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRunTypeComponentsFactory(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		rtc := NewRunTypeComponentsFactory()
		require.NotNil(t, rtc)
	})
}

func TestRunTypeComponentsFactory_Create(t *testing.T) {
	t.Parallel()

	rtcf := NewRunTypeComponentsFactory()
	require.NotNil(t, rtcf)

	rtc := rtcf.Create()
	require.NotNil(t, rtc)
}

func TestRunTypeComponentsFactory_Close(t *testing.T) {
	t.Parallel()

	rtcf := NewRunTypeComponentsFactory()
	require.NotNil(t, rtcf)

	rtc := rtcf.Create()
	require.NotNil(t, rtc)

	require.NoError(t, rtc.Close())
}
