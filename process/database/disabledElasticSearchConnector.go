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

// GetTransactionsByAddress will return error because database connection is offline
func (desc *disabledElasticSearchConnector) GetTransactionsByAddress(_ string) ([]data.DatabaseTransaction, error) {
	return nil, errDatabaseConnectionIsDisabled
}

// GetAtlasBlockByShardIDAndNonce will return error because database connection is offline
func (desc *disabledElasticSearchConnector) GetAtlasBlockByShardIDAndNonce(_ uint32, _ uint64) (data.AtlasBlock, error) {
	return data.AtlasBlock{}, errDatabaseConnectionIsDisabled
}

// IsInterfaceNil -
func (desc *disabledElasticSearchConnector) IsInterfaceNil() bool {
	return desc == nil
}
