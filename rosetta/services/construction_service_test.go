package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func TestConstructionAPIService_ConstructionPreprocess(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	elrondProvider := &mocks.ElrondProviderMock{}

	constructionAPIService := NewConstructionAPIService(elrondProvider, cfg, networkCfg, false)

	senderAddr := "senderAddr"
	receiverAddr := "receiverAddr"
	value := "123456"

	operations := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: senderAddr,
			},
			Amount: &types.Amount{
				Value:    "-" + value,
				Currency: cfg.Currency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{Index: 0},
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: receiverAddr,
			},
			Amount: &types.Amount{
				Value:    value,
				Currency: cfg.Currency,
			},
		},
	}

	feeMultiplier := 1.1
	maxFee := "1234567"

	gasPrice := uint64(100)
	gasLimit := uint64(10000)
	dataField := "data"

	response, err := constructionAPIService.ConstructionPreprocess(context.Background(),
		&types.ConstructionPreprocessRequest{
			Operations: operations,
			MaxFee: []*types.Amount{
				{
					Value:    maxFee,
					Currency: cfg.Currency,
				},
			},
			SuggestedFeeMultiplier: &feeMultiplier,
			Metadata: objectsMap{
				"gasPrice": gasPrice,
				"gasLimit": gasLimit,
				"data":     dataField,
			},
		},
	)
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}{
		"receiver":      receiverAddr,
		"sender":        senderAddr,
		"gasPrice":      gasPrice,
		"gasLimit":      gasLimit,
		"feeMultiplier": feeMultiplier,
		"data":          dataField,
		"value":         value,
		"maxFee":        maxFee,
		"type":          opTransfer,
	}, response.Options)
}

func TestConstructionAPIService_ConstructionMetadata(t *testing.T) {
	t.Parallel()

	senderAddr := "senderAddr"
	receiverAddr := "receiverAddr"
	value := "123456"
	feeMultiplier := 1.1
	maxFee := "1234567"
	gasPrice := uint64(100)
	gasLimit := uint64(10000)
	dataField := "data"
	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	nonce := uint64(5)
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	elrondProvider := &mocks.ElrondProviderMock{
		GetAccountCalled: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: senderAddr,
				Nonce:   nonce,
			}, nil
		},
	}

	constructionAPIService := NewConstructionAPIService(elrondProvider, cfg, networkCfg, false)

	options := map[string]interface{}{
		"receiver":      receiverAddr,
		"sender":        senderAddr,
		"gasPrice":      gasPrice,
		"gasLimit":      gasLimit,
		"feeMultiplier": feeMultiplier,
		"data":          dataField,
		"value":         value,
		"maxFee":        maxFee,
		"type":          opTransfer,
	}
	response, err := constructionAPIService.ConstructionMetadata(context.Background(),
		&types.ConstructionMetadataRequest{
			Options: options,
		},
	)

	expectedSuggestedFee := "1100000"
	require.Nil(t, err)
	require.Equal(t, expectedSuggestedFee, response.SuggestedFee[0].Value)
	require.Equal(t, cfg.Currency, response.SuggestedFee[0].Currency)

	delete(options, "feeMultiplier")
	delete(options, "maxFee")
	delete(options, "type")
	options["chainID"] = networkCfg.ChainID
	options["version"] = networkCfg.MinTxVersion
	options["data"] = []byte(fmt.Sprintf("%v", dataField))
	options["nonce"] = nonce
	options["gasLimit"] = uint64(10000)
	options["gasPrice"] = uint64(110)
	require.Equal(t, options, response.Metadata)
}

func TestConstructionAPIService_ConstructionPayloads(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})

	nonce := uint64(5)
	senderAddr := "senderAddr"
	receiverAddr := "receiverAddr"
	value := "123456"
	gasPrice := uint64(100)
	gasLimit := uint64(10000)
	dataField := "data"
	metadata := map[string]interface{}{
		"nonce":    nonce,
		"receiver": receiverAddr,
		"sender":   senderAddr,
		"gasPrice": gasPrice,
		"gasLimit": gasLimit,
		"data":     []byte(fmt.Sprintf("%v", dataField)),
		"value":    value,
		"chainID":  networkCfg.ChainID,
		"version":  networkCfg.MinTxVersion,
	}

	operations := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: senderAddr,
			},
			Amount: &types.Amount{
				Value:    "-" + value,
				Currency: cfg.Currency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{Index: 0},
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: receiverAddr,
			},
			Amount: &types.Amount{
				Value:    value,
				Currency: cfg.Currency,
			},
		},
	}

	constructionAPIService := NewConstructionAPIService(&mocks.ElrondProviderMock{}, cfg, networkCfg, false)

	response, err := constructionAPIService.ConstructionPayloads(context.Background(),
		&types.ConstructionPayloadsRequest{
			Operations: operations,
			Metadata:   metadata,
		},
	)

	marshalizedTx := []byte("{\"nonce\":5,\"value\":\"123456\",\"receiver\":\"receiverAddr\",\"sender\":\"senderAddr\",\"gasPrice\":100,\"gasLimit\":10000,\"data\":\"ZGF0YQ==\",\"chainID\":\"local-testnet\",\"version\":1}")
	unsignedTx := "7b226e6f6e6365223a352c2276616c7565223a22313233343536222c227265636569766572223a22726563656976657241646472222c2273656e646572223a2273656e64657241646472222c226761735072696365223a3130302c226761734c696d6974223a31303030302c2264617461223a225a47463059513d3d222c22636861696e4944223a226c6f63616c2d746573746e6574222c2276657273696f6e223a317d"

	require.Nil(t, err)
	require.Equal(t, unsignedTx, response.UnsignedTransaction)
	require.Equal(t, marshalizedTx, response.Payloads[0].Bytes)
	require.Equal(t, senderAddr, response.Payloads[0].AccountIdentifier.Address)
	require.Equal(t, types.Ed25519, response.Payloads[0].SignatureType)
}

func TestConstructionAPIService_ConstructionParse(t *testing.T) {
	t.Parallel()

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	constructionAPIService := NewConstructionAPIService(&mocks.ElrondProviderMock{}, cfg, networkCfg, false)
	unsignedTx := "7b226e6f6e6365223a352c2276616c7565223a22313233343536222c227265636569766572223a22726563656976657241646472222c2273656e646572223a2273656e64657241646472222c226761735072696365223a3130302c226761734c696d6974223a31303030302c2264617461223a225a47463059513d3d222c22636861696e4944223a226c6f63616c2d746573746e6574222c2276657273696f6e223a317d"

	senderAddr := "senderAddr"
	receiverAddr := "receiverAddr"
	value := "123456"
	operations := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: senderAddr,
			},
			Amount: &types.Amount{
				Value:    "-" + value,
				Currency: cfg.Currency,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{Index: 0},
			},
			Type: opTransfer,
			Account: &types.AccountIdentifier{
				Address: receiverAddr,
			},
			Amount: &types.Amount{
				Value:    value,
				Currency: cfg.Currency,
			},
		},
	}

	response, err := constructionAPIService.ConstructionParse(context.Background(),
		&types.ConstructionParseRequest{
			Signed:      false,
			Transaction: unsignedTx,
		},
	)
	require.Nil(t, err)
	require.Equal(t, operations, response.Operations)
	require.Nil(t, response.AccountIdentifierSigners)

	response, err = constructionAPIService.ConstructionParse(context.Background(),
		&types.ConstructionParseRequest{
			Signed:      true,
			Transaction: unsignedTx,
		},
	)
	require.Nil(t, err)
	require.Equal(t, operations, response.Operations)
	require.NotNil(t, response.AccountIdentifierSigners)
}

func TestConstructionAPIService_ConstructionCombine(t *testing.T) {
	t.Parallel()

	unsignedTx := "7b226e6f6e6365223a352c2276616c7565223a22313233343536222c227265636569766572223a22726563656976657241646472222c2273656e646572223a2273656e64657241646472222c226761735072696365223a3130302c226761734c696d6974223a31303030302c2264617461223a225a47463059513d3d222c22636861696e4944223a226c6f63616c2d746573746e6574222c2276657273696f6e223a317d"
	signature := []byte("signature-signature-signature")

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	constructionAPIService := NewConstructionAPIService(&mocks.ElrondProviderMock{}, cfg, networkCfg, false)

	response, err := constructionAPIService.ConstructionCombine(context.Background(),
		&types.ConstructionCombineRequest{
			UnsignedTransaction: unsignedTx,
			Signatures: []*types.Signature{
				{
					Bytes: signature,
				},
			},
		},
	)

	signedTx := "7b226e6f6e6365223a352c2276616c7565223a22313233343536222c227265636569766572223a22726563656976657241646472222c2273656e646572223a2273656e64657241646472222c226761735072696365223a3130302c226761734c696d6974223a31303030302c2264617461223a225a47463059513d3d222c227369676e6174757265223a2237333639363736653631373437353732363532643733363936373665363137343735373236353264373336393637366536313734373537323635222c22636861696e4944223a226c6f63616c2d746573746e6574222c2276657273696f6e223a317d"
	require.Nil(t, err)
	require.Equal(t, signedTx, response.SignedTransaction)
}

func TestConstructionAPIService_ConstructionDerive(t *testing.T) {
	t.Parallel()

	encodedAddress := "erd12312321321321123321321"
	elrondProvider := &mocks.ElrondProviderMock{
		EncodeAddressCalled: func(address []byte) (string, error) {
			return encodedAddress, nil
		},
	}

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	constructionAPIService := NewConstructionAPIService(elrondProvider, cfg, networkCfg, false)

	response, err := constructionAPIService.ConstructionDerive(context.Background(),
		&types.ConstructionDeriveRequest{
			PublicKey: &types.PublicKey{
				Bytes:     []byte("blablabla"),
				CurveType: types.Edwards25519,
			},
		},
	)
	require.Nil(t, err)
	require.Equal(t, encodedAddress, response.AccountIdentifier.Address)
}

func TestConstructionAPIService_ConstructionHash(t *testing.T) {
	t.Parallel()

	txHash := "hash-hash-hash"
	signedTx := "7b226e6f6e6365223a352c2276616c7565223a22313233343536222c227265636569766572223a22726563656976657241646472222c2273656e646572223a2273656e64657241646472222c226761735072696365223a3130302c226761734c696d6974223a31303030302c2264617461223a225a47463059513d3d222c227369676e6174757265223a2237333639363736653631373437353732363532643733363936373665363137343735373236353264373336393637366536313734373537323635222c22636861696e4944223a226c6f63616c2d746573746e6574222c2276657273696f6e223a317d"
	elrondProvider := &mocks.ElrondProviderMock{
		ComputeTransactionHashCalled: func(tx *data.Transaction) (string, error) {
			return txHash, nil
		},
	}

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	constructionAPIService := NewConstructionAPIService(elrondProvider, cfg, networkCfg, false)

	response, err := constructionAPIService.ConstructionHash(context.Background(),
		&types.ConstructionHashRequest{
			SignedTransaction: signedTx,
		},
	)
	require.Nil(t, err)
	require.Equal(t, txHash, response.TransactionIdentifier.Hash)
}

func TestConstructionAPIService_ConstructionSubmit(t *testing.T) {
	t.Parallel()

	txHash := "hash-hash-hash"
	signedTx := "7b226e6f6e6365223a352c2276616c7565223a22313233343536222c227265636569766572223a22726563656976657241646472222c2273656e646572223a2273656e64657241646472222c226761735072696365223a3130302c226761734c696d6974223a31303030302c2264617461223a225a47463059513d3d222c227369676e6174757265223a2237333639363736653631373437353732363532643733363936373665363137343735373236353264373336393637366536313734373537323635222c22636861696e4944223a226c6f63616c2d746573746e6574222c2276657273696f6e223a317d"
	elrondProvider := &mocks.ElrondProviderMock{
		SendTxCalled: func(tx *data.Transaction) (string, error) {
			return txHash, nil
		},
	}

	networkCfg := &provider.NetworkConfig{
		GasPerDataByte: 1,
		ClientVersion:  "",
		MinGasPrice:    10,
		MinGasLimit:    100,
		ChainID:        "local-testnet",
		MinTxVersion:   1,
	}
	cfg := configuration.LoadConfiguration(networkCfg, &config.Config{})
	constructionAPIService := NewConstructionAPIService(elrondProvider, cfg, networkCfg, false)

	response, err := constructionAPIService.ConstructionSubmit(context.Background(),
		&types.ConstructionSubmitRequest{
			SignedTransaction: signedTx,
		},
	)
	require.Nil(t, err)
	require.Equal(t, txHash, response.TransactionIdentifier.Hash)
}
