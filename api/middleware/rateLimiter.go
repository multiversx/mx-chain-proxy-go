package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// ReturnCodeRequestError defines a request which hasn't been executed successfully due to a bad request received
const ReturnCodeRequestError string = "bad_request"

type rateLimiter struct {
	requestsMap    map[string]uint64
	mutRequestsMap sync.RWMutex
	limits         map[string]uint64
	countDuration  time.Duration
}

// NewRateLimiter returns a new instance of rateLimiter
func NewRateLimiter(limits map[string]uint64, countDuration time.Duration) (*rateLimiter, error) {
	if limits == nil {
		return nil, ErrNilLimitsMapForEndpoints
	}
	return &rateLimiter{
		requestsMap:   make(map[string]uint64),
		limits:        limits,
		countDuration: countDuration,
	}, nil
}

// MiddlewareHandlerFunc returns the gin middleware for limiting the number of requests for a given endpoint
func (rl *rateLimiter) MiddlewareHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()

		limitForEndpoint, isEndpointLimited := rl.limits[endpoint]
		if !isEndpointLimited {
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("%s_%s", endpoint, clientIP)

		numRequests := rl.addInRequestsMap(key)
		if numRequests >= limitForEndpoint {
			printMessage := fmt.Sprintf("your IP exceeded the limit of %d requests in %v for this endpoint", limitForEndpoint, rl.countDuration)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, data.GenericAPIResponse{
				Data:  nil,
				Error: printMessage,
				Code:  data.ReturnCode(ReturnCodeRequestError),
			})
		}
	}
}

func (rl *rateLimiter) addInRequestsMap(key string) uint64 {
	rl.mutRequestsMap.Lock()
	defer rl.mutRequestsMap.Unlock()

	_, ok := rl.requestsMap[key]
	if !ok {
		rl.requestsMap[key] = 1
		return 1
	}

	rl.requestsMap[key]++

	return rl.requestsMap[key]
}

// ResetMap has to be called from outside at a given interval so the requests map will be cleaned and older restrictions
// would be erased
func (rl *rateLimiter) ResetMap(version string) {
	rl.mutRequestsMap.Lock()
	rl.requestsMap = make(map[string]uint64)
	rl.mutRequestsMap.Unlock()

	log.Info("rate limiter map has been reset", "version", version, "time", time.Now())
}

// IsInterfaceNil returns true if there is no value under the interface
func (rl *rateLimiter) IsInterfaceNil() bool {
	return rl == nil
}
