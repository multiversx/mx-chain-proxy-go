package groups_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

type dataResponse struct {
	Data string `json:"data"`
}

type simpleResponse struct {
	Data  dataResponse `json:"data"`
	Error string       `json:"error"`
}

type vmOutputResponse struct {
	Data      *vm.VMOutputApi `json:"data"`
	BlockInfo data.BlockInfo  `json:"blockInfo"`
}

type vmOutputGenericResponse struct {
	Data  vmOutputResponse `json:"data"`
	Error string           `json:"error"`
}

const vmValuesPath = "/vm-values"
const DummyScAddress = "erd1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsmsgldz"

func TestNewVmValuesGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewVmValuesGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetHex_ShouldWork(t *testing.T) {
	t.Parallel()

	valueBuff, _ := hex.DecodeString("DEADBEEF")

	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			return &vm.VMOutputApi{
				ReturnData: [][]byte{valueBuff},
			}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	response := simpleResponse{}
	statusCode := doPost(t, facade, "/vm-values/hex", request, &response)

	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "", response.Error)
	require.Equal(t, hex.EncodeToString(valueBuff), response.Data.Data)
}

func TestGetString_ShouldWork(t *testing.T) {
	t.Parallel()

	valueBuff := "DEADBEEF"

	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			return &vm.VMOutputApi{
				ReturnData: [][]byte{[]byte(valueBuff)},
			}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	response := simpleResponse{}
	statusCode := doPost(t, facade, "/vm-values/string", request, &response)

	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "", response.Error)
	require.Equal(t, valueBuff, response.Data.Data)
}

func TestGetInt_ShouldWork(t *testing.T) {
	t.Parallel()

	value := "1234567"

	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			returnData := big.NewInt(0)
			returnData.SetString(value, 10)
			return &vm.VMOutputApi{
				ReturnData: [][]byte{returnData.Bytes()},
			}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	response := simpleResponse{}
	statusCode := doPost(t, facade, "/vm-values/int", request, &response)

	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "", response.Error)
	require.Equal(t, value, response.Data.Data)
}

func TestQuery_ShouldWork(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {

			return &vm.VMOutputApi{
				ReturnData: [][]byte{big.NewInt(42).Bytes()},
			}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	response := vmOutputGenericResponse{}
	statusCode := doPost(t, facade, "/vm-values/query", request, &response)

	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "", response.Error)
	require.Equal(t, int64(42), big.NewInt(0).SetBytes(response.Data.Data.ReturnData[0]).Int64())
}

func TestQuery_ShouldWorkWithCoordinates(t *testing.T) {
	t.Parallel()

	providedNonce := uint64(123)
	providedBlockInfo := data.BlockInfo{
		Nonce:    providedNonce,
		Hash:     "block hash",
		RootHash: "block rootHash",
	}
	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			require.Equal(t, providedNonce, query.BlockNonce.Value)
			return &vm.VMOutputApi{
				ReturnData: [][]byte{big.NewInt(42).Bytes()},
			}, providedBlockInfo, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	response := vmOutputGenericResponse{}
	statusCode := doPost(t, facade, "/vm-values/query?blockNonce="+strconv.FormatUint(providedNonce, 10), request, &response)

	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "", response.Error)
	require.Equal(t, int64(42), big.NewInt(0).SetBytes(response.Data.Data.ReturnData[0]).Int64())
	require.Equal(t, providedBlockInfo, response.Data.BlockInfo)
}

func TestCreateSCQuery_ArgumentIsNotHexShouldErr(t *testing.T) {
	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{"bad arg"},
	}

	_, err := createSCQuery(&request)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "'bad arg' is not a valid hex string")
}

func TestAllRoutes_FacadeErrorsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("some random error")
	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			return nil, data.BlockInfo{}, errExpected
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	requireErrorOnAllRoutes(t, facade, request, errExpected)
}

func TestAllRoutes_WhenBadArgumentsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("not a valid hex string")
	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			return &vm.VMOutputApi{}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{"AA", "ZZ"},
	}

	requireErrorOnAllRoutes(t, facade, request, errExpected)
}

func TestAllRoutes_WhenNoVMReturnDataShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("no return data")
	facade := mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			return &vm.VMOutputApi{}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress: DummyScAddress,
		FuncName:  "function",
		Args:      []string{},
	}

	requireErrorOnGetSingleValueRoutes(t, &facade, request, errExpected)
}

func TestAllRoutes_WhenBadJsonShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			return &vm.VMOutputApi{}, data.BlockInfo{}, nil
		},
	}

	requireErrorOnGetSingleValueRoutes(t, &facade, []byte("dummy"), apiErrors.ErrInvalidJSONRequest)
}

func TestAllRoutes_WithSameScStateAndShouldBySyncedFilled(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		ExecuteSCQueryHandler: func(query *data.SCQuery) (vmOutput *vm.VMOutputApi, blockInfo data.BlockInfo, e error) {
			require.True(t, query.ShouldBeSynced)
			require.True(t, query.SameScState)
			return &vm.VMOutputApi{}, data.BlockInfo{}, nil
		},
	}

	request := groups.VMValueRequest{
		ScAddress:      DummyScAddress,
		FuncName:       "function",
		Args:           []string{},
		SameScState:    true,
		ShouldBeSynced: true,
	}

	response := vmOutputGenericResponse{}
	_ = doPost(t, facade, "/vm-values/query", &request, &response)
}

func doPost(t *testing.T, facade interface{}, url string, request interface{}, response interface{}) int {
	// Serialize if not already
	requestAsBytes, ok := request.([]byte)
	if !ok {
		requestAsBytes, _ = json.Marshal(request)
	}

	group, err := groups.NewVmValuesGroup(facade)
	require.NoError(t, err)
	server := startProxyServer(group, vmValuesPath)

	httpRequest, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestAsBytes))

	responseRecorder := httptest.NewRecorder()
	server.ServeHTTP(responseRecorder, httpRequest)

	loadResponse(responseRecorder.Body, &response)
	return responseRecorder.Code
}

func requireErrorOnAllRoutes(t *testing.T, facade interface{}, request interface{}, errExpected error) {
	requireErrorOnGetSingleValueRoutes(t, facade, request, errExpected)

	response := simpleResponse{}
	statusCode := doPost(t, facade, "/vm-values/query", request, &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Contains(t, response.Error, errExpected.Error())
}

func requireErrorOnGetSingleValueRoutes(t *testing.T, facade interface{}, request interface{}, errExpected error) {
	response := simpleResponse{}

	statusCode := doPost(t, facade, "/vm-values/hex", request, &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Contains(t, response.Error, errExpected.Error())

	statusCode = doPost(t, facade, "/vm-values/string", request, &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Contains(t, response.Error, errExpected.Error())

	statusCode = doPost(t, facade, "/vm-values/int", request, &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Contains(t, response.Error, errExpected.Error())
}

func createSCQuery(request *groups.VMValueRequest) (*data.SCQuery, error) {
	arguments := make([][]byte, len(request.Args))
	for i, arg := range request.Args {
		argBytes, err := hex.DecodeString(arg)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid hex string: %s", arg, err.Error())
		}

		arguments[i] = append(arguments[i], argBytes...)
	}

	return &data.SCQuery{
		ScAddress:      request.ScAddress,
		FuncName:       request.FuncName,
		CallerAddr:     request.CallerAddr,
		CallValue:      request.CallValue,
		ShouldBeSynced: request.ShouldBeSynced,
		SameScState:    request.SameScState,
		Arguments:      arguments,
	}, nil
}
