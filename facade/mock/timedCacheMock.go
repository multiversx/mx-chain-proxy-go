package mock

// TimedCacheMock -
type TimedCacheMock struct {
	Cache map[string]interface{}
}

// NewTimedCacheMock -
func NewTimedCacheMock() *TimedCacheMock {
	return &TimedCacheMock{Cache: make(map[string]interface{})}
}

// Put -
func (mock *TimedCacheMock) Put(key []byte, value interface{}, _ int) bool {
	mock.Cache[string(key)] = value
	return false
}

// Get -
func (mock *TimedCacheMock) Get(key []byte) (value interface{}, ok bool) {
	val, found := mock.Cache[string(key)]
	return val, found
}

// Close -
func (mock *TimedCacheMock) Close() error {
	return nil
}

// IsInterfaceNil -
func (mock *TimedCacheMock) IsInterfaceNil() bool {
	return mock == nil
}
