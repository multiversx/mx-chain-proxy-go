package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

type ExternalStorageConnectorStub struct {
	GetAtlasBlockByShardIDAndNonceCalled func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
}

// GetAtlasBlockByShardIDAndNonce -
func (e *ExternalStorageConnectorStub) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	if e.GetAtlasBlockByShardIDAndNonceCalled != nil {
		return e.GetAtlasBlockByShardIDAndNonceCalled(shardID, nonce)
	}

	return data.AtlasBlock{Hash: "hash"}, nil
}

// IsInterfaceNil -
func (e *ExternalStorageConnectorStub) IsInterfaceNil() bool {
	return e == nil
}
