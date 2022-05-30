package groups

import (
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/gin-gonic/gin"
)

func parseBlockQueryOptions(c *gin.Context) (common.BlockQueryOptions, error) {
	withTxs, err := parseBoolUrlParam(c, "withTxs")
	if err != nil {
		return common.BlockQueryOptions{}, err
	}

	withLogs, err := parseBoolUrlParam(c, "withLogs")
	if err != nil {
		return common.BlockQueryOptions{}, err
	}

	options := common.BlockQueryOptions{WithTransactions: withTxs, WithLogs: withLogs}
	return options, nil
}

func parseHyperblockQueryOptions(c *gin.Context) (common.HyperblockQueryOptions, error) {
	withLogs, err := parseBoolUrlParam(c, "withLogs")
	if err != nil {
		return common.HyperblockQueryOptions{}, err
	}

	options := common.HyperblockQueryOptions{WithLogs: withLogs}
	return options, nil
}

func parseBoolUrlParam(c *gin.Context, name string) (bool, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return false, nil
	}

	return strconv.ParseBool(param)
}
