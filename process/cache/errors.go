package cache

import "errors"

// ErrNilHeartbeatsInCache signals that the heartbeats response stored in cache is nil
var ErrNilHeartbeatsInCache = errors.New("nil heartbeat response in cache")

// ErrNilHeartbeatsToStoreInCache signals that the provided heartbeats response is nil
var ErrNilHeartbeatsToStoreInCache = errors.New("nil heartbeat response to store in cache")
