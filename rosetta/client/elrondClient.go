package client

import (
	"encoding/json"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

var log = logger.GetOrCreate("rosetta/client")

type objectMap = map[string]interface{}

// ElrondClient -
type ElrondClient struct {
	client ElrondProxyClient
}

func NewElrondClient(elrondFacade api.ElrondProxyHandler) *ElrondClient {
	elrondProxy := elrondFacade.(ElrondProxyClient)

	return &ElrondClient{
		client: elrondProxy,
	}
}

func (ec *ElrondClient) GetNetworkConfig() (*NetworkConfig, error) {
	networkConfigResponse, err := ec.client.GetNetworkConfigMetrics()
	if err != nil {
		log.Warn("cannot get network metrics", "error", err.Error())

		return nil, err
	}

	if networkConfigResponse.Error != "" {
		log.Warn("cannot get network metrics", "error", networkConfigResponse.Error)

		return nil, err
	}

	networkConfig := &NetworkConfig{}
	responseBytes, _ := json.Marshal(networkConfigResponse.Data.(objectMap)["config"])
	err = json.Unmarshal(responseBytes, networkConfig)

	return networkConfig, err
}

func (ec *ElrondClient) GetNetworkStatus() (*NetworkStatus, error) {
	networkStatusResponse, err := ec.client.GetNetworkStatusMetrics(MetachainID)
	if err != nil {
		log.Warn("cannot get network status", "error", err.Error())

		return nil, err
	}
	if networkStatusResponse.Error != "" {
		log.Warn("cannot get network status", "error", networkStatusResponse.Error)

		return nil, err
	}

	networkStatus := &NetworkStatus{}
	responseBytes, _ := json.Marshal(networkStatusResponse.Data.(objectMap)["status"])
	err = json.Unmarshal(responseBytes, networkStatus)

	return networkStatus, err
}

func (ec *ElrondClient) GetLatestBlockData() (*BlockData, error) {
	networkStatus, err := ec.GetNetworkStatus()
	if err != nil {
		return nil, err
	}

	blockResponse, err := ec.client.GetBlockByNonce(MetachainID, networkStatus.CurrentNonce, false)
	if err != nil {
		log.Warn("cannot get block", "nonce", networkStatus.CurrentNonce,
			"error", err.Error())

		return nil, err
	}

	if blockResponse.Error != "" {
		log.Warn("cannot get block", "nonce", networkStatus.CurrentNonce,
			"error", blockResponse.Error)

		return nil, err
	}

	return &BlockData{
		Nonce:         blockResponse.Data.Block.Nonce,
		Hash:          blockResponse.Data.Block.Hash,
		PrevBlockHash: blockResponse.Data.Block.PrevBlockHash,
		Timestamp:     CalculateBlockTimestampUnix(blockResponse.Data.Block.Round),
	}, nil
}

func (ec *ElrondClient) GetBlockByNonce(nonce int64) (*data.Hyperblock, error) {
	blockResponse, err := ec.client.GetHyperBlockByNonce(uint64(nonce))
	if err != nil {
		log.Warn("cannot get hyper block", "nonce", nonce,
			"error", err.Error())

		return nil, err
	}

	if blockResponse.Error != "" {
		log.Warn("cannot get hyper block", "nonce", nonce,
			"error", blockResponse.Error)

		return nil, err
	}

	return &blockResponse.Data.Hyperblock, nil
}

func (ec *ElrondClient) GetBlockByHash(hash string) (*data.Hyperblock, error) {
	blockResponse, err := ec.client.GetHyperBlockByHash(hash)
	if err != nil {
		log.Warn("cannot get hyper block", "hash", hash,
			"error", err.Error())

		return nil, err
	}

	if blockResponse.Error != "" {
		log.Warn("cannot get  hyper block", "hash", hash,
			"error", blockResponse.Error)

		return nil, err
	}

	return &blockResponse.Data.Hyperblock, nil
}

func (ec *ElrondClient) GetAccount(address string) (*data.Account, error) {
	return ec.client.GetAccount(address)
}

func (ec *ElrondClient) SimulateTx(tx *data.Transaction) (string, error) {
	simulatedTxResponse, err := ec.client.SimulateTransaction(tx)
	if err != nil {
		return "", err
	}

	if simulatedTxResponse.Error != "" {
		log.Warn("cannot simulate", "error", simulatedTxResponse.Error)

		return "", err
	}

	return simulatedTxResponse.Data.Result.Hash, nil
}

func (ec *ElrondClient) SendTx(tx *data.Transaction) (string, error) {
	_, hash, err := ec.client.SendTransaction(tx)
	if err != nil {
		return "", err
	}

	return hash, nil
}
