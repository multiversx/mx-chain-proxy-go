package process

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// HeartBeatPath represents the path where an observer exposes his heartbeat status
const HeartBeatPath = "/node/heartbeatstatus"

// HeartbeatProcessor is able to process transaction requests
type HeartbeatProcessor struct {
	proc Processor
}

// NewHeartbeatProcessor creates a new instance of TransactionProcessor
func NewHeartbeatProcessor(proc Processor) (*HeartbeatProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &HeartbeatProcessor{
		proc: proc,
	}, nil
}

// GetHeartbeatData will simply forward the heartbeat status from an observer
func (hbp *HeartbeatProcessor) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	observer, err := hbp.proc.GetFirstAvailableObserver()
	if err != nil {
		return nil, err
	}

	var heartbeatResponse data.HeartbeatResponse
	err = hbp.proc.CallGetRestEndPoint(observer.Address, HeartBeatPath, &heartbeatResponse)
	if err != nil {
		return nil, err
	}

	return &heartbeatResponse, nil
}
