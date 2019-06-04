package mock

import "github.com/ElrondNetwork/elrond-go-sandbox/data/state"

type AddressConverterStub struct {
	CreateAddressFromPublicKeyBytesCalled func(pubKey []byte) (state.AddressContainer, error)
	ConvertToHexCalled                    func(addressContainer state.AddressContainer) (string, error)
	CreateAddressFromHexCalled            func(hexAddress string) (state.AddressContainer, error)
	PrepareAddressBytesCalled             func(addressBytes []byte) ([]byte, error)
}

func (acs *AddressConverterStub) CreateAddressFromPublicKeyBytes(pubKey []byte) (state.AddressContainer, error) {
	if acs.CreateAddressFromPublicKeyBytesCalled != nil {
		return acs.CreateAddressFromPublicKeyBytesCalled(pubKey)
	}

	return nil, errNotImplemented
}

func (acs *AddressConverterStub) ConvertToHex(addressContainer state.AddressContainer) (string, error) {
	if acs.ConvertToHexCalled != nil {
		return acs.ConvertToHexCalled(addressContainer)
	}

	return "", errNotImplemented
}

func (acs *AddressConverterStub) CreateAddressFromHex(hexAddress string) (state.AddressContainer, error) {
	if acs.CreateAddressFromHexCalled != nil {
		return acs.CreateAddressFromHexCalled(hexAddress)
	}

	return nil, errNotImplemented
}

func (acs *AddressConverterStub) PrepareAddressBytes(addressBytes []byte) ([]byte, error) {
	if acs.PrepareAddressBytesCalled != nil {
		return acs.PrepareAddressBytesCalled(addressBytes)
	}

	return nil, errNotImplemented
}
