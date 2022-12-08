package middleware

import "errors"

// ErrNilLimitsMapForEndpoints signals that a nil limits map has been provided
var ErrNilLimitsMapForEndpoints = errors.New("nil limits map")

// ErrNilStatusMetricsExtractor signals that a nil status metrics extractor has been provided
var ErrNilStatusMetricsExtractor = errors.New("nil status metrics extractor")
