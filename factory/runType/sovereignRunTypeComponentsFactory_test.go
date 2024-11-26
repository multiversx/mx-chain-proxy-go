package runType

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSovereignRunTypeComponentsFactory_CreateAndClose(t *testing.T) {
	t.Parallel()

	srtcf := NewSovereignRunTypeComponentsFactory()
	require.NotNil(t, srtcf)

	srtc := srtcf.Create()
	require.NotNil(t, srtc)

	require.NoError(t, srtc.Close())
}
