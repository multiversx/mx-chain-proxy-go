package process

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// SCQueryServicePath defines the get values path at which the nodes answer
const SCQueryServicePath = "/vm-values/query"

// SCQueryServiceProxy is able to process smart contract queries
type SCQueryServiceProxy struct {
	proc Processor
}

// NewSCQueryServiceProxy creates a new instance of GetValuesProcessor
func NewSCQueryServiceProxy(proc Processor) (*SCQueryServiceProxy, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &SCQueryServiceProxy{
		proc: proc,
	}, nil
}

// ExecuteQuery resolves the request by sending the request to the right observer and replies back the answer
func (proxy *SCQueryServiceProxy) ExecuteQuery(query *process.SCQuery) (*vmcommon.VMOutput, error) {
	addressBytes := query.ScAddress
	shardID, err := proxy.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	observers, err := proxy.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		request := proxy.createRequestFromQuery(query)
		response := &data.ResponseVmValue{}

		httpStatus, err := proxy.proc.CallPostRestEndPoint(observer.Address, SCQueryServicePath, request, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Info(fmt.Sprintf("SC query sent successfully to observer %v from shard %v, received response", observer.Address, shardID))
			return response.Data, nil
		}

		if responseHasExplicitError {
			return nil, fmt.Errorf(response.Error)
		}

		return nil, err
	}

	return nil, ErrSendingRequest
}

func (proxy *SCQueryServiceProxy) createRequestFromQuery(query *process.SCQuery) data.VmValueRequest {
	request := data.VmValueRequest{}
	request.Address = hex.EncodeToString(query.ScAddress)
	request.FuncName = query.FuncName
	request.Args = make([]string, len(query.Arguments))
	for i, argument := range query.Arguments {
		argumentAsHex := hex.EncodeToString(argument)
		request.Args[i] = argumentAsHex
	}

	return request
}
