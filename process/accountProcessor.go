package process

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AddressPath defines the address path at which the nodes answer
const AddressPath = "/address/"

// AccountProcessor is able to process account requests
type AccountProcessor struct {
	proc Processor
}

// NewAccountProcessor creates a new instance of AccountProcessor
func NewAccountProcessor(proc Processor) (*AccountProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}

	return &AccountProcessor{
		proc: proc,
	}, nil
}

// GetAccount resolves the request by sending the request to the right observer and replies back the answer
func (ap *AccountProcessor) GetAccount(address string) (*data.Account, error) {
	addressBytes, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	shardId, err := ap.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	observers, err := ap.proc.GetObservers(shardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		responseAccount := &data.ResponseAccount{}

		err = ap.proc.CallGetRestEndPoint(observer.Address, AddressPath+address, responseAccount)
		if err == nil {
			log.Info(fmt.Sprintf("Got account request from observer %v from shard %v", observer.Address, shardId))
			return &responseAccount.AccountData, nil
		}

		log.LogIfError(err)
	}

	return nil, ErrSendingRequest
}
