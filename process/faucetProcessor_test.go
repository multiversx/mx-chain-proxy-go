package process_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewFaucetProcessor_NilBaseProcessorShouldErr(t *testing.T) {

	t.Parallel()

	fp, err := process.NewFaucetProcessor(
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

	tx, err := fp.GenerateTxForSendUserFunds(senderSk, senderHexPk, senderNonce, receiver, nil, &data.NetworkConfig{})
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

	tx, err := fp.GenerateTxForSendUserFunds(senderSk, senderHexPk, senderNonce, receiver, faucetValue, &data.NetworkConfig{})
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
