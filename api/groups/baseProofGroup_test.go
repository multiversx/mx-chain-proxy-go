package groups_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"

	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
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
	facade := &mock.Facade{
		GetProofCalled: func(rh []byte, addr []byte) ([][]byte, error) {
			assert.Equal(t, []byte(rootHash), rh)
			assert.Equal(t, []byte(address), addr)
			return nil, fmt.Errorf(returnedError)
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	req, err := http.NewRequest("GET", "/proof/root-hash/"+rootHash+"/address/"+address, nil)

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
	proof := [][]byte{[]byte("valid"), []byte("proof")}

	facade := &mock.Facade{
		GetProofCalled: func(rh []byte, addr []byte) ([][]byte, error) {
			assert.Equal(t, []byte(rootHash), rh)
			assert.Equal(t, []byte(address), addr)
			return proof, nil
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	req, err := http.NewRequest("GET", "/proof/root-hash/"+rootHash+"/address/"+address, nil)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := Response{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)

	proofs, ok := response.Data["proof"].([]interface{})
	assert.True(t, ok)

	proof1 := proofs[0].(string)
	proof2 := proofs[1].(string)

	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("valid")), proof1)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("proof")), proof2)
}

func TestVerifyProof_FailWhenFacadeVerifyProofFails(t *testing.T) {
	t.Parallel()

	rootHash := "rootHash"
	address := "address"
	proof := "proof"
	returnedError := "getProof error"
	facade := &mock.Facade{
		VerifyProofCalled: func(rh []byte, addr []byte, p [][]byte) (bool, error) {
			assert.Equal(t, []byte(rootHash), rh)
			assert.Equal(t, []byte(address), addr)
			assert.Equal(t, [][]byte{[]byte(proof)}, p)
			return false, fmt.Errorf(returnedError)
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	varifyProofParams := data.VerifyProofRequest{
		RootHash: []byte(rootHash),
		Address:  []byte(address),
		Proof:    [][]byte{[]byte(proof)},
	}
	verifyProofBytes, _ := json.Marshal(varifyProofParams)
	req, err := http.NewRequest("POST", "/proof/verify", bytes.NewBuffer(verifyProofBytes))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := Response{}
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

	facade := &mock.Facade{
		VerifyProofCalled: func(rh []byte, addr []byte, p [][]byte) (bool, error) {
			assert.Equal(t, []byte(rootHash), rh)
			assert.Equal(t, []byte(address), addr)
			assert.Equal(t, [][]byte{[]byte(proof)}, p)
			return true, nil
		},
	}

	proofGroup, err := groups.NewProofGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(proofGroup, "/proof")

	varifyProofParams := data.VerifyProofRequest{
		RootHash: []byte(rootHash),
		Address:  []byte(address),
		Proof:    [][]byte{[]byte(proof)},
	}
	verifyProofBytes, _ := json.Marshal(varifyProofParams)
	req, err := http.NewRequest("POST", "/proof/verify", bytes.NewBuffer(verifyProofBytes))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := Response{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)

	isValid, ok := response.Data["ok"].(bool)
	assert.True(t, ok)
	assert.True(t, isValid)
}
