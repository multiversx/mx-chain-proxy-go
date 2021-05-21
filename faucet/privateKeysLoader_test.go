package faucet_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/faucet"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPrivateKeysLoader_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(nil, "location", &mock.PubKeyConverterMock{})

	assert.Nil(t, pkl)
	assert.Equal(t, faucet.ErrNilShardCoordinator, err)
}

func TestNewPrivateKeysLoader_InvalidPemFileLocationShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.ShardCoordinatorMock{}, "", &mock.PubKeyConverterMock{})

	assert.Nil(t, pkl)
	assert.Equal(t, faucet.ErrFaucetPemFileDoesNotExist, err)
}

func TestNewPrivateKeysLoader_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.ShardCoordinatorMock{}, "location", nil)

	assert.Nil(t, pkl)
	assert.Equal(t, faucet.ErrNilPubKeyConverter, err)
}

func TestNewPrivateKeysLoader_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.ShardCoordinatorMock{}, "location", &mock.PubKeyConverterMock{})

	assert.NotNil(t, pkl)
	assert.Nil(t, err)
}

func TestPrivateKeysLoader_MapOfPrivateKeysByShardInvalidPemFileContentShouldErr(t *testing.T) {
	t.Parallel()

	pemFileName := "wrong-test.pem"
	pkl, _ := faucet.NewPrivateKeysLoader(
		&mock.ShardCoordinatorMock{},
		pemFileName,
		&mock.PubKeyConverterMock{},
	)

	err := ioutil.WriteFile(pemFileName, []byte(getWrongPemFileContent()), 0644)
	assert.Nil(t, err)

	defer func() {
		_ = os.Remove(pemFileName)
	}()

	retMap, err := pkl.PrivateKeysByShard()
	assert.Nil(t, retMap)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "invalid"))
}

func TestPrivateKeysLoader_MapOfPrivateKeysByShardShouldWork(t *testing.T) {
	t.Parallel()

	pemFileName := "test.pem"
	pkl, _ := faucet.NewPrivateKeysLoader(
		&mock.ShardCoordinatorMock{},
		pemFileName,
		&mock.PubKeyConverterMock{},
	)

	err := ioutil.WriteFile(pemFileName, []byte(getTestPemFileContent()), 0644)
	assert.Nil(t, err)

	defer func() {
		_ = os.Remove(pemFileName)
	}()

	retMap, err := pkl.PrivateKeysByShard()
	assert.NotNil(t, retMap)
	assert.Nil(t, err)
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
