package process

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"strconv"
	"sync"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/kyber/singlesig"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func getSingleSigner() crypto.SingleSigner {
	return &singlesig.SchnorrSigner{}
}

// FaucetProcessor will handle the faucet operation
type FaucetProcessor struct {
	baseProc           Processor
	accMapByShard      map[uint32][]crypto.PrivateKey
	mutMap             sync.RWMutex
	singleSigner       crypto.SingleSigner
	minGasLimit        uint64
	minGasPrice        uint64
	defaultFaucetValue *big.Int
}

// NewFaucetProcessor will return a new instance of FaucetProcessor
func NewFaucetProcessor(
	ecConf *config.EconomicsConfig,
	baseProc Processor,
	privKeysLoader PrivateKeysLoaderHandler,
	defaultFaucetValue *big.Int,
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

	accMap, err := privKeysLoader.MapOfPrivateKeysByShard()
	if err != nil {
		return nil, err
	}

	if len(accMap) == 0 {
		return nil, ErrEmptyMapOfAccountsFromPem
	}

	minGasPrice, minGasLimit, err := parseEconomicsConfig(ecConf)
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
		minGasLimit:        minGasLimit,
		defaultFaucetValue: defaultFaucetValue,
	}, nil
}

// SenderDetailsFromPem will return details for a sender in the same shard with the receiver
func (fp *FaucetProcessor) SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error) {
	receiverBytes, err := hex.DecodeString(receiver)
	if err != nil {
		return nil, "", err
	}

	receiverShardId, err := fp.baseProc.ComputeShardId(receiverBytes)
	if err != nil {
		return nil, "", err
	}

	senderPrivKey := fp.getPrivKeyFromShard(receiverShardId)

	senderPubKeyPubKey := senderPrivKey.GeneratePublic()
	senderPubKeyBytes, err := senderPubKeyPubKey.ToByteArray()
	if err != nil {
		return nil, "", err
	}

	senderPubKeyHex := hex.EncodeToString(senderPubKeyBytes)

	return senderPrivKey, senderPubKeyHex, nil
}

// GenerateTxForSendUserFunds transmits a request to the right observer to load a provided address with some predefined balance
func (fp *FaucetProcessor) GenerateTxForSendUserFunds(
	senderSk crypto.PrivateKey,
	senderPk string,
	senderNonce uint64,
	receiver string,
	value *big.Int,
) (*data.Transaction, error) {

	if value == nil {
		value = fp.defaultFaucetValue
	}

	genTx := data.Transaction{
		Nonce:     senderNonce,
		Value:     value,
		Receiver:  receiver,
		Sender:    senderPk,
		GasPrice:  fp.minGasPrice,
		GasLimit:  fp.minGasLimit,
		Data:      "",
		Signature: "",
		Challenge: "",
	}

	return fp.getSignedTx(&genTx, senderSk)
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
	snrB, err := hex.DecodeString(tx.Sender)
	if err != nil {
		return nil, err
	}

	rcB, err := hex.DecodeString(tx.Receiver)
	if err != nil {
		return nil, err
	}

	erdTx := erdTransaction{
		Nonce:    tx.Nonce,
		Value:    tx.Value,
		RcvAddr:  rcB,
		SndAddr:  snrB,
		GasPrice: tx.GasPrice,
		GasLimit: tx.GasLimit,
		Data:     tx.Data,
	}

	return json.Marshal(erdTx)
}

func (fp *FaucetProcessor) getPrivKeyFromShard(shardId uint32) crypto.PrivateKey {
	fp.mutMap.Lock()
	defer fp.mutMap.Unlock()

	randomPrivKeyIdx := rand.Intn(len(fp.accMapByShard[shardId]))
	return fp.accMapByShard[shardId][randomPrivKeyIdx]
}

func parseEconomicsConfig(ecConf *config.EconomicsConfig) (uint64, uint64, error) {
	conversionBase := 10
	bitConversionSize := 64

	minGasPrice, err := strconv.ParseUint(ecConf.FeeSettings.MinGasPrice, conversionBase, bitConversionSize)
	if err != nil {
		return 0, 0, err
	}

	minGasLimit, err := strconv.ParseUint(ecConf.FeeSettings.MinGasLimit, conversionBase, bitConversionSize)
	if err != nil {
		return 0, 0, err
	}

	return minGasPrice, minGasLimit, nil
}
