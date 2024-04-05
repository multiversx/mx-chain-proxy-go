package groups

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type networkGroup struct {
	facade NetworkFacadeHandler
	*baseGroup
}

// NewNetworkGroup returns a new instance of networkGroup
func NewNetworkGroup(facadeHandler data.FacadeHandler) (*networkGroup, error) {
	facade, ok := facadeHandler.(NetworkFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &networkGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/status/:shard", Handler: ng.getNetworkStatusData, Method: http.MethodGet},
		{Path: "/config", Handler: ng.getNetworkConfigData, Method: http.MethodGet},
		{Path: "/economics", Handler: ng.getEconomicsData, Method: http.MethodGet},
		{Path: "/esdts", Handler: ng.getEsdts, Method: http.MethodGet},
		{Path: "/esdt/fungible-tokens", Handler: ng.getEsdtHandlerFunc(data.FungibleTokens), Method: http.MethodGet},
		{Path: "/esdt/semi-fungible-tokens", Handler: ng.getEsdtHandlerFunc(data.SemiFungibleTokens), Method: http.MethodGet},
		{Path: "/esdt/non-fungible-tokens", Handler: ng.getEsdtHandlerFunc(data.NonFungibleTokens), Method: http.MethodGet},
		{Path: "/esdt/supply/:token", Handler: ng.getESDTSupply, Method: http.MethodGet},
		{Path: "/enable-epochs", Handler: ng.getEnableEpochs, Method: http.MethodGet},
		{Path: "/direct-staked-info", Handler: ng.getDirectStakedInfo, Method: http.MethodGet},
		{Path: "/delegated-info", Handler: ng.getDelegatedInfo, Method: http.MethodGet},
		{Path: "/ratings", Handler: ng.getRatingsConfig, Method: http.MethodGet},
		{Path: "/genesis-nodes", Handler: ng.getGenesisNodes, Method: http.MethodGet},
		{Path: "/gas-configs", Handler: ng.getGasConfigs, Method: http.MethodGet},
		{Path: "/trie-statistics/:shard", Handler: ng.getTrieStatistics, Method: http.MethodGet},
		{Path: "/epoch-start/:shard/by-epoch/:epoch", Handler: ng.getEpochStartData, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getNetworkStatusData will expose the node network metrics for the given shard
func (group *networkGroup) getNetworkStatusData(c *gin.Context) {
	shardIDUint, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrInvalidShardIDParam.Error(), data.ReturnCodeRequestError)
		return
	}

	networkStatusResults, err := group.facade.GetNetworkStatusMetrics(shardIDUint)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkStatusResults)
}

// getNetworkConfigData will expose the node network metrics for the given shard
func (group *networkGroup) getNetworkConfigData(c *gin.Context) {
	networkConfigResults, err := group.facade.GetNetworkConfigMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkConfigResults)
}

// getEconomicsData will expose the economics data metrics from an observer (if any available) in json format
func (group *networkGroup) getEconomicsData(c *gin.Context) {
	economicsData, err := group.facade.GetEconomicsDataMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, economicsData)
}

func (group *networkGroup) getEsdtHandlerFunc(tokenType string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokens, err := group.facade.GetAllIssuedESDTs(tokenType)
		if err != nil {
			shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
			return
		}

		c.JSON(http.StatusOK, tokens)
	}
}

// getDirectStakedInfo will expose the direct staked values from a metachain observer in json format
func (group *networkGroup) getDirectStakedInfo(c *gin.Context) {
	directStakedInfo, err := group.facade.GetDirectStakedInfo()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, directStakedInfo)
}

// getDelegatedInfo will expose the delegated info values from a metachain observer in json format
func (group *networkGroup) getDelegatedInfo(c *gin.Context) {
	delegatedInfo, err := group.facade.GetDelegatedInfo()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, delegatedInfo)
}

// getEsdts will expose all the issued ESDTs
func (group *networkGroup) getEsdts(c *gin.Context) {
	allIssuedESDTs, err := group.facade.GetAllIssuedESDTs("")
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, allIssuedESDTs)
}

func (group *networkGroup) getEnableEpochs(c *gin.Context) {
	enableEpochsMetrics, err := group.facade.GetEnableEpochsMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, enableEpochsMetrics)
}

func (group *networkGroup) getESDTSupply(c *gin.Context) {
	tokenIdentifier := c.Param("token")
	if tokenIdentifier == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrGetESDTTokenData.Error(), errors.ErrEmptyTokenIdentifier.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	esdtSupply, err := group.facade.GetESDTSupply(tokenIdentifier)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, esdtSupply)
}

// getRatingsConfig will expose the ratings configuration
func (group *networkGroup) getRatingsConfig(c *gin.Context) {
	networkConfigResults, err := group.facade.GetRatingsConfig()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkConfigResults)
}

// getGenesisNodes will expose genesis nodes public keys
func (group *networkGroup) getGenesisNodes(c *gin.Context) {
	genesisNodes, err := group.facade.GetGenesisNodesPubKeys()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, genesisNodes)
}

// getGasConfigs will expose gas configs
func (group *networkGroup) getGasConfigs(c *gin.Context) {
	gasConfigs, err := group.facade.GetGasConfigs()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, gasConfigs)
}

// getTrieStatistics will expose trie statistics
func (group *networkGroup) getTrieStatistics(c *gin.Context) {
	shardID, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrInvalidShardIDParam.Error(), data.ReturnCodeRequestError)
		return
	}

	trieStatistics, err := group.facade.GetTriesStatistics(shardID)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, trieStatistics)
}

// getEpochStartData will expose epoch-start data for a given shard and epoch
func (group *networkGroup) getEpochStartData(c *gin.Context) {
	epoch, err := shared.FetchEpochFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, fmt.Sprintf("error while parsing the epoch: %s", err.Error()))
		return
	}

	shardID, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, fmt.Sprintf("error while parsing the shard ID: %s", err.Error()))
		return
	}

	epochStartData, err := group.facade.GetEpochStartData(epoch, shardID)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, epochStartData)
}
