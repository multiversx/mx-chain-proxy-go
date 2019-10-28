package faucet_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPrivateKeysLoader_NilAddressConverterShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(nil, &mock.ShardCoordinatorMock{}, "location")

	assert.Nil(t, pkl)
	assert.Equal(t, faucet.ErrNilAddressConverter, err)
}

func TestNewPrivateKeysLoader_NilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.AddressConverterStub{}, nil, "location")

	assert.Nil(t, pkl)
	assert.Equal(t, faucet.ErrNilShardCoordinator, err)
}

func TestNewPrivateKeysLoader_InvalidPemFileLocationShouldErr(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.AddressConverterStub{}, &mock.ShardCoordinatorMock{}, "")

	assert.Nil(t, pkl)
	assert.Equal(t, faucet.ErrInvalidPemFileLocation, err)
}

func TestNewPrivateKeysLoader_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	pkl, err := faucet.NewPrivateKeysLoader(&mock.AddressConverterStub{}, &mock.ShardCoordinatorMock{}, "location")

	assert.NotNil(t, pkl)
	assert.Nil(t, err)
}

func TestPrivateKeysLoader_MapOfPrivateKeysByShardInvalidPemFileContentShouldErr(t *testing.T) {
	t.Parallel()

	pemFileName := "wrong-test.pem"
	pkl, _ := faucet.NewPrivateKeysLoader(
		&mock.AddressConverterStub{
			CreateAddressFromPublicKeyBytesCalled: func(pubKey []byte) (state.AddressContainer, error) {
				return &mock.AddressContainerMock{BytesField: pubKey}, nil
			},
		},
		&mock.ShardCoordinatorMock{},
		pemFileName,
	)

	err := ioutil.WriteFile(pemFileName, []byte(getWrongPemFileContent()), 0644)
	assert.Nil(t, err)

	defer func() {
		_ = os.Remove(pemFileName)
	}()

	retMap, err := pkl.MapOfPrivateKeysByShard()
	assert.Nil(t, retMap)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "invalid"))
}

func TestPrivateKeysLoader_MapOfPrivateKeysByShardShouldWork(t *testing.T) {
	t.Parallel()

	pemFileName := "test.pem"
	pkl, _ := faucet.NewPrivateKeysLoader(
		&mock.AddressConverterStub{
			CreateAddressFromPublicKeyBytesCalled: func(pubKey []byte) (state.AddressContainer, error) {
				return &mock.AddressContainerMock{BytesField: pubKey}, nil
			},
		},
		&mock.ShardCoordinatorMock{},
		pemFileName,
	)

	err := ioutil.WriteFile(pemFileName, []byte(getTestPemFileContent()), 0644)
	assert.Nil(t, err)

	defer func() {
		_ = os.Remove(pemFileName)
	}()

	retMap, err := pkl.MapOfPrivateKeysByShard()
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
