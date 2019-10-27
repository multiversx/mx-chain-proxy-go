package process_test

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func testEconomicsConfig() *config.EconomicsConfig {
	return &config.EconomicsConfig{
		FeeSettings: config.FeeSettings{
			MinGasPrice: "1",
			MinGasLimit: "5",
		},
	}
}

func TestNewFaucetProcessor_NilBaseProcessorShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		nil,
		&mock.PrivateKeysLoaderStub{},
		big.NewInt(1),
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
	)

	assert.Nil(t, fp)
	assert.Equal(t, process.ErrInvalidDefaultFaucetValue, err)
}

func TestNewFaucetProcessor_EmptyAccMapShouldErr(t *testing.T) {
	t.Parallel()

	fp, err := process.NewFaucetProcessor(
		testEconomicsConfig(),
		&mock.ProcessorStub{},
		&mock.PrivateKeysLoaderStub{
			MapOfPrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				return make(map[uint32][]crypto.PrivateKey), nil
			},
		},
		big.NewInt(1),
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
			MapOfPrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
				mapToReturn := make(map[uint32][]crypto.PrivateKey)
				mapToReturn[0] = append(mapToReturn[0], nil)

				return mapToReturn, nil
			},
		},
		big.NewInt(1),
	)

	assert.NotNil(t, fp)
	assert.Nil(t, err)
}

//func TestFaucetProcessor_GenerateTxForSendUserFunds(t *testing.T) {
//	t.Parallel()
//
//	fp, _ := process.NewFaucetProcessor(
//		testEconomicsConfig(),
//		&mock.ProcessorStub{
//			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
//				return uint32(0), nil
//			},
//		},
//		&mock.PrivateKeysLoaderStub{
//			MapOfPrivateKeysByShardCalled: func() (map[uint32][]crypto.PrivateKey, error) {
//				mapToReturn := make(map[uint32][]crypto.PrivateKey)
//				mapToReturn[0] = append(mapToReturn[0], getPrivKey())
//
//				return mapToReturn, nil
//			},
//		},
//	)
//
//	tx, err := fp.GenerateTxForSendUserFunds("", big.NewInt(10))
//	assert.NotNil(t, tx)
//	assert.Nil(t, err)
//}
//
//func getPrivKey() crypto.PrivateKey {
//	keyGen := signing.NewKeyGenerator(kyber.NewBlakeSHA256Ed25519())
//	sk, _ := keyGen.GeneratePair()
//
//	return sk
//}
