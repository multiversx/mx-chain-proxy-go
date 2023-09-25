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

// MapCountersHolder handles multiple counters map based on the data availability
type MapCountersHolder struct {
	countersMap map[data.ObserverDataAvailabilityType]*mapCounter
}

// NewMapCountersHolder populates the initial map and returns a new instance of MapCountersHolder
func NewMapCountersHolder() *MapCountersHolder {
	availabilityProvider := availabilityCommon.AvailabilityProvider{}
	dataAvailabilityTypes := availabilityProvider.GetAllAvailabilityTypes()

	countersMap := make(map[data.ObserverDataAvailabilityType]*mapCounter)
	for _, availability := range dataAvailabilityTypes {
		countersMap[availability] = newMapCounter()
	}

	return &MapCountersHolder{
		countersMap: countersMap,
	}
}

// ComputeShardPosition returns the shard position based on the availability and the shard
func (mch *MapCountersHolder) ComputeShardPosition(availability data.ObserverDataAvailabilityType, shardID uint32, numNodes uint32) (uint32, error) {
	counterMap, exists := mch.countersMap[availability]
	if !exists {
		return 0, errInvalidAvailability
	}

	if numNodes == 0 {
		return 0, errNumNodesMustBeGreaterThanZero
	}

	position := counterMap.computePositionForShard(shardID, numNodes)
	return position, nil
}

// ComputeAllNodesPosition returns the all nodes position based on the availability
func (mch *MapCountersHolder) ComputeAllNodesPosition(availability data.ObserverDataAvailabilityType, numNodes uint32) (uint32, error) {
	counterMap, exists := mch.countersMap[availability]
	if !exists {
		return 0, errInvalidAvailability
	}

	if numNodes == 0 {
		return 0, errNumNodesMustBeGreaterThanZero
	}

	position := counterMap.computePositionForAllNodes(numNodes)
	return position, nil
}
