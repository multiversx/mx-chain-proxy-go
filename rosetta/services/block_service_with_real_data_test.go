package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func TestBlockAPIService_GetBlockByNonce_ShouldWorkWithRealWorldData(t *testing.T) {
	t.Skip("Will be enabled upon re-fetching test data from the localnet.")
	t.Parallel()

	startNonce := int64(46)
	stopNonce := int64(56)

	for nonce := startNonce; nonce < stopNonce; nonce++ {
		checkBlock(t, nonce)
	}
}

func checkBlock(t *testing.T, nonce int64) {
	fmt.Printf("checkBlock(%d)\n", nonce)

	service := createService()

	expectedBlockResponse, err := readRosettaBlock(nonce)
	require.Nil(t, err)
	require.NotNil(t, expectedBlockResponse)

	actualBlockResponse, typedError := service.getBlockByNonce(nonce)
	require.Nil(t, typedError)
	require.NotNil(t, actualBlockResponse)

	expectedJson, _ := marshalPrettyJson(expectedBlockResponse)
	actualJson, _ := marshalPrettyJson(actualBlockResponse)

	require.Equal(t, expectedJson, actualJson, fmt.Sprintf("check failed: nonce = %d", nonce))
}

func createService() *blockAPIService {
	networkConfig := &provider.NetworkConfig{
		ChainID:        "localnet",
		GasPerDataByte: 1500,
		MinGasPrice:    1000000000,
		MinGasLimit:    50000,
		StartTime:      1653209603,
		RoundDuration:  6000,
	}

	configuration := &configuration.Configuration{
		ElrondNetworkConfig: networkConfig,
		Currency: &types.Currency{
			Symbol:   "XeGLD",
			Decimals: 18,
		},
	}

	providerMock := &mocks.ElrondProviderMock{
		GetBlockByNonceCalled: func(nonce int64) (*data.Hyperblock, error) {
			response, err := readHyperblock(nonce)
			if err != nil {
				return nil, err
			}

			return &response.Data.Hyperblock, nil
		},
		CalculateBlockTimestampUnixCalled: func(round uint64) int64 {
			return int64(networkConfig.StartTime)*1000 + int64(round*networkConfig.RoundDuration)
		},
		DecodeAddressCalled: func(address string) ([]byte, error) {
			var publicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32)
			return publicKeyConverter.Decode(address)
		},
	}

	return &blockAPIService{
		elrondProvider: providerMock,
		txsParser:      newTransactionParser(providerMock, configuration, networkConfig),
	}
}

func readHyperblock(nonce int64) (*data.HyperblockApiResponse, error) {
	filePath := fmt.Sprintf("testdata/localnet_%d_hyperblock.json", nonce)
	response := &data.HyperblockApiResponse{}

	err := readJson(filePath, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func readRosettaBlock(nonce int64) (*types.BlockResponse, error) {
	filePath := fmt.Sprintf("testdata/localnet_%d_rosetta.json", nonce)
	response := &types.BlockResponse{}

	err := readJson(filePath, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func readJson(filePath string, value interface{}) error {
	file, err := core.OpenFile(filePath)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, value)
	if err != nil {
		return err
	}

	return nil
}

func marshalPrettyJson(value interface{}) (string, error) {
	const emptyString = ""
	const fourSpaces = "    "
	content, err := json.MarshalIndent(value, emptyString, fourSpaces)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
