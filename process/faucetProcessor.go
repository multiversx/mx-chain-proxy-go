package process

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"strconv"
	"sync"

	erdConfig "github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/crypto"
	ed25519SingleSigner "github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519/singlesig"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/economics"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/disabled"
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
	minGasPrice        uint64
	defaultFaucetValue *big.Int
	econData           process.FeeHandler
	pubKeyConverter    core.PubkeyConverter
}

// NewFaucetProcessor will return a new instance of FaucetProcessor
func NewFaucetProcessor(
	ecConf *erdConfig.EconomicsConfig,
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

	econData, minGasPrice, err := parseEconomicsConfig(ecConf)
	if err != nil {
		return nil, ErrInvalidEconomicsConfig
	}

	singleSigner := getSingleSigner()
	return &FaucetProcessor{
		baseProc:           baseProc,
		accMapByShard:      accMap,
		mutMap:             sync.RWMutex{},
		singleSigner:       singleSigner,
		minGasPrice:        minGasPrice,
		defaultFaucetValue: defaultFaucetValue,
		econData:           econData,
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

	senderPubKeyString := fp.pubKeyConverter.Encode(senderPubKeyBytes)

	return senderPrivKey, senderPubKeyString, nil
}

// GenerateTxForSendUserFunds transmits a request to the right observer to load a provided address with some predefined balance
func (fp *FaucetProcessor) GenerateTxForSendUserFunds(
	senderSk crypto.PrivateKey,
	senderPk string,
	senderNonce uint64,
	receiver string,
	value *big.Int,
	chainID string,
	version uint32,
) (*data.Transaction, error) {

	if value == nil {
		value = fp.defaultFaucetValue
	}

	genTx := data.Transaction{
		Nonce:     senderNonce,
		Value:     value.String(),
		Receiver:  receiver,
		Sender:    senderPk,
		GasPrice:  fp.minGasPrice,
		Data:      []byte(""),
		Signature: "",
		ChainID:   chainID,
		Version:   version,
	}

	wrappedTx, err := data.NewTransactionWrapper(&genTx, fp.pubKeyConverter)
	if err != nil {
		return nil, err
	}

	gasLimit := fp.econData.ComputeGasLimit(wrappedTx)
	genTx.GasLimit = gasLimit

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

func parseEconomicsConfig(ecConf *erdConfig.EconomicsConfig) (process.FeeHandler, uint64, error) {
	argsNewEconomics := economics.ArgsNewEconomicsData{
		Economics:                      ecConf,
		PenalizedTooMuchGasEnableEpoch: 0,
		EpochNotifier:                  &disabled.EpochStartNotifier{},
	}
	econData, err := economics.NewEconomicsData(argsNewEconomics)
	if err != nil {
		return nil, 0, err
	}
	conversionBase := 10
	bitConversionSize := 64

	minGasPrice, err := strconv.ParseUint(ecConf.FeeSettings.MinGasPrice, conversionBase, bitConversionSize)
	if err != nil {
		return nil, 0, err
	}

	return econData, minGasPrice, nil
}
