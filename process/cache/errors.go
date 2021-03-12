package cache

import "errors"

// ErrNilHeartbeatsInCache signals that the heartbeats response stored in cache is nil
var ErrNilHeartbeatsInCache = errors.New("nil heartbeat response in cache")

// ErrNilHeartbeatsToStoreInCache signals that the provided heartbeats response is nil
var ErrNilHeartbeatsToStoreInCache = errors.New("nil heartbeat response to store in cache")

// ErrNilValidatorStatsInCache signals that the heartbeats response stored in cache is nil
var ErrNilValidatorStatsInCache = errors.New("nil validator statistics response in cache")

// ErrNilValidatorStatsToStoreInCache signals that the provided validator statistics is nil
var ErrNilValidatorStatsToStoreInCache = errors.New("nil validator statistics to store in cache")

// ErrNilGenericApiResponseInCache signals that the generic api response stored in cache is nil
var ErrNilGenericApiResponseInCache = errors.New("nil generic api response in cache")

// ErrNilGenericApiResponseToStoreInCache signals that the provided generic api response is nil
var ErrNilGenericApiResponseToStoreInCache = errors.New("nil generic api response to store in cache")
