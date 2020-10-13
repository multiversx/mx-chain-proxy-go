package shared

import "github.com/gin-gonic/gin"

type EndpointHandlerData struct {
	Handler gin.HandlerFunc
	Method  string
}
