package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// ESDTSuppliesProcessorStub -
type ESDTSuppliesProcessorStub struct {
	GetESDTSupplyCalled func(token string) (*data.ESDTSupplyResponse, error)
}

// GetESDTSupply -
func (e *ESDTSuppliesProcessorStub) GetESDTSupply(token string) (*data.ESDTSupplyResponse, error) {
	if e.GetESDTSupplyCalled != nil {
		return e.GetESDTSupplyCalled(token)
	}

	return nil, nil
}
