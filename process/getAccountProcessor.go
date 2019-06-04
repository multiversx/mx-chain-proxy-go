package process

import (
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GetAccountProcessor is able to process account requests
type GetAccountProcessor struct {
}

// GetAccount resolves the request by sending the request to the right observer and replies back the answer
func (gap *GetAccountProcessor) GetAccount(address string) (*data.Account, error) {
	//TODO, fix this mock with a real impl

	return &data.Account{
		Nonce:    1,
		Balance:  "2",
		Address:  address,
		RootHash: []byte("ROOT_HASH"),
		CodeHash: []byte("CODE_HASH"),
	}, nil
}

func (gap *GetAccountProcessor) ApplyConfig(cfg *config.Config) {

}
