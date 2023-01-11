package groups_test

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

var emptyGinHandler = func(_ *gin.Context) {}

func init() {
	gin.SetMode(gin.TestMode)
}

func startProxyServer(group data.GroupHandler, path string) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	routes := ws.Group(path)
	group.RegisterRoutes(routes, data.ApiRoutesConfig{}, emptyGinHandler, emptyGinHandler, emptyGinHandler)
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		fmt.Println(err.Error())
	}
}
