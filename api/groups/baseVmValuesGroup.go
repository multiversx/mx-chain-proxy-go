package groups

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// VMValueRequest represents the structure on which user input for generating a new transaction will validate against
type VMValueRequest struct {
	ScAddress  string   `form:"scAddress" json:"scAddress"`
	FuncName   string   `form:"funcName" json:"funcName"`
	CallerAddr string   `form:"caller" json:"caller"`
	CallValue  string   `form:"value" json:"value"`
	Args       []string `form:"args"  json:"args"`
}

type vmValuesGroup struct {
	facade VmValuesFacadeHandler
	*baseGroup
}

// NewVmValuesGroup returns a new instance of vmValuesGroup
func NewVmValuesGroup(facadeHandler data.FacadeHandler) (*vmValuesGroup, error) {
	facade, ok := facadeHandler.(VmValuesFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	vvg := &vmValuesGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/hex":    {Handler: vvg.getHex, Method: http.MethodPost},
		"/string": {Handler: vvg.getString, Method: http.MethodPost},
		"/int":    {Handler: vvg.getInt, Method: http.MethodPost},
		"/query":  {Handler: vvg.executeQuery, Method: http.MethodPost},
	}
	vvg.baseGroup.endpoints = baseRoutesHandlers

	return vvg, nil
}

// getHex returns the data as bytes, hex-encoded
func (group *vmValuesGroup) getHex(context *gin.Context) {
	group.doGetVMValue(context, vmcommon.AsHex)
}

// getString returns the data as string
func (group *vmValuesGroup) getString(context *gin.Context) {
	group.doGetVMValue(context, vmcommon.AsString)
}

// getInt returns the data as big int
func (group *vmValuesGroup) getInt(context *gin.Context) {
	group.doGetVMValue(context, vmcommon.AsBigIntString)
}

func (group *vmValuesGroup) doGetVMValue(context *gin.Context, asType vmcommon.ReturnDataKind) {
	vmOutput, status, err := group.doExecuteQuery(context)
	if err != nil {
		message := fmt.Sprintf("%s: %s", "doGetVMValue", err)
		shared.RespondWith(context, status, nil, message)
		return
	}

	returnData, err := vmOutput.GetFirstReturnData(asType)
	if err != nil {
		message := fmt.Sprintf("%s: %s", "doGetVMValue", err)
		shared.RespondWith(context, http.StatusBadRequest, nil, message)
		return
	}

	returnOkResponse(context, returnData)
}

// executeQuery returns the data as string
func (group *vmValuesGroup) executeQuery(context *gin.Context) {
	vmOutput, status, err := group.doExecuteQuery(context)
	if err != nil {
		message := fmt.Sprintf("%s: %s", "executeQuery", err)
		shared.RespondWith(context, status, nil, message)
		return
	}

	returnOkResponse(context, vmOutput)
}

func (group *vmValuesGroup) doExecuteQuery(context *gin.Context) (*vm.VMOutputApi, int, error) {
	request := VMValueRequest{}
	err := context.ShouldBindJSON(&request)
	if err != nil {
		return nil, http.StatusBadRequest, apiErrors.ErrInvalidJSONRequest
	}

	command, err := createSCQuery(&request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	vmOutput, status, err := group.facade.ExecuteSCQuery(command)
	if err != nil {
		return nil, status, err
	}

	return vmOutput, status, nil
}

func createSCQuery(request *VMValueRequest) (*data.SCQuery, error) {
	arguments := make([][]byte, len(request.Args))
	for i, arg := range request.Args {
		argBytes, err := hex.DecodeString(arg)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid hex string: %s", arg, err.Error())
		}

		arguments[i] = append(arguments[i], argBytes...)
	}

	return &data.SCQuery{
		ScAddress:  request.ScAddress,
		FuncName:   request.FuncName,
		CallerAddr: request.CallerAddr,
		CallValue:  request.CallValue,
		Arguments:  arguments,
	}, nil
}

func returnOkResponse(context *gin.Context, dataToReturn interface{}) {
	shared.RespondWith(context, http.StatusOK, gin.H{"data": dataToReturn}, "")
}
