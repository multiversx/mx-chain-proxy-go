package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type ExternalStorageConnectorStub struct {
	GetTransactionsByAddressCalled       func(address string) ([]data.DatabaseTransaction, error)
	GetAtlasBlockByShardIDAndNonceCalled func(shardID uint32, nonce uint64) (data.ApiBlock, error)
}

// GetTransactionsByAddress -
func (e *ExternalStorageConnectorStub) GetTransactionsByAddress(address string) ([]data.DatabaseTransaction, error) {
	if e.GetTransactionsByAddressCalled != nil {
		return e.GetTransactionsByAddressCalled(address)
	}

	return []data.DatabaseTransaction{{Fee: "0"}}, nil
}

// GetAtlasBlockByShardIDAndNonce -
func (e *ExternalStorageConnectorStub) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error) {
	if e.GetAtlasBlockByShardIDAndNonceCalled != nil {
		return e.GetAtlasBlockByShardIDAndNonceCalled(shardID, nonce)
	}

	return data.ApiBlock{Hash: "hash"}, nil
}

// IsInterfaceNil -
func (e *ExternalStorageConnectorStub) IsInterfaceNil() bool {
	return e == nil
}
