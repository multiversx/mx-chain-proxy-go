package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterHandler defines the actions that an implementation of rate limiter handler should do
type RateLimiterHandler interface {
	MiddlewareProcessor
	ResetMap(version string)
}

// StatusMetricsExtractor defines what a status metrics extractor should do
type StatusMetricsExtractor interface {
	AddRequestData(path string, withError bool, duration time.Duration)
	IsInterfaceNil() bool
}

// MiddlewareProcessor defines a processor used internally by the web server when processing requests
type MiddlewareProcessor interface {
	MiddlewareHandlerFunc() gin.HandlerFunc
	IsInterfaceNil() bool
}
