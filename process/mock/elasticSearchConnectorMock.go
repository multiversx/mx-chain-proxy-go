package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type ElasticSearchConnectorMock struct {
}

// GetTransactionsByAddress -
func (escm *ElasticSearchConnectorMock) GetTransactionsByAddress(_ string) ([]data.DatabaseTransaction, error) {
	return nil, nil
}

// GetAtlasBlockByShardIDAndNonce -
func (escm *ElasticSearchConnectorMock) GetAtlasBlockByShardIDAndNonce(_ uint32, _ uint64) (data.AtlasBlock, error) {
	return data.AtlasBlock{}, nil
}

// IsInterfaceNil -
func (escm *ElasticSearchConnectorMock) IsInterfaceNil() bool {
	return escm == nil
}
