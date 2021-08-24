package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// ESDTSuppliesProcessorStub -
type ESDTSuppliesProcessorStub struct {
	GetESDTSupplyCalled func(token string) (*data.GenericAPIResponse, error)
}

// GetESDTSupply -
func (e *ESDTSuppliesProcessorStub) GetESDTSupply(token string) (*data.GenericAPIResponse, error) {
	if e.GetESDTSupplyCalled != nil {
		return e.GetESDTSupplyCalled(token)
	}

	return nil, nil
}
