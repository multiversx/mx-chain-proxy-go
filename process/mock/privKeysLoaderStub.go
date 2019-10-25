package mock

import "github.com/ElrondNetwork/elrond-go/crypto"

type PrivateKeysLoaderStub struct {
	MapOfPrivateKeysByShardCalled func() (map[uint32][]crypto.PrivateKey, error)
}

func (pkls *PrivateKeysLoaderStub) MapOfPrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error) {
	return pkls.MapOfPrivateKeysByShardCalled()
}
