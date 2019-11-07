package mock

type AddressContainerMock struct {
	BytesField []byte
}

func (adr *AddressContainerMock) Bytes() []byte {
	return adr.BytesField
}

func (adr *AddressContainerMock) IsInterfaceNil() bool {
	return adr == nil
}
