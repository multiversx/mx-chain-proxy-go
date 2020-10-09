package client

import (
	"encoding/json"
	"errors"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ElrondClient is able to process requests
type ElrondClient struct {
	client                    ElrondProxyClient
	genesisTime               uint64
	roundDurationMilliseconds uint64
}

const (
	MaxRetriesGetNetworkConfig = 20
	DelayBetweenRetries        = 5 * time.Second
)

var (
	log = logger.GetOrCreate("rosetta/client")
	// ErrInvalidElrondProxyHandler signals that provided elrond proxy handler is not a elrond proxy client
	ErrInvalidElrondProxyHandler = errors.New("invalid elrond proxy handler")
)

//NewElrondClient will create a new instance of ElrondClient
func NewElrondClient(elrondFacade api.ElrondProxyHandler) (*ElrondClient, error) {
	elrondProxy, ok := elrondFacade.(ElrondProxyClient)
	if !ok {
		return nil, ErrInvalidElrondProxyHandler
	}

	elrondClient := &ElrondClient{
		client: elrondProxy,
	}

	err := elrondClient.initializeElrondClient()
	if err != nil {
		return nil, err
	}

	return elrondClient, nil
}

func (ec *ElrondClient) initializeElrondClient() error {
	var err error

	networkConfig := &NetworkConfig{}
	for count := 0; count < MaxRetriesGetNetworkConfig; count++ {
		networkConfig, err = ec.GetNetworkConfig()
		if err != nil {
			time.Sleep(DelayBetweenRetries)
			continue
		}

		break
	}
	// if maxRetries is reached we should return error here because we did maxRetries to get network config
	// but the observers did not answer
	if err != nil {
		return err
	}

	ec.genesisTime = networkConfig.StartTime
	ec.roundDurationMilliseconds = networkConfig.RoundDuration

	return nil
}

// GetNetworkConfig will return the network config
func (ec *ElrondClient) GetNetworkConfig() (*NetworkConfig, error) {
	networkConfigResponse, err := ec.client.GetNetworkConfigMetrics()
	if err != nil {
		log.Warn("cannot get network metrics", "error", err.Error())

		return nil, err
	}

	if networkConfigResponse.Error != "" {
		log.Warn("cannot get network metrics", "error", networkConfigResponse.Error)

		return nil, errors.New(networkConfigResponse.Error)
	}

	networkConfig := &NetworkConfig{}

	responseDataI, ok := networkConfigResponse.Data.(map[string]interface{})
	if !ok {
		return nil, errors.New("response data is invalid")
	}
	responseData, ok := responseDataI["config"]
	if !ok {
		return nil, errors.New("response data is invalid network config is not in response")
	}

	responseBytes, err := json.Marshal(responseData)
	if err != nil {
		log.Warn("cannot marshal network config response", "error", err.Error())

		return nil, err
	}

	err = json.Unmarshal(responseBytes, networkConfig)

	return networkConfig, err
}

// GetLatestBlockData will return latest block data
func (ec *ElrondClient) GetLatestBlockData() (*BlockData, error) {
	latestBlockNonce, err := ec.client.GetLatestFullySynchronizedHyperblockNonce()
	if err != nil {
		return nil, err
	}

	blockResponse, err := ec.client.GetBlockByNonce(MetachainID, latestBlockNonce, false)
	if err != nil {
		log.Warn("cannot get block", "nonce", latestBlockNonce,
			"error", err.Error())

		return nil, err
	}

	if blockResponse.Error != "" {
		log.Warn("cannot get block", "nonce", latestBlockNonce,
			"error", blockResponse.Error)

		return nil, err
	}

	return &BlockData{
		Nonce:         blockResponse.Data.Block.Nonce,
		Hash:          blockResponse.Data.Block.Hash,
		PrevBlockHash: blockResponse.Data.Block.PrevBlockHash,
		Timestamp:     ec.CalculateBlockTimestampUnix(blockResponse.Data.Block.Round),
	}, nil
}

// GetBlockByNonce will return a block by nonce
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

		return nil, errors.New(blockResponse.Error)
	}

	return &blockResponse.Data.Hyperblock, nil
}

// GetBlockByHash will return a hyper block by hash
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

		return nil, errors.New(blockResponse.Error)
	}

	return &blockResponse.Data.Hyperblock, nil
}

// GetAccount will return an account by address
func (ec *ElrondClient) GetAccount(address string) (*data.Account, error) {
	return ec.client.GetAccount(address)
}

// ComputeTransactionHash will compute hash of provided transaction
func (ec *ElrondClient) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	return ec.client.ComputeTransactionHash(tx)
}

// EncodeAddress will encode an address
func (ec *ElrondClient) EncodeAddress(address []byte) (string, error) {
	pubKeyConverter, err := ec.client.GetAddressConverter()
	if err != nil {
		return "", err
	}

	return pubKeyConverter.Encode(address), nil
}

// SendTx will send a transaction
func (ec *ElrondClient) SendTx(tx *data.Transaction) (string, error) {
	_, hash, err := ec.client.SendTransaction(tx)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// CalculateBlockTimestampUnix will calculate block timestamp
func (ec *ElrondClient) CalculateBlockTimestampUnix(round uint64) int64 {
	startTimeMilliseconds := ec.genesisTime * 1000

	return int64(startTimeMilliseconds) + int64(round*ec.roundDurationMilliseconds)
}

// GetTransactionByHashFromPool will return a transaction only if is in pool
func (ec *ElrondClient) GetTransactionByHashFromPool(txHash string) (*data.FullTransaction, bool) {
	tx, _, err := ec.client.GetTransactionByHashAndSenderAddress(txHash, "")
	if err != nil {
		log.Debug("elrond clinet: cannot get transaction by hash", "error", err.Error())
		return nil, false
	}

	if !isTxFromPool(tx) {
		return nil, false
	}

	return tx, true
}

func isTxFromPool(tx *data.FullTransaction) bool {
	acceptedTxStatuses := []transaction.TxStatus{transaction.TxStatusReceived, transaction.TxStatusPartiallyExecuted}
	for idx := 0; idx < len(acceptedTxStatuses); idx++ {
		if acceptedTxStatuses[idx] == tx.Status {
			return true
		}
	}

	return false
}
