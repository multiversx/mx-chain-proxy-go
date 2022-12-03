package process

import (
	"time"

	proxyData "github.com/ElrondNetwork/elrond-proxy-go/data"
)

// SetDelayForCheckingNodesSyncState -
func (bp *BaseProcessor) SetDelayForCheckingNodesSyncState(delay time.Duration) {
	bp.delayForCheckingNodesSyncState = delay
}

// SetNodeStatusFetcher -
func (bp *BaseProcessor) SetNodeStatusFetcher(fetcher func(url string) (*proxyData.NodeStatusAPIResponse, int, error)) {
	bp.nodeStatusFetcher = fetcher
}

// ComputeTokenStorageKey -
func ComputeTokenStorageKey(tokenID string, nonce uint64) string {
	return computeTokenStorageKey(tokenID, nonce)
}
