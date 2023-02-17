package mock

import crypto "github.com/multiversx/mx-chain-crypto-go"

type PrivateKeysLoaderStub struct {
	PrivateKeysByShardCalled func() (map[uint32][]crypto.PrivateKey, error)
}

func (pkls *PrivateKeysLoaderStub) PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error) {
	return pkls.PrivateKeysByShardCalled()
}
