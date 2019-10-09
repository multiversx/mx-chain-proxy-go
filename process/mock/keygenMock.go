package mock

import "github.com/ElrondNetwork/elrond-go/crypto"

type KeygenStub struct {
	GeneratePairCalled            func() (crypto.PrivateKey, crypto.PublicKey)
	PrivateKeyFromByteArrayCalled func(b []byte) (crypto.PrivateKey, error)
	PublicKeyFromByteArrayCalled  func(b []byte) (crypto.PublicKey, error)
	SuiteCalled                   func() crypto.Suite
}

func (kgs *KeygenStub) GeneratePair() (crypto.PrivateKey, crypto.PublicKey) {
	return kgs.GeneratePairCalled()
}

func (kgs *KeygenStub) PrivateKeyFromByteArray(b []byte) (crypto.PrivateKey, error) {
	return kgs.PrivateKeyFromByteArrayCalled(b)
}

func (kgs *KeygenStub) PublicKeyFromByteArray(b []byte) (crypto.PublicKey, error) {
	return kgs.PublicKeyFromByteArrayCalled(b)
}

func (kgs *KeygenStub) Suite() crypto.Suite {
	return kgs.SuiteCalled()
}

func (kgs *KeygenStub) IsInterfaceNil() bool {
	return kgs == nil
}
