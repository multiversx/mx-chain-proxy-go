package process

import (
	"encoding/hex"
	"fmt"

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
func (gvp *VmValuesProcessor) GetVmValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
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

		err = gvp.proc.CallPostRestEndPoint(observer.Address, GetValuesPath, vvr, vmValueResponse)
		if err == nil {
			log.Info(fmt.Sprintf("VmValues sent successfully to observer %v from shard %v, received value %s",
				observer.Address,
				shardId,
				vmValueResponse.HexData,
			))

			getValBytes, err := hex.DecodeString(vmValueResponse.HexData)
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
