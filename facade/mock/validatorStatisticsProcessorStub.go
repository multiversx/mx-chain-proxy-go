package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// ValidatorStatisticsProcessorStub -
type ValidatorStatisticsProcessorStub struct {
	GetValidatorStatisticsCalled func() (*data.ValidatorStatisticsResponse, error)
}

// GetValidatorStatistics -
func (v *ValidatorStatisticsProcessorStub) GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error) {
	return v.GetValidatorStatisticsCalled()
}
