package process

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	ed25519SingleSigner "github.com/multiversx/mx-chain-crypto-go/signing/ed25519/singlesig"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

func getSingleSigner() crypto.SingleSigner {
	return &ed25519SingleSigner.Ed25519Signer{}
}

// FaucetProcessor will handle the faucet operation
type FaucetProcessor struct {
	baseProc           Processor
	accMapByShard      map[uint32][]crypto.PrivateKey
	mutMap             sync.RWMutex
	singleSigner       crypto.SingleSigner
	defaultFaucetValue *big.Int
	pubKeyConverter    core.PubkeyConverter
}

// NewFaucetProcessor will return a new instance of FaucetProcessor
func NewFaucetProcessor(
	baseProc Processor,
	privKeysLoader PrivateKeysLoaderHandler,
	defaultFaucetValue *big.Int,
	pubKeyConverter core.PubkeyConverter,
) (*FaucetProcessor, error) {
	if baseProc == nil {
		return nil, ErrNilCoreProcessor
	}
	if privKeysLoader == nil {
		return nil, ErrNilPrivateKeysLoader
	}
	if defaultFaucetValue == nil {
		return nil, ErrNilDefaultFaucetValue
	}
	if defaultFaucetValue.Cmp(big.NewInt(0)) <= 0 {
		return nil, ErrInvalidDefaultFaucetValue
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	accMap, err := privKeysLoader.PrivateKeysByShard()
	if err != nil {
		return nil, err
	}

	if len(accMap) == 0 {
		return nil, ErrEmptyMapOfAccountsFromPem
	}

	singleSigner := getSingleSigner()
	return &FaucetProcessor{
		baseProc:           baseProc,
		accMapByShard:      accMap,
		mutMap:             sync.RWMutex{},
		singleSigner:       singleSigner,
		defaultFaucetValue: defaultFaucetValue,
		pubKeyConverter:    pubKeyConverter,
	}, nil
}

// IsEnabled returns true
func (fp *FaucetProcessor) IsEnabled() bool {
	return true
}

// SenderDetailsFromPem will return details for a sender in the same shard with the receiver
func (fp *FaucetProcessor) SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error) {
	receiverBytes, err := fp.pubKeyConverter.Decode(receiver)
	if err != nil {
		return nil, "", err
	}

	receiverShardID, err := fp.baseProc.ComputeShardId(receiverBytes)
	if err != nil {
		return nil, "", err
	}

	senderPrivKey, err := fp.getPrivKeyFromShard(receiverShardID)
	if err != nil {
		return nil, "", err
	}

	senderPubKeyPubKey := senderPrivKey.GeneratePublic()
	senderPubKeyBytes, err := senderPubKeyPubKey.ToByteArray()
	if err != nil {
		return nil, "", err
	}

	senderPubKeyString := fp.pubKeyConverter.SilentEncode(senderPubKeyBytes, log)

	return senderPrivKey, senderPubKeyString, nil
}

// GenerateTxForSendUserFunds transmits a request to the right observer to load a provided address with some predefined balance
func (fp *FaucetProcessor) GenerateTxForSendUserFunds(
	senderSk crypto.PrivateKey,
	senderPk string,
	senderNonce uint64,
	receiver string,
	value *big.Int,
	networkConfig *data.NetworkConfig,
) (*data.Transaction, error) {
	if value == nil {
		value = fp.defaultFaucetValue
	}

	genTx := data.Transaction{
		Nonce:     senderNonce,
		Value:     value.String(),
		Receiver:  receiver,
		Sender:    senderPk,
		Data:      []byte(""),
		Signature: "",
		ChainID:   networkConfig.Config.ChainID,
		Version:   networkConfig.Config.MinTransactionVersion,
		GasPrice:  networkConfig.Config.MinGasPrice,
		GasLimit:  networkConfig.Config.MinGasLimit,
	}

	signedTx, err := fp.getSignedTx(&genTx, senderSk)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (fp *FaucetProcessor) getSignedTx(tx *data.Transaction, privKey crypto.PrivateKey) (*data.Transaction, error) {
	marshalizedTxBeforeSigning, err := fp.marshalTxForSigning(tx)
	if err != nil {
		return nil, err
	}

	signature, err := fp.singleSigner.Sign(privKey, marshalizedTxBeforeSigning)
	if err != nil {
		return nil, err
	}

	signHex := hex.EncodeToString(signature)
	tx.Signature = signHex

	return tx, nil
}

func (fp *FaucetProcessor) marshalTxForSigning(tx *data.Transaction) ([]byte, error) {
	erdTx := erdTransaction{
		Nonce:    tx.Nonce,
		Value:    tx.Value,
		RcvAddr:  tx.Receiver,
		SndAddr:  tx.Sender,
		GasPrice: tx.GasPrice,
		GasLimit: tx.GasLimit,
		Data:     tx.Data,
		ChainID:  tx.ChainID,
		Version:  tx.Version,
	}

	return json.Marshal(erdTx)
}

func (fp *FaucetProcessor) getPrivKeyFromShard(shardID uint32) (crypto.PrivateKey, error) {
	fp.mutMap.Lock()
	defer fp.mutMap.Unlock()

	accountsInShard, ok := fp.accMapByShard[shardID]
	if !ok || len(accountsInShard) == 0 {
		return nil, ErrNoFaucetAccountForGivenShard
	}

	randomPrivKeyIdx := rand.Intn(len(accountsInShard))
	return fp.accMapByShard[shardID][randomPrivKeyIdx], nil
}
