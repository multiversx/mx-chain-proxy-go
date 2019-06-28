package process

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GetValuesPath defines the get values path at which the nodes answer
const GetValuesPath = "/get-values/"

// GetValuesProcessor is able to process get values requests
type GetValuesProcessor struct {
	proc Processor
}

// NewGetValuesProcessor creates a new instance of GetValuesProcessor
func NewGetValuesProcessor(proc Processor) (*GetValuesProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &GetValuesProcessor{
		proc: proc,
	}, nil
}

// GetDataValue resolves the request by sending the request to the right observer and replies back the answer
func (gvp *GetValuesProcessor) GetDataValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	addressBytes, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	shardId, err := gvp.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	observers, err := gvp.proc.GetObservers(shardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		gvr := &data.GetValuesRequest{}
		gvr.Address = address
		gvr.FuncName = funcName

		hexArgs := make([]string, len(argsBuff))
		for idx, arg := range argsBuff {
			argHex := hex.EncodeToString(arg)
			hexArgs[idx] = argHex
		}
		gvr.Args = hexArgs

		getValuesResponse := &data.ResponseGetValues{}

		err = gvp.proc.CallPostRestEndPoint(observer.Address, GetValuesPath, gvr, getValuesResponse)
		if err == nil {
			log.Info(fmt.Sprintf("GetValues sent successfully to observer %v from shard %v, received value %s",
				observer.Address,
				shardId,
				getValuesResponse.HexData,
			))

			getValBytes, err := hex.DecodeString(getValuesResponse.HexData)
			if err != nil {
				log.LogIfError(err)
				//we move to the next observer. It might give a good answer
				continue
			}

			return getValBytes, nil
		}

		log.LogIfError(err)
	}

	return nil, ErrSendingRequest
}
