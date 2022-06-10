package middleware

import (
	"time"

	"github.com/ElrondNetwork/elrond-go/api/shared"
)

// RateLimiterHandler defines the actions that an implementation of rate limiter handler should do
type RateLimiterHandler interface {
	shared.MiddlewareProcessor
	ResetMap(version string)
}

// StatusMetricsExtractor defines what a status metrics extractor should do
type StatusMetricsExtractor interface {
	AddRequestData(path string, withError bool, duration time.Duration)
	IsInterfaceNil() bool
}
