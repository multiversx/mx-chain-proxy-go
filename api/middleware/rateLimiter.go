package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/api/shared"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requestsMap    map[string]uint64
	mutRequestsMap sync.RWMutex
	limits         map[string]uint64
	mutLimits      sync.RWMutex
	countDuration  uint64
}

// NewRateLimiter returns a new instance of rateLimiter
func NewRateLimiter(limits map[string]uint64, countDuration uint64) (*rateLimiter, error) {
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

		limitForEndpoint, isEndpointLimited := rl.getEndpointLimit(endpoint)
		if !isEndpointLimited {
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("%s_%s", endpoint, clientIP)

		numRequests := rl.loadFromRequestsMap(key)
		rl.addInRequestsMap(key)
		if numRequests >= limitForEndpoint {
			printMessage := fmt.Sprintf("your IP exceeded the limit of %d requests in %ds for this endpoint", limitForEndpoint, rl.countDuration)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, shared.GenericAPIResponse{
				Data:  nil,
				Error: printMessage,
				Code:  shared.ReturnCodeRequestError,
			})
		}
	}
}

func (rl *rateLimiter) loadFromRequestsMap(key string) uint64 {
	rl.mutRequestsMap.RLock()
	entry, ok := rl.requestsMap[key]
	rl.mutRequestsMap.RUnlock()

	if !ok {
		return 0
	}

	return entry
}

func (rl *rateLimiter) addInRequestsMap(key string) {
	rl.mutRequestsMap.Lock()
	defer rl.mutRequestsMap.Unlock()

	_, ok := rl.requestsMap[key]
	if !ok {
		rl.requestsMap[key] = 1
		return
	}

	rl.requestsMap[key]++
}

func (rl *rateLimiter) getEndpointLimit(endpoint string) (uint64, bool) {
	limit, entryExists := rl.limits[endpoint]
	if !entryExists {
		return 0, false
	}

	return limit, true
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
