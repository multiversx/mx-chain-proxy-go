package availabilityCommon

import (
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// AvailabilityProvider is a stateless component that aims to group common operations regarding observers' data availability
type AvailabilityProvider struct {
}

// AvailabilityForAccountQueryOptions returns the availability needed for the provided query options
func (ap *AvailabilityProvider) AvailabilityForAccountQueryOptions(options common.AccountQueryOptions) data.ObserverDataAvailabilityType {
	availability := data.AvailabilityRecent
	if options.AreHistoricalCoordinatesSet() {
		availability = data.AvailabilityAll
	}
	return availability
}

// AvailabilityForVmQuery returns the availability needed for the provided query options
func (ap *AvailabilityProvider) AvailabilityForVmQuery(query *data.SCQuery) data.ObserverDataAvailabilityType {
	availability := data.AvailabilityRecent
	if query.BlockNonce.HasValue || len(query.BlockHash) > 0 {
		availability = data.AvailabilityAll
	}
	return availability
}

// IsNodeValid returns true if the provided node is valid based on the availability
func (ap *AvailabilityProvider) IsNodeValid(node *data.NodeData, availability data.ObserverDataAvailabilityType) bool {
	isInvalidSnapshotlessNode := availability == data.AvailabilityRecent && !node.IsSnapshotless
	isInvalidRegularNode := availability == data.AvailabilityAll && node.IsSnapshotless
	isInvalidNode := isInvalidSnapshotlessNode || isInvalidRegularNode
	return !isInvalidNode
}

// GetDescriptionForAvailability returns a short description string about the provided availability
func (ap *AvailabilityProvider) GetDescriptionForAvailability(availability data.ObserverDataAvailabilityType) string {
	switch availability {
	case data.AvailabilityAll:
		return "regular nodes"
	case data.AvailabilityRecent:
		return "snapshotless nodes"
	default:
		return "N/A"
	}
}

// GetAllAvailabilityTypes returns all data availability types
func (ap *AvailabilityProvider) GetAllAvailabilityTypes() []data.ObserverDataAvailabilityType {
	return []data.ObserverDataAvailabilityType{data.AvailabilityAll, data.AvailabilityRecent}
}
