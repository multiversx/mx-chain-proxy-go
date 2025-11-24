package cache

import (
	"context"
	"sync"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("cache")

const minDuration = time.Second

type entry struct {
	timestamp time.Time
	value     interface{}
}
type timeCacher struct {
	*sync.RWMutex
	data       map[string]*entry
	duration   time.Duration
	cancelFunc func()
}

// NewTimeCacher creates a new timeCacher
func NewTimeCacher(cacheExpiry time.Duration) (*timeCacher, error) {
	if cacheExpiry < minDuration {
		return nil, errInvalidCacheExpiry
	}

	tc := &timeCacher{
		RWMutex:  &sync.RWMutex{},
		data:     make(map[string]*entry),
		duration: cacheExpiry,
	}

	var ctx context.Context
	ctx, tc.cancelFunc = context.WithCancel(context.Background())
	go tc.startSweeping(ctx)

	return tc, nil
}

// startSweeping handles sweeping the time cache
func (tc *timeCacher) startSweeping(ctx context.Context) {
	timer := time.NewTimer(tc.duration)
	defer timer.Stop()

	for {
		timer.Reset(tc.duration)

		select {
		case <-timer.C:
			tc.sweep()
		case <-ctx.Done():
			log.Info("closing mapTimeCacher's sweep go routine...")
			return
		}
	}
}

// Put will add the key, value and provided duration, overriding values if the data already existed
// It also operates on the locker so the call is concurrent safe
func (tc *timeCacher) Put(key []byte, value interface{}) error {
	if len(key) == 0 {
		return errEmptyKey
	}

	tc.Lock()
	tc.data[string(key)] = &entry{
		timestamp: time.Now(),
		value:     value,
	}
	tc.Unlock()

	return nil
}

// Get returns a key's value from the cache
func (tc *timeCacher) Get(key []byte) (interface{}, bool) {
	tc.RLock()
	defer tc.RUnlock()

	v, ok := tc.data[string(key)]
	if !ok {
		return nil, ok
	}

	return v.value, ok
}

// sweep iterates over all contained elements checking if the element is still valid to be kept
// It also operates on the locker so the call is concurrent safe
func (tc *timeCacher) sweep() {
	tc.Lock()
	defer tc.Unlock()

	for key, element := range tc.data {
		isOldElement := time.Since(element.timestamp) > tc.duration
		if isOldElement {
			delete(tc.data, key)
		}
	}
}

// Close will close the internal sweep go routine
func (tc *timeCacher) Close() error {
	if tc.cancelFunc != nil {
		tc.cancelFunc()
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (tc *timeCacher) IsInterfaceNil() bool {
	return tc == nil
}
