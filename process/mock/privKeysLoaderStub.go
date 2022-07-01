package mock

import "github.com/ElrondNetwork/elrond-go-crypto"

type PrivateKeysLoaderStub struct {
	PrivateKeysByShardCalled func() (map[uint32][]crypto.PrivateKey, error)
}

func (pkls *PrivateKeysLoaderStub) PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error) {
	return pkls.PrivateKeysByShardCalled()
}
