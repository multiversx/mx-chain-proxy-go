package database

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

var errDatabaseConnectionIsDisabled = errors.New("database connection is disabled")

type disabledElasticSearchConnector struct{}

func NewDisabledElasticSearchConnector() *disabledElasticSearchConnector {
	return new(disabledElasticSearchConnector)
}

// GetTransactionsByAddress will return error because database connection is disabled
func (desc *disabledElasticSearchConnector) GetTransactionsByAddress(_ string) ([]data.DatabaseTransaction, error) {
	return nil, errDatabaseConnectionIsDisabled
}

// GetLatestBlockHeight will return error because database connection is disabled
func (desc *disabledElasticSearchConnector) GetLatestBlockHeight() (uint64, error) {
	return 0, errDatabaseConnectionIsDisabled
}

// GetBlockByNonce will return error because database connection is disabled
func (desc *disabledElasticSearchConnector) GetBlockByNonce(_ uint64) (data.ApiBlock, error) {
	return data.ApiBlock{}, errDatabaseConnectionIsDisabled
}

// IsInterfaceNil -
func (desc *disabledElasticSearchConnector) IsInterfaceNil() bool {
	return desc == nil
}
