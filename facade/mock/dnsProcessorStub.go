package mock

type DnsProcessorStub struct {
	GetDnsAddressesCalled          func() ([]string, error)
	GetDnsAddressForUsernameCalled func(username string) (string, error)
}

// GetDnsAddresses -
func (d *DnsProcessorStub) GetDnsAddresses() ([]string, error) {
	if d.GetDnsAddressesCalled != nil {
		return d.GetDnsAddressesCalled()
	}

	return nil, nil
}

// GetDnsAddressForUsername -
func (d *DnsProcessorStub) GetDnsAddressForUsername(username string) (string, error) {
	if d.GetDnsAddressForUsernameCalled != nil {
		return d.GetDnsAddressForUsernameCalled(username)
	}

	return "", nil
}
