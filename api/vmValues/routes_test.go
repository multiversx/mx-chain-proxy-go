package vmValues_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	apiErrors "github.com/ElrondNetwork/elrond-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/vmValues"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type GeneralResponse struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

func init() {
	gin.SetMode(gin.TestMode)
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		logError(err)
	}
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func startNodeServer(handler vmValues.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	getValuesRoute := ws.Group("/vm-values")

	if handler != nil {
		getValuesRoute.Use(api.WithElrondProxyFacade(handler))
	}
	vmValues.Routes(getValuesRoute)

	return ws
}

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	getValuesRoute := ws.Group("/vm-values")
	vmValues.Routes(getValuesRoute)

	return ws
}

//------- GetDataValueAsHexBytes

func TestGetDataValueAsHexBytes_WithWrongFacadeShouldErr(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()

	jsonStr := `{"scAddress":"DEADBEEF","funcName":"DEADBEEF","args":[]}`
	req, _ := http.NewRequest("POST", "/vm-values/hex", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Contains(t, response.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetDataValueAsHexBytes_BadRequestShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetVmValueHandler: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
			assert.Fail(t, "should have not called this")
			return nil, nil
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := `{"this should error"}`
	req, _ := http.NewRequest("POST", "/vm-values/hex", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Contains(t, response.Error, "invalid character")
}

func TestGetDataValueAsHexBytes_ArgumentIsNotHexShouldErr(t *testing.T) {
	t.Parallel()

	scAddress := "sc address"
	fName := "function"
	args := []string{"not a hex argument"}
	errUnexpected := errors.New("unexpected error")
	valueBuff, _ := hex.DecodeString("DEADBEEF")

	facade := mock.Facade{
		GetVmValueHandler: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
			if address == scAddress && funcName == fName && len(argsBuff) == len(args) {
				paramsOk := true
				for idx, arg := range args {
					if arg != string(argsBuff[idx]) {
						paramsOk = false
					}
				}

				if paramsOk {
					return valueBuff, nil
				}
			}

			return nil, errUnexpected
		},
	}

	ws := startNodeServer(&facade)

	argsJson, _ := json.Marshal(args)

	jsonStr := fmt.Sprintf(
		`{"scAddress":"%s", "funcName":"%s", "args":%s}`,
		scAddress,
		fName,
		argsJson)
	fmt.Printf("Request: %s\n", jsonStr)

	req, _ := http.NewRequest("POST", "/vm-values/hex", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, "not a hex argument")
}

func testGetValueFacadeErrors(t *testing.T, route string) {
	t.Parallel()

	errExpected := errors.New("expected error")
	facade := mock.Facade{
		GetVmValueHandler: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
			return nil, errExpected
		},
	}

	ws := startNodeServer(&facade)

	jsonStr := `{}`
	fmt.Printf("Request: %s\n", jsonStr)

	req, _ := http.NewRequest("POST", route, bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, errExpected.Error())
}

func TestGetDataValueAsHexBytes_FacadeErrorsShouldErr(t *testing.T) {
	testGetValueFacadeErrors(t, "/vm-values/hex")
}

func TestGetDataValueAsHexBytes_WithParametersShouldReturnValueAsHex(t *testing.T) {
	t.Parallel()

	scAddress := "sc address"
	fName := "function"
	args := []string{"argument 1", "argument 2"}
	errUnexpected := errors.New("unexpected error")
	valueBuff, _ := hex.DecodeString("DEADBEEF")

	facade := mock.Facade{
		GetVmValueHandler: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
			if address == scAddress && funcName == fName && len(argsBuff) == len(args) {
				paramsOk := true
				for idx, arg := range args {
					if arg != string(argsBuff[idx]) {
						paramsOk = false
					}
				}

				if paramsOk {
					return valueBuff, nil
				}
			}

			return nil, errUnexpected
		},
	}

	ws := startNodeServer(&facade)

	argsHex := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		argsHex[i] = hex.EncodeToString([]byte(args[i]))
	}
	argsJson, _ := json.Marshal(argsHex)

	jsonStr := fmt.Sprintf(
		`{"scAddress":"%s", "funcName":"%s", "args":%s}`,
		scAddress,
		fName,
		argsJson)
	fmt.Printf("Request: %s\n", jsonStr)

	req, _ := http.NewRequest("POST", "/vm-values/hex", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", response.Error)
	assert.Equal(t, hex.EncodeToString(valueBuff), response.Data)
}

//------- GetDataValueAsString

func TestGetDataValueAsString_FacadeErrorsShouldErr(t *testing.T) {
	testGetValueFacadeErrors(t, "/vm-values/string")
}

func TestGetDataValueAsString_WithParametersShouldReturnValueAsHex(t *testing.T) {
	t.Parallel()

	scAddress := "sc address"
	fName := "function"
	args := []string{"argument 1", "argument 2"}
	errUnexpected := errors.New("unexpected error")
	valueBuff := "DEADBEEF"

	facade := mock.Facade{
		GetVmValueHandler: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
			if address == scAddress && funcName == fName && len(argsBuff) == len(args) {
				paramsOk := true
				for idx, arg := range args {
					if arg != string(argsBuff[idx]) {
						paramsOk = false
					}
				}

				if paramsOk {
					return []byte(valueBuff), nil
				}
			}

			return nil, errUnexpected
		},
	}

	ws := startNodeServer(&facade)

	argsHex := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		argsHex[i] = hex.EncodeToString([]byte(args[i]))
	}
	argsJson, _ := json.Marshal(argsHex)

	jsonStr := fmt.Sprintf(
		`{"scAddress":"%s", "funcName":"%s", "args":%s}`,
		scAddress,
		fName,
		argsJson)
	fmt.Printf("Request: %s\n", jsonStr)

	req, _ := http.NewRequest("POST", "/vm-values/string", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", response.Error)
	assert.Equal(t, valueBuff, response.Data)
}

//------- GetDataValueAsInt

func TestGetDataValueAsInt_FacadeErrorsShouldErr(t *testing.T) {
	testGetValueFacadeErrors(t, "/vm-values/int")
}

func TestGetDataValueAsInt_WithParametersShouldReturnValueAsHex(t *testing.T) {
	t.Parallel()

	scAddress := "sc address"
	fName := "function"
	args := []string{"argument 1", "argument 2"}
	errUnexpected := errors.New("unexpected error")
	valueBuff := "1234567"

	facade := mock.Facade{
		GetVmValueHandler: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
			if address == scAddress && funcName == fName && len(argsBuff) == len(args) {
				paramsOk := true
				for idx, arg := range args {
					if arg != string(argsBuff[idx]) {
						paramsOk = false
					}
				}

				if paramsOk {
					val := big.NewInt(0)
					val.SetString(valueBuff, 10)
					return val.Bytes(), nil
				}
			}

			return nil, errUnexpected
		},
	}

	ws := startNodeServer(&facade)

	argsHex := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		argsHex[i] = hex.EncodeToString([]byte(args[i]))
	}
	argsJson, _ := json.Marshal(argsHex)

	jsonStr := fmt.Sprintf(
		`{"scAddress":"%s", "funcName":"%s", "args":%s}`,
		scAddress,
		fName,
		argsJson)
	fmt.Printf("Request: %s\n", jsonStr)

	req, _ := http.NewRequest("POST", "/vm-values/int", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "", response.Error)
	assert.Equal(t, valueBuff, response.Data)
}