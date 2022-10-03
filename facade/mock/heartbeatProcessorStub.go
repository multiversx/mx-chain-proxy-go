package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// HeartbeatProcessorStub represents a stub implementation of a HeartbeatProcessor
type HeartbeatProcessorStub struct {
	GetHeartbeatDataCalled     func() (*data.HeartbeatResponse, error)
	IsOldStorageForTokenCalled func(tokenID string, nonce uint64) (bool, error)
}

// IsOldStorageForToken -
func (hbps *HeartbeatProcessorStub) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	return hbps.IsOldStorageForTokenCalled(tokenID, nonce)
}

// GetHeartbeatData will call the handler func
func (hbps *HeartbeatProcessorStub) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return hbps.GetHeartbeatDataCalled()
}
