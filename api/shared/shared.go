package shared

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// RespondWith will respond with the generic API response
func RespondWith(c *gin.Context, status int, dataField interface{}, error string, code data.ReturnCode) {
	c.JSON(
		status,
		data.GenericAPIResponse{
			Data:  dataField,
			Error: error,
			Code:  code,
		},
	)
}

// RespondWithInvalidAppContext will be called when the application's context is invalid
func RespondWithInvalidAppContext(c *gin.Context) {
	RespondWith(c, http.StatusInternalServerError, nil, errors.ErrInvalidAppContext.Error(), data.ReturnCodeInternalError)
}
