package process

import (
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	proxyData "github.com/multiversx/mx-chain-proxy-go/data"
)

// RelayedTxV2DataMarker -
const RelayedTxV2DataMarker = relayedTxV2DataMarker

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

// GetShortHashSize -
func GetShortHashSize() int {
	return shortHashSize
}

// ComputeTransactionStatus -
func (tp *TransactionProcessor) ComputeTransactionStatus(tx *transaction.ApiTransactionResult, withResults bool) *proxyData.ProcessStatusResponse {
	return tp.computeTransactionStatus(tx, withResults)
}

// CheckIfFailed -
func CheckIfFailed(logs []*transaction.ApiLogs) (bool, string) {
	return checkIfFailed(logs)
}
