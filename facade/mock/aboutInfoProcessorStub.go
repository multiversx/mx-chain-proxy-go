package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// AboutInfoProcessorStub -
type AboutInfoProcessorStub struct {
	GetAboutInfoCalled     func() *data.GenericAPIResponse
	GetNodesVersionsCalled func() (*data.GenericAPIResponse, error)
}

// GetAboutInfo -
func (stub *AboutInfoProcessorStub) GetAboutInfo() *data.GenericAPIResponse {
	if stub.GetAboutInfoCalled != nil {
		return stub.GetAboutInfoCalled()
	}

	return nil
}

// GetNodesVersions -
func (stub *AboutInfoProcessorStub) GetNodesVersions() (*data.GenericAPIResponse, error) {
	if stub.GetNodesVersionsCalled != nil {
		return stub.GetNodesVersionsCalled()
	}

	return nil, nil
}
