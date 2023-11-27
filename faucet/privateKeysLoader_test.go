package faucet_test

import (
	"os"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/faucet"
	"github.com/multiversx/mx-chain-proxy-go/faucet/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPrivateKeysLoader_NilShardCoordinatorShouldErr(t *testing.T) {
	pemFileName := "testWallet.pem"
	err := os.WriteFile(pemFileName, []byte(getWrongPemFileContent()), 0644)
	require.Nil(t, err)
	defer func() {
		_ = os.Remove(pemFileName)
	}()

	pkl, err := faucet.NewPrivateKeysLoader(nil, pemFileName, &mock.PubKeyConverterMock{})

	require.Nil(t, pkl)
	require.Equal(t, faucet.ErrNilShardCoordinator, err)
}

func TestNewPrivateKeysLoader_PemFileNotFoundShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.ShardCoordinatorMock{}, "", &mock.PubKeyConverterMock{})

	require.Nil(t, pkl)
	require.Equal(t, faucet.ErrFaucetPemFileDoesNotExist, err)
}

func TestNewPrivateKeysLoader_NilPubKeyConverterShouldErr(t *testing.T) {
	pemFileName := "testWallet.pem"
	err := os.WriteFile(pemFileName, []byte(getWrongPemFileContent()), 0644)
	require.Nil(t, err)
	defer func() {
		_ = os.Remove(pemFileName)
	}()
	pkl, err := faucet.NewPrivateKeysLoader(&mock.ShardCoordinatorMock{}, pemFileName, nil)

	require.Nil(t, pkl)
	require.Equal(t, faucet.ErrNilPubKeyConverter, err)
}

func TestNewPrivateKeysLoader_OkValsShouldWork(t *testing.T) {
	pemFileName := "testWallet.pem"
	err := os.WriteFile(pemFileName, []byte(getWrongPemFileContent()), 0644)
	require.Nil(t, err)
	defer func() {
		_ = os.Remove(pemFileName)
	}()
	pkl, err := faucet.NewPrivateKeysLoader(&mock.ShardCoordinatorMock{}, pemFileName, &mock.PubKeyConverterMock{})

	require.NotNil(t, pkl)
	require.Nil(t, err)
}

func TestPrivateKeysLoader_MapOfPrivateKeysByShardInvalidPemFileContentShouldErr(t *testing.T) {
	pemFileName := "wrong-test.pem"
	err := os.WriteFile(pemFileName, []byte(getWrongPemFileContent()), 0644)
	require.Nil(t, err)
	defer func() {
		_ = os.Remove(pemFileName)
	}()

	pkl, _ := faucet.NewPrivateKeysLoader(
		&mock.ShardCoordinatorMock{},
		pemFileName,
		&mock.PubKeyConverterMock{},
	)

	retMap, err := pkl.PrivateKeysByShard()
	require.Nil(t, retMap)
	require.NotNil(t, err)
	require.True(t, strings.Contains(err.Error(), "invalid"))
}

func TestPrivateKeysLoader_MapOfPrivateKeysByShardShouldWork(t *testing.T) {
	pemFileName := "test.pem"
	err := os.WriteFile(pemFileName, []byte(getTestPemFileContent()), 0644)
	require.Nil(t, err)
	defer func() {
		_ = os.Remove(pemFileName)
	}()

	pkl, err := faucet.NewPrivateKeysLoader(
		&mock.ShardCoordinatorMock{},
		pemFileName,
		&mock.PubKeyConverterMock{},
	)
	require.NoError(t, err)

	retMap, err := pkl.PrivateKeysByShard()
	require.NotNil(t, retMap)
	require.Nil(t, err)
}

func getTestPemFileContent() string {
	return `-----BEGIN PRIVATE KEY for 8a2ee461bd72652fc33d8705b33cf240dc4a1531c5bf80bd4f4d92b6d83636ae-----
YTQwZTQ1YzY2YmM5NTA2YjY1ZWZjNjc2YTlkODRhNGRmMzk1ZWVhYzMwZGI2MjA2
NzlmYWIxNzQ4ZDZlNzIwMQ==
-----END PRIVATE KEY for 8a2ee461bd72652fc33d8705b33cf240dc4a1531c5bf80bd4f4d92b6d83636ae-----
-----BEGIN PRIVATE KEY for e45b0dcd13663a6a21d9252882556966400d7940c4112e5b9d28d72aa1e853fd-----
ODU0OGJjYmE3MDMwMDhkNDVkMTgyODRlZjZmZjYxYjM1ZGRiNmY0N2VhZTVmYTQx
YzQ0ZDY5MTQ1MDUzYzEwYg==
-----END PRIVATE KEY for e45b0dcd13663a6a21d9252882556966400d7940c4112e5b9d28d72aa1e853fd-----
`
}

func getWrongPemFileContent() string {
	return `invalid pem file content`
}
