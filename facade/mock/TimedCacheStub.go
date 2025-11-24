package mock

// TimedCacheMock -
type TimedCacheMock struct {
	cache map[string]interface{}
}

// NewTimedCacheMock -
func NewTimedCacheMock() *TimedCacheMock {
	return &TimedCacheMock{cache: make(map[string]interface{})}
}

// Put -
func (mock *TimedCacheMock) Put(key []byte, value interface{}) error {
	mock.cache[string(key)] = value
	return nil
}

// Get -
func (mock *TimedCacheMock) Get(key []byte) (value interface{}, ok bool) {
	val, found := mock.cache[string(key)]
	return val, found
}
