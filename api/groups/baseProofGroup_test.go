package groups_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Response -
type Response struct {
	Data  map[string]interface{} `json:"data"`
	Error string                 `json:"error"`
	Code  string                 `json:"code"`
}

func TestNewProofGroup_WrongFacadeShouldErr(t *testing.T) {
	t.Parallel()

	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewProofGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetProof_FailWhenFacadeGetProofFails(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	returnedError := "getProof error"
	facade := &mock.FacadeStub{
		GetProofCalled: func(rh string, addr string) (*data.GenericAPIResponse, error) {
			assert.Equal(t, rootHash, rh)
			assert.Equal(t, address, addr)
			return nil, fmt.Errorf(returnedError)
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	req, err := http.NewRequest("GET", "/proof/root-hash/"+rootHash+"/address/"+address, nil)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := Response{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, returnedError, response.Error)
	assert.Empty(t, response.Data)
}

func TestGetProof(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	proof := []string{"valid", "proof"}

	facade := &mock.FacadeStub{
		GetProofCalled: func(rh string, addr string) (*data.GenericAPIResponse, error) {
			assert.Equal(t, rootHash, rh)
			assert.Equal(t, address, addr)
			return &data.GenericAPIResponse{Data: proof}, nil
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	req, err := http.NewRequest("GET", "/proof/root-hash/"+rootHash+"/address/"+address, nil)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)

	proofs, ok := response.Data.([]interface{})
	assert.True(t, ok)

	proof1 := proofs[0].(string)
	proof2 := proofs[1].(string)

	assert.Equal(t, "valid", proof1)
	assert.Equal(t, "proof", proof2)
}

func TestVerifyProof_FailWhenFacadeVerifyProofFails(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	proof := "proof"
	returnedError := "getProof error"
	facade := &mock.FacadeStub{
		VerifyProofCalled: func(rh string, addr string, p []string) (*data.GenericAPIResponse, error) {
			assert.Equal(t, rootHash, rh)
			assert.Equal(t, address, addr)
			assert.Equal(t, []string{proof}, p)
			return nil, fmt.Errorf(returnedError)
		},
	}
	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	varifyProofParams := data.VerifyProofRequest{
		RootHash: rootHash,
		Address:  address,
		Proof:    []string{proof},
	}
	verifyProofBytes, _ := json.Marshal(varifyProofParams)
	req, err := http.NewRequest("POST", "/proof/verify", bytes.NewBuffer(verifyProofBytes))
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, returnedError, response.Error)
	assert.Empty(t, response.Data)
}

func TestVerifyProof(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	proof := "proof"

	facade := &mock.FacadeStub{
		VerifyProofCalled: func(rh string, addr string, p []string) (*data.GenericAPIResponse, error) {
			assert.Equal(t, rootHash, rh)
			assert.Equal(t, address, addr)
			assert.Equal(t, []string{proof}, p)
			return &data.GenericAPIResponse{Data: true}, nil
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	varifyProofParams := data.VerifyProofRequest{
		RootHash: rootHash,
		Address:  address,
		Proof:    []string{proof},
	}
	verifyProofBytes, _ := json.Marshal(varifyProofParams)
	req, err := http.NewRequest("POST", "/proof/verify", bytes.NewBuffer(verifyProofBytes))
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)

	isValid, ok := response.Data.(bool)
	assert.True(t, ok)
	assert.True(t, isValid)
}

func TestGetProofDataTrie_FailWhenFacadeGetProofFails(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	key := "key"
	returnedError := "getProofDataTrie error"
	facade := &mock.FacadeStub{
		GetProofDataTrieCalled: func(rh string, addr string, k string) (*data.GenericAPIResponse, error) {
			assert.Equal(t, rootHash, rh)
			assert.Equal(t, address, addr)
			assert.Equal(t, key, k)
			return nil, fmt.Errorf(returnedError)
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	endpoint := fmt.Sprintf("/proof/root-hash/%s/address/%s/key/%s", rootHash, address, key)
	req, err := http.NewRequest("GET", endpoint, nil)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, returnedError, response.Error)
	assert.Empty(t, response.Data)
}

func TestGetProofDataTrie(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	key := "key"
	proof := []string{"valid", "proof"}

	facade := &mock.FacadeStub{
		GetProofDataTrieCalled: func(rh string, addr string, k string) (*data.GenericAPIResponse, error) {
			assert.Equal(t, rootHash, rh)
			assert.Equal(t, address, addr)
			assert.Equal(t, key, k)
			return &data.GenericAPIResponse{Data: proof}, nil
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	endpoint := fmt.Sprintf("/proof/root-hash/%s/address/%s/key/%s", rootHash, address, key)
	req, err := http.NewRequest("GET", endpoint, nil)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)

	proofs, ok := response.Data.([]interface{})
	assert.True(t, ok)

	proof1 := proofs[0].(string)
	proof2 := proofs[1].(string)

	assert.Equal(t, "valid", proof1)
	assert.Equal(t, "proof", proof2)
}
