package mock

// TimedCacheStub -
type TimedCacheStub struct {
}

// Put -
func (stub *TimedCacheStub) Put(_ []byte, _ interface{}, _ int) bool {
	return false
}

// Get -
func (stub *TimedCacheStub) Get(_ []byte) (value interface{}, ok bool) {
	return nil, false
}

// Close -
func (stub *TimedCacheStub) Close() error {
	return nil
}

// IsInterfaceNil -
func (stub *TimedCacheStub) IsInterfaceNil() bool {
	return stub == nil
}
