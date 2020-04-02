package process

import (
	"encoding/hex"

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
			log.Info("account request", "address", address, "shard id", shardId, "observer", observer.Address)
			return &responseAccount.AccountData, nil
		}

		log.Error("account request", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, ErrSendingRequest
}

//// ValStatsResponse respects the format the validator statistics are received from the observers
//type ValStatsResponse struct {
//	Statistics map[string]*data.ValidatorApiResponse `json:"statistics"`
//}
//
//// ValidatorStatistics will fetch from the observers details about validators statistics
//func (ap *AccountProcessor) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
//	observers, err := ap.proc.GetObservers(core.MetachainShardId)
//	if err != nil {
//		return nil, err
//	}
//
//	valStatsMap := &ValStatsResponse{}
//
//	for _, observer := range observers {
//		err = ap.proc.CallGetRestEndPoint(observer.Address, ValidatorStatisticsPath, valStatsMap)
//		if err == nil {
//			log.Info("validator statistics request", "observer", observer.Address)
//			return valStatsMap.Statistics, nil
//		}
//
//		log.Error("validator statistics request", "observer", observer.Address, "error", err.Error())
//	}
//
//	return nil, ErrSendingRequest
//}
