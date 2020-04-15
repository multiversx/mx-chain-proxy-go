package mock

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go/data/state"
)

type PubKeyConverterMock struct {
	CreateAddressFromBytesCalled func(pkBytes []byte) (state.AddressContainer, error)
}

func (p *PubKeyConverterMock) Len() int {
	return 32
}

func (p *PubKeyConverterMock) Decode(humanReadable string) ([]byte, error) {
	return hex.DecodeString(humanReadable)
}

func (p *PubKeyConverterMock) Encode(pkBytes []byte) string {
	return hex.EncodeToString(pkBytes)
}

func (p *PubKeyConverterMock) CreateAddressFromString(humanReadable string) (state.AddressContainer, error) {
	return nil, nil
}

func (p *PubKeyConverterMock) CreateAddressFromBytes(pkBytes []byte) (state.AddressContainer, error) {
	if p.CreateAddressFromBytesCalled != nil {
		return p.CreateAddressFromBytesCalled(pkBytes)
	}

	return nil, nil
}

func (p *PubKeyConverterMock) IsInterfaceNil() bool {
	return p == nil
}
