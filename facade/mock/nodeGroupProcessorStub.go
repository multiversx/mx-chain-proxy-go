package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// NodeGroupProcessorStub represents a stub implementation of a NodeGroupProcessor
type NodeGroupProcessorStub struct {
	GetHeartbeatDataCalled                 func() (*data.HeartbeatResponse, error)
	IsOldStorageForTokenCalled             func(tokenID string, nonce uint64) (bool, error)
	GetWaitingEpochsLeftForPublicKeyCalled func(publicKey string) (*data.WaitingEpochsLeftApiResponse, error)
}

// IsOldStorageForToken -
func (hbps *NodeGroupProcessorStub) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	return hbps.IsOldStorageForTokenCalled(tokenID, nonce)
}

// GetHeartbeatData -
func (hbps *NodeGroupProcessorStub) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return hbps.GetHeartbeatDataCalled()
}

// GetWaitingEpochsLeftForPublicKey -
func (hbps *NodeGroupProcessorStub) GetWaitingEpochsLeftForPublicKey(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
	if hbps.GetWaitingEpochsLeftForPublicKeyCalled != nil {
		return hbps.GetWaitingEpochsLeftForPublicKeyCalled(publicKey)
	}
	return &data.WaitingEpochsLeftApiResponse{}, nil
}
