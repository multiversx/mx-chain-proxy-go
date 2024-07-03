package transactionDecoder

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
)

// GetTransactionMetadata will decode the Data field from the transaction and populate it in the TransactionMetadata.
func GetTransactionMetadata(tx TransactionToDecode, pubKeyConverter core.PubkeyConverter) (*TransactionMetadata, error) {
	metadata, err := getNormalTransactionMetadata(tx, pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction metadata: %w", err)
	}

	esdtMetadata, err := getESDTTransactionMetadata(*metadata)
	if err == nil {
		return esdtMetadata, nil
	}

	nftMetadata, err := getNFTTransferMetadata(*metadata, pubKeyConverter)
	if err == nil {
		return nftMetadata, nil
	}

	multiMetadata, err := getMultiTransferMetadata(*metadata, pubKeyConverter)
	if err == nil {
		return multiMetadata, nil
	}

	return metadata, nil
}

func getNormalTransactionMetadata(tx TransactionToDecode, pubKeyConverter core.PubkeyConverter) (*TransactionMetadata, error) {
	v := "0"
	if tx.Value != "" {
		v = tx.Value
	}
	value, ok := big.NewInt(0).SetString(v, 10)
	if !ok {
		return nil, ErrValueSet
	}

	metadata := TransactionMetadata{
		Sender:   tx.Sender,
		Receiver: tx.Receiver,
		Value:    value,
	}

	if tx.Data != "" {
		decodedData, err := base64.StdEncoding.DecodeString(tx.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transaction data: %w", err)
		}

		dataComponents := strings.Split(string(decodedData), "@")

		everyCheck := true
		args := dataComponents[1:]
		for _, arg := range args {
			if !isSmartContractArgument(arg) {
				everyCheck = false
			}
		}

		if everyCheck {
			metadata.FunctionName = dataComponents[0]
			metadata.FunctionArgs = args
		}

		if metadata.FunctionName == "relayedTx" && metadata.FunctionArgs != nil && len(metadata.FunctionArgs) == 1 {
			relayedTx, relayErr := parseRelayedV1(metadata, pubKeyConverter)
			if relayErr != nil {
				return nil, fmt.Errorf("failed to parse relayed v1 transaction metadata: %w", relayErr)
			}

			return getNormalTransactionMetadata(*relayedTx, pubKeyConverter)
		}

		if metadata.FunctionName == "relayedTxV2" &&
			metadata.FunctionArgs != nil &&
			len(metadata.FunctionArgs) == 4 {

			relayedTxV2, relayErr := parseRelayedV2(metadata, pubKeyConverter)
			if relayErr != nil {
				return nil, fmt.Errorf("failed to parse relayed v2 transaction metadata: %w", relayErr)
			}

			return getNormalTransactionMetadata(*relayedTxV2, pubKeyConverter)
		}

	}

	return &metadata, nil

}

func parseRelayedV1(metadata TransactionMetadata, pubKeyConverter core.PubkeyConverter) (*TransactionToDecode, error) {
	data, err := hex.DecodeString(metadata.FunctionArgs[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode data field: %w", err)
	}

	var innerTx = struct {
		Nonce     int    `json:"nonce"`
		Sender    string `json:"sender"`
		Receiver  string `json:"receiver"`
		Value     int    `json:"value"`
		GasPrice  int    `json:"gasPrice"`
		GasLimit  int    `json:"gasLimit"`
		Data      string `json:"data"`
		Signature string `json:"signature"`
		ChainID   string `json:"chainID"`
		Version   int    `json:"version"`
	}{}
	err = json.Unmarshal(data, &innerTx)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal inner transaction: %w", err)
	}
	sender, err := retrieveAddressBech32FromBase64(innerTx.Sender, pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sender field: %w", err)
	}
	receiver, err := retrieveAddressBech32FromBase64(innerTx.Receiver, pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode receiver field: %w", err)
	}

	relayedTx := &TransactionToDecode{
		Sender:   sender,
		Receiver: receiver,
		Data:     innerTx.Data,
		Value:    strconv.FormatInt(int64(innerTx.Value), 10),
	}
	return relayedTx, nil
}

func parseRelayedV2(metadata TransactionMetadata, pubKeyConverter core.PubkeyConverter) (*TransactionToDecode, error) {
	relayedTx := &TransactionToDecode{}

	relayedTx.Sender = metadata.Receiver
	receiver, err := retrieveAddressBech32FromHex(metadata.FunctionArgs[0], pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode receiver field: %w", err)
	}
	relayedTx.Receiver = receiver

	decodeString, err := hex.DecodeString(metadata.FunctionArgs[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode data field: %w", err)
	}
	relayedTx.Data = base64.StdEncoding.EncodeToString(decodeString)
	relayedTx.Value = "0"
	return relayedTx, nil
}

func getESDTTransactionMetadata(metadata TransactionMetadata) (*TransactionMetadata, error) {
	if metadata.FunctionName != "ESDTTransfer" {
		return nil, ErrNotESDTTransfer
	}

	args := metadata.FunctionArgs
	if args == nil {
		return nil, ErrNoArgs
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments. required [%d], found [%d]", 2, len(args))
	}

	tokenIdentifier, err := hex.DecodeString(args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token identifier: %w", err)
	}

	value, ok := big.NewInt(0).SetString(args[1], 16)
	if !ok {
		return nil, fmt.Errorf("failed to set value: %w", err)
	}

	result := &TransactionMetadata{}
	result.Sender = metadata.Sender
	result.Receiver = metadata.Receiver
	result.Value = value

	if len(args) > 2 {
		functionName, err := hex.DecodeString(args[2])
		if err != nil {
			return nil, fmt.Errorf("failed to decode function name: %w", err)
		}
		result.FunctionName = string(functionName)
		result.FunctionArgs = args[3:]
	}

	result.Transfers = []TransactionMetadataTransfer{
		{
			Value: value,
			Properties: TokenTransferProperties{
				Identifier: string(tokenIdentifier),
			},
		},
	}

	return result, nil
}

func getNFTTransferMetadata(metadata TransactionMetadata, pubKeyConverter core.PubkeyConverter) (*TransactionMetadata, error) {
	if metadata.Sender != metadata.Receiver {
		return nil, ErrSenderReceiver
	}

	if metadata.FunctionName != "ESDTNFTTransfer" {
		return nil, ErrNotESDTNFTTransfer
	}

	args := metadata.FunctionArgs
	if args == nil {
		return nil, ErrNoArgs
	}

	if len(args) < 4 {
		return nil, fmt.Errorf("not enough arguments. required [%d], found [%d]", 2, len(args))
	}

	if !isAddressValid(args[3]) {
		return nil, ErrInvalidAddress
	}

	collectionIdentifier, err := hex.DecodeString(args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode collection identifier: %w", err)
	}

	nonce := args[1]
	value, ok := big.NewInt(0).SetString(args[2], 16)
	if !ok {
		return nil, ErrValueSet
	}

	receiver, err := retrieveAddressBech32FromHex(args[3], pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode receiver field: %w", err)
	}

	result := &TransactionMetadata{
		Sender:   metadata.Sender,
		Receiver: receiver,
		Value:    value,
	}

	if len(args) > 4 {
		functionName, err := hex.DecodeString(args[4])
		if err != nil {
			return nil, fmt.Errorf("failed to decode function name: %w", err)
		}
		result.FunctionName = string(functionName)
		result.FunctionArgs = args[5:]
	}

	result.Transfers = []TransactionMetadataTransfer{
		{
			Value: value,
			Properties: TokenTransferProperties{
				Collection: string(collectionIdentifier),
				Identifier: fmt.Sprintf("%s-%s", collectionIdentifier, nonce),
			},
		},
	}

	return result, nil
}

func getMultiTransferMetadata(metadata TransactionMetadata, pubKeyConverter core.PubkeyConverter) (*TransactionMetadata, error) {
	if metadata.Sender != metadata.Receiver {
		return nil, ErrSenderReceiver
	}

	if metadata.FunctionName != "MultiESDTNFTTransfer" {
		return nil, ErrNotMultiESDTNFTTransfer
	}

	if metadata.FunctionArgs == nil {
		return nil, ErrNoArgs
	}

	if len(metadata.FunctionArgs) < 3 {
		return nil, fmt.Errorf("not enough arguments. required [%d], found [%d]", 3, len(metadata.FunctionArgs))
	}

	args := metadata.FunctionArgs
	if !isAddressValid(args[0]) {
		return nil, ErrInvalidAddress
	}

	receiver, err := retrieveAddressBech32FromHex(args[0], pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to decode receiver field: %w", err)
	}

	transferCount, err := strconv.ParseInt(args[1], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transfer count: %w", err)
	}

	result := &TransactionMetadata{}
	transfers := make([]TransactionMetadataTransfer, 0)
	index := 2
	for i := int64(0); i < transferCount; i++ {
		identifier, err := hex.DecodeString(args[index])
		if err != nil {
			return nil, fmt.Errorf("failed to decode transfer identifier: %w", err)
		}
		index++

		n, nonceErr := strconv.ParseInt(args[index], 16, 64)
		nonce := args[index]
		index++

		value, ok := big.NewInt(0).SetString(args[index], 16)
		if !ok {
			return nil, fmt.Errorf("failed to set value: %w", err)
		}
		index++
		if nonceErr == nil && n > 0 {
			transfers = append(transfers, TransactionMetadataTransfer{
				Value: value,
				Properties: TokenTransferProperties{
					Collection: string(identifier),
					Identifier: fmt.Sprintf("%s-%s", identifier, nonce),
				},
			})
		} else {
			transfers = append(transfers, TransactionMetadataTransfer{
				Value: value,
				Properties: TokenTransferProperties{
					Token: string(identifier),
				},
			})
		}
	}

	result.Sender = metadata.Sender
	result.Receiver = receiver
	result.Transfers = transfers
	result.Value = big.NewInt(0)

	if len(args) > index {
		functionName, err := hex.DecodeString(args[index])
		if err != nil {
			return nil, fmt.Errorf("failed to decode function name: %w", err)
		}
		result.FunctionName = string(functionName)
		index++
		result.FunctionArgs = args[index:]
	}

	return result, nil
}

func isSmartContractArgument(arg string) bool {
	if _, err := hex.DecodeString(arg); err != nil {
		return false
	}

	if len(arg)%2 != 0 {
		return false
	}

	return true
}

func isAddressValid(address string) bool {
	decodeString, err := hex.DecodeString(address)
	if err != nil {
		return false
	}

	if len(decodeString) != 32 {
		return false
	}

	return true
}

func retrieveAddressBech32FromBase64(base64Encoded string, pubKeyConverter core.PubkeyConverter) (string, error) {
	hexEncoded, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return "", err
	}

	encode, err := pubKeyConverter.Encode(hexEncoded)
	if err != nil {
		return "", err
	}

	return encode, nil
}

func retrieveAddressBech32FromHex(hexEncoded string, pubKeyConverter core.PubkeyConverter) (string, error) {
	bytes, err := hex.DecodeString(hexEncoded)
	if err != nil {
		return "", err
	}
	address, err := pubKeyConverter.Encode(bytes)
	if err != nil {
		return "", err
	}

	return address, nil
}
