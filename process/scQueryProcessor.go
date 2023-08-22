package process

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// scQueryServicePath defines the get values path at which the nodes answer
const scQueryServicePath = "/vm-values/query"
const blockNonceSuffix = "?blockNonce="
const blockHashSuffix = "?blockHash="

// SCQueryProcessor is able to process smart contract queries
type SCQueryProcessor struct {
	proc            Processor
	pubKeyConverter core.PubkeyConverter
}

// NewSCQueryProcessor creates a new instance of SCQueryProcessor
func NewSCQueryProcessor(proc Processor, pubKeyConverter core.PubkeyConverter) (*SCQueryProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &SCQueryProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// ExecuteQuery resolves the request by sending the request to the right observer and replies back the answer
func (scQueryProcessor *SCQueryProcessor) ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	addressBytes, err := scQueryProcessor.pubKeyConverter.Decode(query.ScAddress)
	if err != nil {
		return nil, data.BlockInfo{}, err
	}

	shardID, err := scQueryProcessor.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, data.BlockInfo{}, err
	}

	observers, err := scQueryProcessor.proc.GetObservers(shardID)
	if err != nil {
		return nil, data.BlockInfo{}, err
	}

	for _, observer := range observers {
		request := scQueryProcessor.createRequestFromQuery(query)
		response := &data.ResponseVmValue{}

		path := scQueryServicePath
		if len(query.BlockHash) > 0 {
			path = fmt.Sprintf("%s%s%s", path, blockHashSuffix, hex.EncodeToString(query.BlockHash))
		}
		if query.BlockNonce.HasValue {
			path = fmt.Sprintf("%s%s%d", path, blockNonceSuffix, query.BlockNonce.Value)
		}

		httpStatus, err := scQueryProcessor.proc.CallPostRestEndPoint(observer.Address, path, request, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Debug("SC query sent successfully, received response", "observer", observer.Address, "shard", shardID)
			return response.Data.Data, response.Data.BlockInfo, nil
		}

		if responseHasExplicitError {
			return nil, data.BlockInfo{}, fmt.Errorf(response.Error)
		}

		return nil, data.BlockInfo{}, err
	}

	return nil, data.BlockInfo{}, ErrSendingRequest
}

func (scQueryProcessor *SCQueryProcessor) createRequestFromQuery(query *data.SCQuery) data.VmValueRequest {
	request := data.VmValueRequest{}
	request.Address = query.ScAddress
	request.FuncName = query.FuncName
	request.CallValue = query.CallValue
	request.CallerAddr = query.CallerAddr
	request.Args = make([]string, len(query.Arguments))
	for i, argument := range query.Arguments {
		argumentAsHex := hex.EncodeToString(argument)
		request.Args[i] = argumentAsHex
	}

	return request
}

// IsInterfaceNil returns true if the value under the interface is nil
func (scQueryProcessor *SCQueryProcessor) IsInterfaceNil() bool {
	return scQueryProcessor == nil
}
