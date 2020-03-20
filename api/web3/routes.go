package web3

import (
	"encoding/json"
	"github.com/ElrondNetwork/elrond-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.POST("", GetData)
}

func GetData(c *gin.Context) {
	receivedBody, err := getDataFromBody(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	epf, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	response, err := epf.PrepareDataForRequest(receivedBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jsonrpc": response.JsonRpc,
		"id":      response.Id,
		"result":  response.Result,
	})
}

func getDataFromBody(request *http.Request) (data.RequestBodyWeb3, error) {
	buf := make([]byte, 1024)
	num, err := request.Body.Read(buf)
	if err != nil && err != io.EOF {
		return data.RequestBodyWeb3{}, err
	}

	reqBody := append([]byte(nil), buf[:num]...)

	var receivedBody data.RequestBodyWeb3

	err = json.Unmarshal(reqBody, &receivedBody)
	if err != nil {
		return data.RequestBodyWeb3{}, err
	}

	return receivedBody, nil
}
