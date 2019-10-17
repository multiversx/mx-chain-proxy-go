package mock

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type HeartbeatCacherMock struct {
	Data *data.HeartbeatResponse
}

func (hcm *HeartbeatCacherMock) Heartbeats() (*data.HeartbeatResponse, error) {
	if hcm.Data == nil {
		return nil, errors.New("nil Data")
	}

	return hcm.Data, nil
}

func (hcm *HeartbeatCacherMock) StoreHeartbeats(data *data.HeartbeatResponse) error {
	hcm.Data = data
	return nil
}

func (hcm *HeartbeatCacherMock) IsInterfaceNil() bool {
	return hcm == nil
}
