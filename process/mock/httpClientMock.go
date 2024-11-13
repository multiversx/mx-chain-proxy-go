package mock

import "net/http"

// HttpClientMock -
type HttpClientMock struct {
	DoCalled func(req *http.Request) (*http.Response, error)
}

// Do -
func (mock *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	if mock.DoCalled != nil {
		return mock.DoCalled(req)
	}
	return &http.Response{}, nil
}
