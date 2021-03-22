package middleware

import "github.com/ElrondNetwork/elrond-go/api"

// RateLimiterHandler defines the actions that an implementation of rate limiter handler should do
type RateLimiterHandler interface {
	api.MiddlewareProcessor
	ResetMap(version string)
}
