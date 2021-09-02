package disabled

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
)

// EpochStartNotifier represents a offline struct that implements the EpochStartNotifier interface
type EpochStartNotifier struct {
}

// RegisterNotifyHandler won't do anything as this is a offline component
func (e *EpochStartNotifier) RegisterNotifyHandler(_ core.EpochSubscriberHandler) {
}

// CurrentEpoch returns 0 as this is a offline component
func (e *EpochStartNotifier) CurrentEpoch() uint32 {
	return 0
}

// CheckEpoch won't do anything as this a offline component
func (e *EpochStartNotifier) CheckEpoch(_ data.HeaderHandler) {
}

// IsInterfaceNil returns true if there is no value under the interface
func (e *EpochStartNotifier) IsInterfaceNil() bool {
	return e == nil
}
