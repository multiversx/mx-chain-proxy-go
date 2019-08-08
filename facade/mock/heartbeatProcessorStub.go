package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// HeartbeatProcessorStub represents a stub implementation of a HeartbeatProcessor
type HeartbeatProcessorStub struct {
	GetHeartbeatDataCalled func() (*data.HeartbeatResponse, error)
}

// GetHeartbeatData will call the handler func
func (hbps *HeartbeatProcessorStub) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return hbps.GetHeartbeatDataCalled()
}
