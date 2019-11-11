package process

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GetValuesPath defines the get values path at which the nodes answer
const GetValuesPath = "/vm-values/"

// VmValuesProcessor is able to process get values requests
type VmValuesProcessor struct {
	proc Processor
}

// NewVmValuesProcessor creates a new instance of GetValuesProcessor
func NewVmValuesProcessor(proc Processor) (*VmValuesProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &VmValuesProcessor{
		proc: proc,
	}, nil
}

// GetVmValue resolves the request by sending the request to the right observer and replies back the answer
func (gvp *VmValuesProcessor) GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
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
		vvr := &data.VmValueRequest{}
		vvr.Address = address
		vvr.FuncName = funcName

		hexArgs := make([]string, len(argsBuff))
		for idx, arg := range argsBuff {
			argHex := hex.EncodeToString(arg)
			hexArgs[idx] = argHex
		}
		vvr.Args = hexArgs

		vmValueResponse := &data.ResponseVmValue{}

		respCode, err := gvp.proc.CallPostRestEndPoint(observer.Address, GetValuesPath+resType, vvr, vmValueResponse)
		if respCode == http.StatusOK && err == nil {
			log.Info(fmt.Sprintf("VmValues sent successfully to observer %v from shard %v, received value %s",
				observer.Address,
				shardId,
				vmValueResponse.HexData,
			))

			return []byte(vmValueResponse.HexData), nil
		}

		// if observer was down (or didn't respond in time), skip to the next one
		if respCode == http.StatusNotFound || respCode == http.StatusRequestTimeout {
			log.LogIfError(err)
			continue
		}

		// if the request was bad, return the error message
		return nil, err
	}

	return nil, ErrSendingRequest
}
