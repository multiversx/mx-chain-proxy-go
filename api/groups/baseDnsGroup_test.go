package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dnsPath = "/dns"

type dnsAddressesResponseData struct {
	Addresses []string `json:"addresses"`
}

type dnsAddressesResponse struct {
	Data  dnsAddressesResponseData `json:"data"`
	Error string                   `json:"error"`
	Code  string                   `json:"code"`
}

type dnsAddressForUsernameResponseData struct {
	Address string `json:"address"`
}

type dnsAddressesForUsernameResponse struct {
	Data  dnsAddressForUsernameResponseData `json:"data"`
	Error string                            `json:"error"`
	Code  string                            `json:"code"`
}

func TestNewDnsGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewDnsGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestDnsGroup_getAllDnsAddresses_FacadeErrShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")

	facade := &mock.Facade{
		GetDnsAddressesCalled: func() ([]string, error) {
			return nil, localErr
		},
	}

	dnsGroup, err := groups.NewDnsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(dnsGroup, dnsPath)

	req, _ := http.NewRequest("GET", "/dns/all", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var result dnsAddressesResponse
	loadResponse(resp.Body, &result)

	assert.Contains(t, result.Error, localErr.Error())
}

func TestDnsGroup_getAllDnsAddresses_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedAddresses := []string{"address0", "address1"}
	facade := &mock.Facade{
		GetDnsAddressesCalled: func() ([]string, error) {
			return expectedAddresses, nil
		},
	}

	dnsGroup, err := groups.NewDnsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(dnsGroup, dnsPath)

	req, _ := http.NewRequest("GET", "/dns/all", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result dnsAddressesResponse
	loadResponse(resp.Body, &result)
	assert.Empty(t, result.Error)
	assert.Equal(t, expectedAddresses, result.Data.Addresses)
}

func TestDnsGroup_getDnsAddressForUsername_FacadeErrShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")

	facade := &mock.Facade{
		GetDnsAddressForUsernameCalled: func(_ string) (string, error) {
			return "", localErr
		},
	}

	dnsGroup, err := groups.NewDnsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(dnsGroup, dnsPath)

	req, _ := http.NewRequest("GET", "/dns/username/test", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var result dnsAddressesForUsernameResponse
	loadResponse(resp.Body, &result)

	assert.Contains(t, result.Error, localErr.Error())
}

func TestDnsGroup_getDnsAddressForUsername_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedAddress := "address37"
	facade := &mock.Facade{
		GetDnsAddressForUsernameCalled: func(_ string) (string, error) {
			return expectedAddress, nil
		},
	}

	dnsGroup, err := groups.NewDnsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(dnsGroup, dnsPath)

	req, _ := http.NewRequest("GET", "/dns/username/test", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result dnsAddressesForUsernameResponse
	loadResponse(resp.Body, &result)
	assert.Empty(t, result.Error)
	assert.Equal(t, expectedAddress, result.Data.Address)
}
