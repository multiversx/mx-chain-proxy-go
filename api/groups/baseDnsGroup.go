package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type dnsGroup struct {
	facade DnsFacadeHandler
	*baseGroup
}

// NewDnsGroup returns a new instance of dnsGroup
func NewDnsGroup(facadeHandler data.FacadeHandler) (*dnsGroup, error) {
	facade, ok := facadeHandler.(DnsFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &dnsGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/all":                {Handler: ng.getAllDnsAddresses, Method: http.MethodGet},
		"/username/:username": {Handler: ng.getDnsAddressForUsername, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func (group *dnsGroup) getAllDnsAddresses(c *gin.Context) {
	dnsAddresses, err := group.facade.GetDnsAddresses()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"addresses": dnsAddresses}, "", data.ReturnCodeSuccess)
}

// getHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func (group *dnsGroup) getDnsAddressForUsername(c *gin.Context) {
	username := c.Param("username")
	if len(username) == 0 {
		shared.RespondWithBadRequest(c, "empty username provided")
	}

	dnsAddress, err := group.facade.GetDnsAddressForUsername(username)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"address": dnsAddress}, "", data.ReturnCodeSuccess)
}
