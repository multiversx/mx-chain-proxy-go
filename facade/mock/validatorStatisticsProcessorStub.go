package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// ValidatorStatisticsProcessorStub -
type ValidatorStatisticsProcessorStub struct {
	GetValidatorStatisticsCalled func() (*data.ValidatorStatisticsResponse, error)
}

// GetValidatorStatistics -
func (v *ValidatorStatisticsProcessorStub) GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error) {
	return v.GetValidatorStatisticsCalled()
}

// GetAuctionList -
func (v *ValidatorStatisticsProcessorStub) GetAuctionList() (*data.AuctionListResponse, error) {
	return nil, nil
}
