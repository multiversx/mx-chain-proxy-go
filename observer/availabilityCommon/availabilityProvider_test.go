package availabilityCommon

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestAvailabilityForAccountQueryOptions(t *testing.T) {
	t.Parallel()

	ap := &AvailabilityProvider{}

	// Test with historical coordinates set
	options := common.AccountQueryOptions{BlockHash: []byte("hash")}
	require.Equal(t, data.AvailabilityAll, ap.AvailabilityForAccountQueryOptions(options))

	// Test without historical coordinates set
	options = common.AccountQueryOptions{}
	require.Equal(t, data.AvailabilityRecent, ap.AvailabilityForAccountQueryOptions(options))
}

func TestAvailabilityForVmQuery(t *testing.T) {
	t.Parallel()

	ap := &AvailabilityProvider{}

	// Test with BlockNonce set
	query := &data.SCQuery{BlockNonce: core.OptionalUint64{HasValue: true, Value: 37}}
	require.Equal(t, data.AvailabilityAll, ap.AvailabilityForVmQuery(query))

	// Test without BlockNonce set but with BlockHash
	query = &data.SCQuery{BlockHash: []byte("hash")}
	require.Equal(t, data.AvailabilityAll, ap.AvailabilityForVmQuery(query))

	// Test without BlockNonce and BlockHash
	query = &data.SCQuery{}
	require.Equal(t, data.AvailabilityRecent, ap.AvailabilityForVmQuery(query))
}

func TestIsNodeValid(t *testing.T) {
	t.Parallel()

	ap := &AvailabilityProvider{}

	// Test with AvailabilityRecent and snapshotless node
	node := &data.NodeData{IsSnapshotless: true}
	require.True(t, ap.IsNodeValid(node, data.AvailabilityRecent))

	// Test with AvailabilityRecent and regular node
	node = &data.NodeData{}
	require.False(t, ap.IsNodeValid(node, data.AvailabilityRecent))

	// Test with AvailabilityAll and regular node
	node = &data.NodeData{}
	require.True(t, ap.IsNodeValid(node, data.AvailabilityAll))

	// Test with AvailabilityAll and Snapshotless node
	node = &data.NodeData{IsSnapshotless: true}
	require.False(t, ap.IsNodeValid(node, data.AvailabilityAll))
}

func TestGetDescriptionForAvailability(t *testing.T) {
	t.Parallel()

	ap := &AvailabilityProvider{}

	require.Equal(t, "regular nodes", ap.GetDescriptionForAvailability(data.AvailabilityAll))
	require.Equal(t, "snapshotless nodes", ap.GetDescriptionForAvailability(data.AvailabilityRecent))
	require.Equal(t, "N/A", ap.GetDescriptionForAvailability("invalid")) // Invalid value
}

func TestAvailabilityProvider_GetAllAvailabilityTypes(t *testing.T) {
	t.Parallel()

	ap := &AvailabilityProvider{}
	require.Equal(t, []data.ObserverDataAvailabilityType{data.AvailabilityAll, data.AvailabilityRecent}, ap.GetAllAvailabilityTypes())
}
