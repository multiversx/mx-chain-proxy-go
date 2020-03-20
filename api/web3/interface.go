package web3

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	PrepareDataForRequest(requestBody data.RequestBodyWeb3) (data.ResponseWeb3, error)
}
