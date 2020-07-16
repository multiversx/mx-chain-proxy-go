package process_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"testing"

	erdConfig "github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewFaucetProcessor_NilBaseProcessorShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		nil,
		&mock.PrivateKeysLoaderStub{},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewFaucetProcessor_NilPrivateKeysLoaderShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		nil,
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrNilPrivateKeysLoader, err)
}

func TestNewFaucetProcessor_NilDefaultFaucetValueShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{},
		nil,
		&mock.PubKeyConverterMock{},
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrNilDefaultFaucetValue, err)
}

func TestNewFaucetProcessor_ZeroDefaultFaucetValueShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{},
		big.NewInt(0),
		&mock.PubKeyConverterMock{},
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrInvalidDefaultFaucetValue, err)
}

func TestNewFaucetProcessor_NegativeDefaultFaucetValueShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{},
		big.NewInt(-1),
		&mock.PubKeyConverterMock{},
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrInvalidDefaultFaucetValue, err)
}

func TestNewFaucetProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{},
		big.NewInt(10),
		nil,
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewFaucetProcessor_EmptyAccMapShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				return make(map[uint32][]crypto.PrivateKey), nil
			},
		},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrEmptyMapOfAccountsFromPem, err)
}

func TestNewFaucetProcessor_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], nil)

				return mapToReturn, nil
			},
		},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	assert.NotNil(t, fp)
	assert.Nil(t, err)
}

func TestFaucetProcessor_SenderDetailsFromPemWrongReceiverHexShouldErr(t *testing.T) {
	t.Parallel()

	receiver := "wrong receiver public key hex"
	fp, _ := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return uint32(0), nil
			},
		},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], nil)

				return mapToReturn, nil
			},
		},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	sk, pkHex, err := fp.SenderDetailsFromPem(receiver)
	assert.Nil(t, sk)
	assert.Equal(t, "", pkHex)
	assert.NotNil(t, err)
}

func TestFaucetProcessor_SenderDetailsFromPemShardIdComputationWrongShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("error computing shard id")
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	fp, _ := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return uint32(0), expectedErr
			},
		},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], nil)

				return mapToReturn, nil
			},
		},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	sk, pkHex, err := fp.SenderDetailsFromPem(receiver)
	assert.Nil(t, sk)
	assert.Equal(t, "", pkHex)
	assert.Equal(t, expectedErr, err)
}

func TestFaucetProcessor_SenderDetailsFromPemComputedShardIdNotFoundInAccountsShouldErr(t *testing.T) {
	t.Parallel()

	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	fp, _ := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return uint32(37), nil
			},
		},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], nil)

				return mapToReturn, nil
			},
		},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	sk, pkHex, err := fp.SenderDetailsFromPem(receiver)
	assert.Nil(t, sk)
	assert.Equal(t, "", pkHex)
	assert.Equal(t, process.ErrNoFaucetAccountForGivenShard, err)
}

func TestFaucetProcessor_SenderDetailsFromPemShouldWork(t *testing.T) {
	t.Parallel()

	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	expectedPrivKey := getPrivKey()
	fp, _ := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return uint32(0), nil
			},
		},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], expectedPrivKey)

				return mapToReturn, nil
			},
		},
		big.NewInt(1),
		&mock.PubKeyConverterMock{},
	)

	sk, pkHex, err := fp.SenderDetailsFromPem(receiver)
	assert.Equal(t, expectedPrivKey, sk)
	assert.NotEqual(t, "", pkHex)
	assert.Nil(t, err)
}

func TestFaucetProcessor_GenerateTxForSendUserFundsNilFaucetValueShouldUseDefault(t *testing.T) {
	t.Parallel()

	senderSk := getPrivKey()
	senderHexPk := hexPubKeyFromSk(senderSk)
	senderNonce := uint64(25)
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	defaultFaucetValue := big.NewInt(100000000000)

	fp, _ := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return uint32(0), nil
			},
		},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], getPrivKey())

				return mapToReturn, nil
			},
		},
		defaultFaucetValue,
		&mock.PubKeyConverterMock{},
	)

	tx, err := fp.GenerateTxForSendUserFunds(senderSk, senderHexPk, senderNonce, receiver, nil, "", 1)
	assert.Nil(t, err)
	assert.Equal(t, senderHexPk, tx.Sender)
	assert.Equal(t, receiver, tx.Receiver)
	assert.Equal(t, defaultFaucetValue.String(), tx.Value)
}

func TestFaucetProcessor_GenerateTxForSendUserFundsShouldWork(t *testing.T) {
	t.Parallel()

	senderSk := getPrivKey()
	senderHexPk := hexPubKeyFromSk(senderSk)
	senderNonce := uint64(25)
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	defaultFaucetValue := big.NewInt(100000000000)
	faucetValue := big.NewInt(12345)

	fp, _ := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return uint32(0), nil
			},
		},
		&mock.PrivateKeysLoaderStub{
			PrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], getPrivKey())

				return mapToReturn, nil
			},
		},
		defaultFaucetValue,
		&mock.PubKeyConverterMock{},
	)

	tx, err := fp.GenerateTxForSendUserFunds(senderSk, senderHexPk, senderNonce, receiver, faucetValue, "", 1)
	assert.Nil(t, err)
	assert.Equal(t, senderHexPk, tx.Sender)
	assert.Equal(t, receiver, tx.Receiver)
	assert.Equal(t, faucetValue.String(), tx.Value)
}

func getPrivKey() crypto.PrivateKey {
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	sk, _ := keyGen.GeneratePair()

	return sk
}

func hexPubKeyFromSk(sk crypto.PrivateKey) string {
	senderPk := sk.GeneratePublic()
	senderPkBytes, _ := senderPk.ToByteArray()
	senderPkHex := hex.EncodeToString(senderPkBytes)

	return senderPkHex
}

func testEconomicsConfig() *erdConfig.EconomicsConfig {
	maxGasLimitPerBlock := strconv.FormatUint(uint64(1500000000), 10)
	maxGasLimitPerMetablock := strconv.FormatUint(uint64(1500000000), 10)
	minGasPrice := strconv.FormatUint(uint64(200000000000), 10)
	minGasLimit := strconv.FormatUint(uint64(50000), 10)

	return &erdConfig.EconomicsConfig{
		GlobalSettings: erdConfig.GlobalSettings{
			GenesisTotalSupply: "2000000000000000000000",
			MinimumInflation:   0,
			YearSettings: []*erdConfig.YearSetting{
				{
					Year:             0,
					MaximumInflation: 0.01,
				},
			},
		},
		RewardsSettings: erdConfig.RewardsSettings{
			LeaderPercentage:    0.10,
			DeveloperPercentage: 0.10,
			CommunityPercentage: 0.10,
			CommunityAddress:    "erd1932eft30w753xyvme8d49qejgkjc09n5e49w4mwdjtm0neld797su0dlxp",
		},
		FeeSettings: erdConfig.FeeSettings{
			MaxGasLimitPerBlock:     maxGasLimitPerBlock,
			MaxGasLimitPerMetaBlock: maxGasLimitPerMetablock,
			MinGasPrice:             minGasPrice,
			MinGasLimit:             minGasLimit,
			GasPerDataByte:          "1",
			DataLimitForBaseCalc:    "10000",
		},
		ValidatorSettings: erdConfig.ValidatorSettings{
			GenesisNodePrice:         "500000000",
			UnBondPeriod:             "5",
			TotalSupply:              "200000000000",
			MinStepValue:             "100000",
			AuctionEnableNonce:       "100000",
			StakeEnableNonce:         "0",
			NumRoundsWithoutBleed:    "1000",
			MaximumPercentageToBleed: "0.5",
			BleedPercentagePerRound:  "0.00001",
			UnJailValue:              "1000",
		},
	}
}
