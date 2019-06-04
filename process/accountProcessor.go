package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AddressPath defines the address path at which the nodes answer
const AddressPath = "/address/"

// AccountProcessor is able to process account requests
type AccountProcessor struct {
	proc CoreProcessor
}

// NewAccountProcessor creates a new instance of AccountProcessor
func NewAccountProcessor(proc CoreProcessor) (*AccountProcessor, error) {
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
		account := &data.Account{}
		err = ap.proc.CallGetRestEndPoint(observer.Address, AddressPath+address, account)
		if err == nil {
			return account, nil
		}

		log.LogIfError(err)
	}

	return nil, ErrSendingRequest
}
