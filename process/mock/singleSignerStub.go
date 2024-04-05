package mock

import crypto "github.com/multiversx/mx-chain-crypto-go"

type SignerStub struct {
	SignCalled   func(private crypto.PrivateKey, msg []byte) ([]byte, error)
	VerifyCalled func(public crypto.PublicKey, msg []byte, sig []byte) error
}

func (s *SignerStub) Sign(private crypto.PrivateKey, msg []byte) ([]byte, error) {
	return s.SignCalled(private, msg)
}

func (s *SignerStub) Verify(public crypto.PublicKey, msg []byte, sig []byte) error {
	return s.VerifyCalled(public, msg, sig)
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *SignerStub) IsInterfaceNil() bool {
	return s == nil
}
