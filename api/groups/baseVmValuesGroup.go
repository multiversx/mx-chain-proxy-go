package groups

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// VMValueRequest represents the structure on which user input for generating a new transaction will validate against
type VMValueRequest struct {
	ScAddress      string   `json:"scAddress"`
	FuncName       string   `json:"funcName"`
	CallerAddr     string   `json:"caller"`
	CallValue      string   `json:"value"`
	SameScState    bool     `json:"sameScState"`
	ShouldBeSynced bool     `json:"shouldBeSynced"`
	Args           []string `json:"args"`
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

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/hex", Handler: vvg.getHex, Method: http.MethodPost},
		{Path: "/string", Handler: vvg.getString, Method: http.MethodPost},
		{Path: "/int", Handler: vvg.getInt, Method: http.MethodPost},
		{Path: "/query", Handler: vvg.executeQuery, Method: http.MethodPost},
	}
	vvg.baseGroup.endpoints = baseRoutesHandlers

	return vvg, nil
}

// getHex returns the data as bytes, hex-encoded
func (group *vmValuesGroup) getHex(context *gin.Context) {
	group.doGetVMValue(context, vm.AsHex)
}

// getString returns the data as string
func (group *vmValuesGroup) getString(context *gin.Context) {
	group.doGetVMValue(context, vm.AsString)
}

// getInt returns the data as big int
func (group *vmValuesGroup) getInt(context *gin.Context) {
	group.doGetVMValue(context, vm.AsBigIntString)
}

func (group *vmValuesGroup) doGetVMValue(context *gin.Context, asType vm.ReturnDataKind) {
	vmOutput, blockInfo, err := group.doExecuteQuery(context)

	if err != nil {
		returnBadRequest(context, "doGetVMValue", err)
		return
	}

	returnData, err := vmOutput.GetFirstReturnData(vm.ReturnDataKind(asType))
	if err != nil {
		returnBadRequest(context, "doGetVMValue", err)
		return
	}

	returnOkResponse(context, returnData, blockInfo)
}

// executeQuery returns the data as string
func (group *vmValuesGroup) executeQuery(context *gin.Context) {
	vmOutput, blockInfo, err := group.doExecuteQuery(context)
	if err != nil {
		returnBadRequest(context, "executeQuery", err)
		return
	}

	returnOkResponse(context, vmOutput, blockInfo)
}

func (group *vmValuesGroup) doExecuteQuery(context *gin.Context) (*vm.VMOutputApi, data.BlockInfo, error) {
	request := VMValueRequest{}
	err := context.ShouldBindJSON(&request)
	if err != nil {
		return nil, data.BlockInfo{}, apiErrors.ErrInvalidJSONRequest
	}

	command, err := createSCQuery(&request)
	if err != nil {
		return nil, data.BlockInfo{}, err
	}

	command.BlockNonce, command.BlockHash, err = extractBlockCoordinates(context)
	if err != nil {
		return nil, data.BlockInfo{}, err
	}

	vmOutput, blockInfo, err := group.facade.ExecuteSCQuery(command)
	if err != nil {
		return nil, data.BlockInfo{}, err
	}

	return vmOutput, blockInfo, nil
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
		ScAddress:      request.ScAddress,
		FuncName:       request.FuncName,
		CallerAddr:     request.CallerAddr,
		CallValue:      request.CallValue,
		SameScState:    request.SameScState,
		ShouldBeSynced: request.ShouldBeSynced,
		Arguments:      arguments,
	}, nil
}

func extractBlockCoordinates(context *gin.Context) (core.OptionalUint64, []byte, error) {
	blockNonce, err := parseUint64UrlParam(context, common.UrlParameterBlockNonce)
	if err != nil {
		return core.OptionalUint64{}, nil, fmt.Errorf("%w for block nonce", err)
	}

	blockHash, err := parseHexBytesUrlParam(context, common.UrlParameterBlockHash)
	if err != nil {
		return core.OptionalUint64{}, nil, fmt.Errorf("%w for block hash", err)
	}

	return blockNonce, blockHash, nil
}

func returnBadRequest(context *gin.Context, errScope string, err error) {
	message := fmt.Sprintf("%s: %s", errScope, err)
	shared.RespondWith(context, http.StatusBadRequest, nil, message, data.ReturnCodeRequestError)
}

func returnOkResponse(context *gin.Context, dataToReturn interface{}, blockInfo interface{}) {
	shared.RespondWith(context, http.StatusOK, gin.H{"data": dataToReturn, "blockInfo": blockInfo}, "", data.ReturnCodeSuccess)
}
