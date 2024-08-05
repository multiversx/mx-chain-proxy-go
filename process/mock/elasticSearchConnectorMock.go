package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

type ElasticSearchConnectorMock struct {
}

// GetAtlasBlockByShardIDAndNonce -
func (escm *ElasticSearchConnectorMock) GetAtlasBlockByShardIDAndNonce(_ uint32, _ uint64) (data.AtlasBlock, error) {
	return data.AtlasBlock{}, nil
}

// IsInterfaceNil -
func (escm *ElasticSearchConnectorMock) IsInterfaceNil() bool {
	return escm == nil
}
