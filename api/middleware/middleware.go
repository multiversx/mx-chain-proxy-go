package middleware

import "github.com/gin-gonic/gin"

// ElrondProxyHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type ElrondProxyHandler interface {
}

// WithElrondProxyFacade middleware will set up an ElrondFacade object in the gin context
func WithElrondProxyFacade(elrondProxyFacade ElrondProxyHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("elrondProxyFacade", elrondProxyFacade)
		c.Next()
	}
}
