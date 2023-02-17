package disabled

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data"
)

// EpochStartNotifier represents a disabled struct that implements the EpochStartNotifier interface
type EpochStartNotifier struct {
}

// RegisterNotifyHandler won't do anything as this is a disabled component
func (e *EpochStartNotifier) RegisterNotifyHandler(_ core.EpochSubscriberHandler) {
}

// CurrentEpoch returns 0 as this is a disabled component
func (e *EpochStartNotifier) CurrentEpoch() uint32 {
	return 0
}

// CheckEpoch won't do anything as this a disabled component
func (e *EpochStartNotifier) CheckEpoch(_ data.HeaderHandler) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (e *EpochStartNotifier) IsInterfaceNil() bool {
	return e == nil
}
