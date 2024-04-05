package mapCounters

import (
	"errors"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/availabilityCommon"
)

var (
	errInvalidAvailability           = errors.New("invalid data availability type")
	errNumNodesMustBeGreaterThanZero = errors.New("the number of nodes must be greater than 0")
)

// mapCountersHolder handles multiple counters map based on the data availability
type mapCountersHolder struct {
	countersMap map[data.ObserverDataAvailabilityType]*mapCounter
}

// NewMapCountersHolder populates the initial map and returns a new instance of mapCountersHolder
func NewMapCountersHolder() *mapCountersHolder {
	availabilityProvider := availabilityCommon.AvailabilityProvider{}
	dataAvailabilityTypes := availabilityProvider.GetAllAvailabilityTypes()

	countersMap := make(map[data.ObserverDataAvailabilityType]*mapCounter, len(dataAvailabilityTypes))
	for _, availability := range dataAvailabilityTypes {
		countersMap[availability] = newMapCounter()
	}

	return &mapCountersHolder{
		countersMap: countersMap,
	}
}

// ComputeShardPosition returns the shard position based on the availability and the shard
func (mch *mapCountersHolder) ComputeShardPosition(availability data.ObserverDataAvailabilityType, shardID uint32, numNodes uint32) (uint32, error) {
	if numNodes == 0 {
		return 0, errNumNodesMustBeGreaterThanZero
	}
	counterMap, exists := mch.countersMap[availability]
	if !exists {
		return 0, errInvalidAvailability
	}

	position := counterMap.computePositionForShard(shardID, numNodes)
	return position, nil
}

// ComputeAllNodesPosition returns the all nodes position based on the availability
func (mch *mapCountersHolder) ComputeAllNodesPosition(availability data.ObserverDataAvailabilityType, numNodes uint32) (uint32, error) {
	if numNodes == 0 {
		return 0, errNumNodesMustBeGreaterThanZero
	}
	counterMap, exists := mch.countersMap[availability]
	if !exists {
		return 0, errInvalidAvailability
	}

	position := counterMap.computePositionForAllNodes(numNodes)
	return position, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (mch *mapCountersHolder) IsInterfaceNil() bool {
	return mch == nil
}
