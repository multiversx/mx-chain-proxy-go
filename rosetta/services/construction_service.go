package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type constructionAPIService struct {
	elrondClient client.ElrondClientHandler
	config       *configuration.Configuration
	txsParser    *transactionsParser
}

// NewConstructionAPIService creates a new instance of an constructionAPIService.
func NewConstructionAPIService(elrondClient client.ElrondClientHandler, cfg *configuration.Configuration) server.ConstructionAPIServicer {
	return &constructionAPIService{
		elrondClient: elrondClient,
		config:       cfg,
		txsParser:    newTransactionParser(cfg),
	}
}

func checkOperationsAndMeta(ops []*types.Operation, meta map[string]interface{}) *types.Error {
	terr := ErrConstructionCheck
	if len(ops) == 0 {
		terr.Message += "invalid number of operations"
		return terr
	}

	for _, op := range ops {
		if !checkOperationsType(op) {
			terr.Message += "unsupported operation type"
			return terr
		}
	}

	if meta["gasLimit"] != nil {
		if _, ok := meta["gasLimit"].(uint64); ok {
			terr.Message += "invalid gas limit"
			return terr
		}
	}
	if meta["gasPrice"] != nil {
		if _, ok := meta["gasPrice"].(uint64); ok {
			terr.Message += "invalid gas price"
			return terr
		}
	}

	return nil
}

func checkOperationsType(op *types.Operation) bool {
	for _, supOp := range SupportedOperationTypes {
		if supOp == op.Type {
			return true
		}
	}

	return false
}

func getOptionsFromOperations(ops []*types.Operation) objectsMap {
	options := make(objectsMap)
	options["sender"] = ops[0].Account.Address
	options["receiver"] = ops[1].Account.Address
	options["type"] = ops[0].Type
	options["value"] = ops[1].Amount.Value

	return options
}

//ConstructionPreprocess will preprocess data that in provided in request
func (cas *constructionAPIService) ConstructionPreprocess(
	_ context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	if err := checkOperationsAndMeta(request.Operations, request.Metadata); err != nil {
		return nil, err
	}

	options := getOptionsFromOperations(request.Operations)

	if len(request.MaxFee) > 0 {
		maxFee := request.MaxFee[0]
		if maxFee.Currency.Symbol != cas.config.Currency.Symbol ||
			maxFee.Currency.Decimals != cas.config.Currency.Decimals {
			terr := ErrConstructionCheck
			terr.Message += "invalid currency"
			return nil, terr
		}

		options["maxFee"] = maxFee.Value
	}

	if request.SuggestedFeeMultiplier != nil {
		options["feeMultiplier"] = *request.SuggestedFeeMultiplier
	}

	if request.Metadata["gasLimit"] != nil {
		options["gasLimit"] = request.Metadata["gasLimit"]
	}
	if request.Metadata["gasPrice"] != nil {
		options["gasPrice"] = request.Metadata["gasPrice"]
	}
	if request.Metadata["data"] != nil {
		options["data"] = request.Metadata["data"]
	}

	return &types.ConstructionPreprocessResponse{
		Options: options,
	}, nil
}

// ConstructionMetadata construct metadata for a transaction
func (cas *constructionAPIService) ConstructionMetadata(
	_ context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	txType, ok := request.Options["type"].(string)
	if !ok {
		terr := ErrInvalidInputParam
		terr.Message += "transaction type"
		return nil, terr
	}

	networkConfig, err := cas.elrondClient.GetNetworkConfig()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetNetworkConfig, err)
	}

	metadata, errS := cas.computeMetadata(request.Options, networkConfig)
	if errS != nil {
		return nil, errS
	}

	suggestedFee, gasPrice, gasLimit, errS := computeSuggestedFeeAndGas(txType, request.Options, networkConfig)
	if errS != nil {
		return nil, errS
	}

	metadata["gasLimit"] = gasLimit
	metadata["gasPrice"] = gasPrice

	return &types.ConstructionMetadataResponse{
		Metadata: metadata,
		SuggestedFee: []*types.Amount{
			{
				Value:    suggestedFee.String(),
				Currency: cas.config.Currency,
			},
		},
	}, nil
}

func (cas *constructionAPIService) computeMetadata(options objectsMap, networkConfig *client.NetworkConfig) (objectsMap, *types.Error) {
	metadata := make(objectsMap)
	if dataField, ok := options["data"]; ok {
		// convert string to byte array
		metadata["data"] = []byte(fmt.Sprintf("%v", dataField))
	}

	var ok bool
	if metadata["sender"], ok = options["sender"]; !ok {
		return nil, wrapErr(ErrMalformedValue, errors.New("sender address missing"))
	}
	if metadata["receiver"], ok = options["receiver"]; !ok {
		return nil, wrapErr(ErrMalformedValue, errors.New("receiver address missing"))
	}
	if metadata["value"], ok = options["value"]; !ok {
		return nil, wrapErr(ErrMalformedValue, errors.New("value missing"))
	}

	metadata["chainID"] = networkConfig.ChainID
	metadata["version"] = networkConfig.MinTxVersion

	account, err := cas.elrondClient.GetAccount(options["sender"].(string))
	if err != nil {
		return nil, wrapErr(ErrUnableToGetAccount, err)
	}

	metadata["nonce"] = account.Nonce

	return metadata, nil
}

// ConstructionPayloads will prepare a transaction for signing
func (cas *constructionAPIService) ConstructionPayloads(
	_ context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	erdTx, err := createTransaction(request)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	mtx, err := json.Marshal(erdTx)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	unsignedTx := hex.EncodeToString(mtx)

	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedTx,
		Payloads: []*types.SigningPayload{
			{
				AccountIdentifier: &types.AccountIdentifier{
					Address: request.Operations[0].Account.Address,
				},
				SignatureType: types.Ed25519,
				Bytes:         mtx,
			},
		},
	}, nil
}

// ConstructionParse will check if a transaction is correctly formatted
func (cas *constructionAPIService) ConstructionParse(
	_ context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {
	elrondTx, err := getTxFromRequest(request.Transaction)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	var signers []*types.AccountIdentifier
	if request.Signed {
		signers = []*types.AccountIdentifier{
			{
				Address: elrondTx.Sender,
			},
		}
	}

	return &types.ConstructionParseResponse{
		Operations:               cas.txsParser.createOperationsFromPreparedTx(elrondTx),
		AccountIdentifierSigners: signers,
	}, nil
}

func createTransaction(request *types.ConstructionPayloadsRequest) (*data.Transaction, error) {
	tx := &data.Transaction{}

	requestMetadataBytes, err := json.Marshal(request.Metadata)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestMetadataBytes, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func getTxFromRequest(txString string) (*data.Transaction, error) {
	txBytes, err := hex.DecodeString(txString)
	if err != nil {
		return nil, err
	}

	var elrondTx data.Transaction
	err = json.Unmarshal(txBytes, &elrondTx)
	if err != nil {
		return nil, err
	}

	return &elrondTx, nil
}

//ConstructionCombine will create a signed transaction for transaction bytes and signature
func (cas *constructionAPIService) ConstructionCombine(
	_ context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	elrondTx, err := getTxFromRequest(request.UnsignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	if len(request.Signatures) != 1 {
		return nil, ErrInvalidInputParam
	}

	txSignature := hex.EncodeToString(request.Signatures[0].Bytes)
	elrondTx.Signature = txSignature

	signedTxBytes, err := json.Marshal(elrondTx)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	signedTx := hex.EncodeToString(signedTxBytes)

	return &types.ConstructionCombineResponse{
		SignedTransaction: signedTx,
	}, nil
}

// ConstructionDerive return a bech32 address from public key bytes
func (cas *constructionAPIService) ConstructionDerive(
	_ context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	if request.PublicKey.CurveType != types.Edwards25519 {
		return nil, ErrUnsupportedCurveType
	}

	pubKey := request.PublicKey.Bytes
	address, err := cas.elrondClient.EncodeAddress(pubKey)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	return &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: address,
		},
		Metadata: nil,
	}, nil
}

// ConstructionHash will calculate transaction hash
func (cas *constructionAPIService) ConstructionHash(
	_ context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	elrondTx, err := getTxFromRequest(request.SignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	txHash, err := cas.elrondClient.ComputeTransactionHash(elrondTx)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txHash,
		},
	}, nil
}

// ConstructionSubmit will submit transaction and return hash
func (cas *constructionAPIService) ConstructionSubmit(
	_ context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	elrondTx, err := getTxFromRequest(request.SignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrMalformedValue, err)
	}

	txHash, err := cas.elrondClient.SendTx(elrondTx)
	if err != nil {
		return nil, wrapErr(ErrUnableToSubmitTransaction, err)
	}

	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: txHash,
		},
	}, nil
}
