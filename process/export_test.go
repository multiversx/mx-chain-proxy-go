package process

import "time"

// SetDelayForCheckingNodesSyncState -
func (bp *BaseProcessor) SetDelayForCheckingNodesSyncState(delay time.Duration) {
	bp.delayForCheckingNodesSyncState = delay
}
