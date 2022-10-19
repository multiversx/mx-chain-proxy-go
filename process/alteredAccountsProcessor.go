package process

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
)

const (
	alteredAccountByBlockNonce = "altered-accounts/by-nonce"
)

type alteredAccountsProcessor struct {
	proc Processor
}

func NewAlteredAccountsProcessor(proc Processor) (*alteredAccountsProcessor, error) {
	return &alteredAccountsProcessor{
		proc: proc,
	}, nil
}

func (aap *alteredAccountsProcessor) GetAlteredAccountsByNonce(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) ([]*outport.AlteredAccount, error) {
	observers, err := aap.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}
	path := common.BuildUrlWithAlteredAccountsQueryOptions(fmt.Sprintf("%s/%d", alteredAccountByBlockNonce, nonce), options)

	for _, observer := range observers {
		var response []*outport.AlteredAccount

		_, err := aap.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("altered accounts request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("altered accounts request", "shard id", observer.ShardId, "nonce", nonce, "observer", observer.Address)
		return response, nil

	}

	return nil, ErrSendingRequest
}

func (aap *alteredAccountsProcessor) GetAlteredAccountsByHash(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) ([]*outport.AlteredAccount, error) {
	return nil, nil
}
