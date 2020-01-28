package process

import (
	"encoding/hex"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AddressPath defines the address path at which the nodes answer
const AddressPath = "/address/"

// ValidatorStatisticsPath defines the validator statistics path at which the nodes answer
const ValidatorStatisticsPath = "/validator/statistics"

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

	if ap.proc.AreObserversBalanced() {
		observersRing, err := ap.proc.GetObserversRing(shardId)
		if err != nil {
			return nil, err
		}

		numTries := 0
		for numTries < observersRing.Len() {
			acc, err := ap.callApiEndpointForAccount(observersRing.Next(), address, shardId)
			if err == nil {
				return acc, nil
			}
			numTries++
		}
		return nil, ErrSendingRequest
	} else {
		observers, err := ap.proc.GetObservers(shardId)
		if err != nil {
			return nil, err
		}

		for _, observer := range observers {
			acc, err := ap.callApiEndpointForAccount(observer.Address, address, shardId)
			if err == nil {
				return acc, nil
			}
		}
		return nil, ErrSendingRequest
	}
}

func (ap *AccountProcessor) callApiEndpointForAccount(observerAddress string, address string, shardId uint32) (*data.Account, error) {
	responseAccount := &data.ResponseAccount{}

	err := ap.proc.CallGetRestEndPoint(observerAddress, AddressPath+address, responseAccount)
	if err == nil {
		log.Info("account request", "address", address, "shard id", shardId, "observer", observerAddress)
		return &responseAccount.AccountData, nil
	}

	log.Error("account request", "observer", observerAddress, "address", address, "error", err.Error())

	return nil, ErrNoResponseFromObserver
}

// ValStatsResponse respects the format the validator statistics are received from the observers
type ValStatsResponse struct {
	Statistics map[string]*data.ValidatorApiResponse `json:"statistics"`
}

// ValidatorStatistics will fetch from the observers details about validators statistics
func (ap *AccountProcessor) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	if ap.proc.AreObserversBalanced() {
		observersRing := ap.proc.GetAllObserversRing()
		numTries := 0
		for numTries < observersRing.Len() {
			retMap, err := ap.callApiEndpointForValidatorStatistics(observersRing.Next())
			if err == nil {
				return retMap, nil
			}
			numTries++
		}

		return nil, ErrSendingRequest
	} else {
		observers, err := ap.proc.GetAllObservers()
		if err != nil {
			return nil, err
		}

		for _, observer := range observers {
			retMap, err := ap.callApiEndpointForValidatorStatistics(observer.Address)
			if err == nil {
				return retMap, nil
			}
		}

		return nil, ErrSendingRequest
	}
}

func (ap *AccountProcessor) callApiEndpointForValidatorStatistics(observerAddress string) (map[string]*data.ValidatorApiResponse, error) {
	valStatsMap := &ValStatsResponse{}

	err := ap.proc.CallGetRestEndPoint(observerAddress, ValidatorStatisticsPath, valStatsMap)
	if err == nil {
		log.Info("validator statistics request", "observer", observerAddress)
		return valStatsMap.Statistics, nil
	}

	log.Error("validator statistics request", "observer", observerAddress, "error", err.Error())

	return nil, ErrNoResponseFromObserver
}
