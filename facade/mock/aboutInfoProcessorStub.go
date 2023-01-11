package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// AboutInfoProcessorStub -
type AboutInfoProcessorStub struct {
	GetAboutInfoCalled func() *data.GenericAPIResponse
}

// GetAboutInfo -
func (stub *AboutInfoProcessorStub) GetAboutInfo() *data.GenericAPIResponse {
	if stub.GetAboutInfoCalled != nil {
		return stub.GetAboutInfoCalled()
	}

	return nil
}
