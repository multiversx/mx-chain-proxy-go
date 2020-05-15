package database

import (
	"errors"

	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

var errDatabaseConnectionIsDisabled = errors.New("database connection is disabled")

type disabledElasticSearchConnector struct{}

func NewDisabledElasticSearchConnector() *disabledElasticSearchConnector {
	return new(disabledElasticSearchConnector)
}

// GetTransactionsByAddress will return error because database connection is disabled
func (desc *disabledElasticSearchConnector) GetTransactionsByAddress(_ string) ([]indexer.Transaction, error) {
	return nil, errDatabaseConnectionIsDisabled
}

// GetBlockByShardIDAndNonce will return error because database connection is disabled
func (desc *disabledElasticSearchConnector) GetBlockByShardIDAndNonce(_ uint32, _ uint64) (data.ApiBlock, error) {
	return data.ApiBlock{}, errDatabaseConnectionIsDisabled
}

// IsInterfaceNil -
func (desc *disabledElasticSearchConnector) IsInterfaceNil() bool {
	return desc == nil
}
