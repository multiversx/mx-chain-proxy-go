package process

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AddressPath defines the address path at which the nodes answer
const AddressPath = "/address/"

// AccountProcessor is able to process account requests
type AccountProcessor struct {
	proc            Processor
	pubKeyConverter state.PubkeyConverter
}

// NewAccountProcessor creates a new instance of AccountProcessor
func NewAccountProcessor(proc Processor, pubKeyConverter state.PubkeyConverter) (*AccountProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}
	if pubKeyConverter == nil {
		return nil, ErrNilPubKeyConverter
	}

	return &AccountProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// GetAccount resolves the request by sending the request to the right observer and replies back the answer
func (ap *AccountProcessor) GetAccount(address string) (*data.Account, error) {
	addressBytes, err := ap.pubKeyConverter.Decode(address)
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
			log.Info("account request", "address", address, "shard id", shardId, "observer", observer.Address)
			return &responseAccount.AccountData, nil
		}

		log.Error("account request", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, ErrSendingRequest
}
