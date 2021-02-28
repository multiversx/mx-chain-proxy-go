package factory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const dirLocation = "testdata"

func TestNewApiConfigParser_IncorrectPath(t *testing.T) {
	acp, err := NewApiConfigParser("wrong path")
	require.Nil(t, acp)
	require.Error(t, err)
}

func TestNewApiConfigParser_FileInsteadOfDirectory(t *testing.T) {
	acp, err := NewApiConfigParser(fmt.Sprintf("%s/%s", dirLocation, "vx_x.toml"))
	require.Error(t, err)
	require.Nil(t, acp)
}

func TestNewApiConfigParser(t *testing.T) {
	acp, err := NewApiConfigParser(dirLocation)
	require.NoError(t, err)
	require.NotNil(t, acp)
}

func TestApiConfigParser_GetConfigForVersionInvalidVersion(t *testing.T) {
	acp, _ := NewApiConfigParser(dirLocation)

	res, err := acp.GetConfigForVersion("wrong version")
	require.Nil(t, res)
	require.Error(t, err)
}

func TestApiConfigParser_GetConfigForVersion(t *testing.T) {
	acp, _ := NewApiConfigParser(dirLocation)

	res, err := acp.GetConfigForVersion("vx_x")
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, 1, len(res.APIPackages))

	endpointConfig, ok := res.APIPackages["testendpoint"]
	require.True(t, ok)
	require.Equal(t, 3, len(endpointConfig.Routes))
}
