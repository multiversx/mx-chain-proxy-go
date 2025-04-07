package process_test

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/core/sharding"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccountProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(nil, &mock.PubKeyConverterMock{})

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewAccountProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, nil)

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewAccountProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})

	assert.NotNil(t, ap)
	assert.Nil(t, err)
}

//------- GetAccount

func TestAccountProcessor_GetAccountInvalidHexAddressShouldErr(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	accnt, err := ap.GetAccount("invalid hex number", common.AccountQueryOptions{})

	assert.Nil(t, accnt)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestAccountProcessor_GetAccountComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address, common.AccountQueryOptions{})

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return nil, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address, common.AccountQueryOptions{})

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address1", ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				return 0, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address, common.AccountQueryOptions{})

	assert.Nil(t, accnt)
	assert.True(t, errors.Is(err, process.ErrSendingRequest))
}

func TestAccountProcessor_GetAccountSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := errors.New("expected error")
	respondedAccount := &data.AccountModel{
		Account: data.Account{
			Address: "an address",
		},
	}
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: addressFail, ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := value.(*data.AccountApiResponse)
				valRespond.Data.Account = respondedAccount.Account
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accountModel, err := ap.GetAccount(address, common.AccountQueryOptions{})

	assert.Equal(t, respondedAccount.Account, accountModel.Account)
	assert.Nil(t, err)
}

func TestAccountProcessor_GetValueForAKeyShouldWork(t *testing.T) {
	t.Parallel()

	expectedValue := "dummyValue"
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				valRespond := value.(*data.AccountKeyValueResponse)
				valRespond.Data.Value = expectedValue
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	key := "key"
	addr1 := "DEADBEEF"
	value, err := ap.GetValueForKey(addr1, key, common.AccountQueryOptions{})
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestAccountProcessor_GetValueForAKeyShouldError(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("err")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				return 0, expectedError
			},
		},
		&mock.PubKeyConverterMock{},
	)

	key := "key"
	addr1 := "DEADBEEF"
	value, err := ap.GetValueForKey(addr1, key, common.AccountQueryOptions{})
	assert.Equal(t, "", value)
	assert.True(t, errors.Is(err, process.ErrSendingRequest))
}

func TestAccountProcessor_GetShardIForAddressShouldWork(t *testing.T) {
	t.Parallel()

	shardC, err := sharding.NewMultiShardCoordinator(uint32(2), 0)
	require.NoError(t, err)

	bech32C, _ := pubkeyConverter.NewBech32PubkeyConverter(32, "erd")

	// this addressShard0 should be in shard 0 for a 2 shards configuration
	addressShard0 := "erd1ffqlrryvwrnfh2523wmzrhvx5d8p2wmxeau64fps4lnqq5qex68q7ax8k5"

	// this addressShard1 should be in shard 1 for a 2 shards configuration
	addressShard1 := "erd1qqe9qll7n66lv4cuuml2wxsv3sd2t0eyajkyjr7rvtqmhha0cgsse4pel3"

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return shardC.ComputeId(addressBuff), nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return observers, nil
			},
		},
		bech32C,
	)

	shardID, err := ap.GetShardIDForAddress(addressShard1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), shardID)

	shardID, err = ap.GetShardIDForAddress(addressShard0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), shardID)
}

func TestAccountProcessor_GetShardIDForAddressShouldError(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("err")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, expectedError
			},
		},
		&mock.PubKeyConverterMock{},
	)

	shardID, err := ap.GetShardIDForAddress("aaaa")
	assert.Equal(t, uint32(0), shardID)
	assert.Equal(t, expectedError, err)
}

func TestAccountProcessor_GetESDTsWithRoleGetObserversFails(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("cannot get observers")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			GetObserversCalled: func(_ uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return nil, expectedErr
			},
		},
		&mock.PubKeyConverterMock{},
	)

	result, err := ap.GetESDTsWithRole("address", "role", common.AccountQueryOptions{})
	require.Equal(t, expectedErr, err)
	require.Nil(t, result)
}

func TestAccountProcessor_GetESDTsWithRoleApiCallFails(t *testing.T) {
	t.Parallel()

	expectedApiErr := errors.New("cannot get observers")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{
						Address: "observer0",
						ShardId: core.MetachainShardId,
					},
				}, nil
			},

			CallGetRestEndPointCalled: func(_ string, _ string, _ interface{}) (int, error) {
				return 0, expectedApiErr
			},
		},
		&mock.PubKeyConverterMock{},
	)

	result, err := ap.GetESDTsWithRole("address", "role", common.AccountQueryOptions{})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "sending request error"))
	require.Nil(t, result)
}

func TestAccountProcessor_GetESDTsWithRoleShouldWork(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(_ []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				tokensResponse := value.(*data.GenericAPIResponse)
				tokensResponse.Data = []string{"token0"}
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	response, err := ap.GetESDTsWithRole(address, "role", common.AccountQueryOptions{})
	require.NoError(t, err)
	require.Equal(t, "token0", response.Data.([]string)[0])
}

func TestAccountProcessor_GetESDTsRolesGetObserversFails(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("cannot get observers")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return nil, expectedErr
			},
		},
		&mock.PubKeyConverterMock{},
	)

	result, err := ap.GetESDTsRoles("address", common.AccountQueryOptions{})
	require.Equal(t, expectedErr, err)
	require.Nil(t, result)
}

func TestAccountProcessor_GetESDTsRolesApiCallFails(t *testing.T) {
	t.Parallel()

	expectedApiErr := errors.New("cannot get observers")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{
						Address: "observer0",
						ShardId: core.MetachainShardId,
					},
				}, nil
			},

			CallGetRestEndPointCalled: func(_ string, _ string, _ interface{}) (int, error) {
				return 0, expectedApiErr
			},
		},
		&mock.PubKeyConverterMock{},
	)

	result, err := ap.GetESDTsRoles("address", common.AccountQueryOptions{})
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "sending request error"))
	require.Nil(t, result)
}

func TestAccountProcessor_GetESDTsRolesShouldWork(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(_ []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				tokensResponse := value.(*data.GenericAPIResponse)
				tokensResponse.Data = []string{"token0"}
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	response, err := ap.GetESDTsRoles(address, common.AccountQueryOptions{})
	require.NoError(t, err)
	require.Equal(t, "token0", response.Data.([]string)[0])
}

func TestAccountProcessor_GetCodeHash(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(_ []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				codeHashResponse := value.(*data.GenericAPIResponse)
				codeHashResponse.Data = []string{"code-hash"}
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	response, err := ap.GetCodeHash(address, common.AccountQueryOptions{})
	require.NoError(t, err)
	require.Equal(t, "code-hash", response.Data.([]string)[0])
}

func TestAccountProcessor_IsDataTrieMigrated(t *testing.T) {
	t.Parallel()

	t.Run("should return error when cannot get observers", func(t *testing.T) {
		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return nil, errors.New("cannot get observers")
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.IsDataTrieMigrated("address", common.AccountQueryOptions{})
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("should return error when cannot get data trie migrated", func(t *testing.T) {
		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{
							Address: "observer0",
							ShardId: 0,
						},
					}, nil
				},

				CallGetRestEndPointCalled: func(_ string, _ string, _ interface{}) (int, error) {
					return 0, errors.New("cannot get data trie migrated")
				},
				ComputeShardIdCalled: func(_ []byte) (uint32, error) {
					return 0, nil
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.IsDataTrieMigrated("DEADBEEF", common.AccountQueryOptions{})
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("should work", func(t *testing.T) {
		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{
							Address: "observer0",
							ShardId: 0,
						},
					}, nil
				},

				CallGetRestEndPointCalled: func(_ string, _ string, value interface{}) (int, error) {
					dataTrieMigratedResponse := value.(*data.GenericAPIResponse)
					dataTrieMigratedResponse.Data = true
					return 0, nil
				},
				ComputeShardIdCalled: func(_ []byte) (uint32, error) {
					return 0, nil
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.IsDataTrieMigrated("DEADBEEF", common.AccountQueryOptions{})
		require.NoError(t, err)
		require.True(t, result.Data.(bool))
	})
}

func TestAccountProcessor_GetAccounts(t *testing.T) {
	t.Parallel()

	t.Run("should return error if a shard returns error", func(t *testing.T) {
		t.Parallel()

		expectedError := "expected error message"
		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(shardID uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					address := "observer0"
					if shardID == 1 {
						address = "observer1"
					}
					return []*data.NodeData{
						{
							Address: address,
							ShardId: shardID,
						},
					}, nil
				},

				CallPostRestEndPointCalled: func(obsAddr string, _ string, _ interface{}, value interface{}) (int, error) {
					response := value.(*data.AccountsApiResponse)
					if obsAddr == "observer1" {
						response.Error = expectedError
					}
					return 0, nil
				},
				ComputeShardIdCalled: func(addr []byte) (uint32, error) {
					if hex.EncodeToString(addr) == "aabb" {
						return 0, nil
					}

					return 1, nil
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.GetAccounts([]string{"aabb", "bbaa"}, common.AccountQueryOptions{})
		require.Equal(t, expectedError, err.Error())
		require.Empty(t, result)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(shardID uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					address := "observer0"
					if shardID == 1 {
						address = "observer1"
					}
					return []*data.NodeData{
						{
							Address: address,
							ShardId: shardID,
						},
					}, nil
				},

				CallPostRestEndPointCalled: func(obsAddr string, _ string, _ interface{}, value interface{}) (int, error) {
					address := "shard0Address"
					if obsAddr == "observer1" {
						address = "shard1Address"
					}
					response := value.(*data.AccountsApiResponse)
					response.Data.Accounts = map[string]*data.Account{
						address: {Address: address, Balance: "37"},
					}
					return 0, nil
				},
				ComputeShardIdCalled: func(addr []byte) (uint32, error) {
					if hex.EncodeToString(addr) == "aabb" {
						return 0, nil
					}

					return 1, nil
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.GetAccounts([]string{"aabb", "bbaa"}, common.AccountQueryOptions{})
		require.NoError(t, err)

		require.Equal(t, map[string]*data.Account{
			"shard0Address": {
				Address: "shard0Address",
				Balance: "37",
			},
			"shard1Address": {
				Address: "shard1Address",
				Balance: "37",
			},
		}, result.Accounts)
	})
}

func TestAccountProcessor_IterateKeys(t *testing.T) {
	t.Parallel()

	t.Run("should return error when cannot get observers", func(t *testing.T) {
		t.Parallel()

		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return nil, errors.New("cannot get observers")
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.IterateKeys("address", 0, nil, common.AccountQueryOptions{})
		require.Error(t, err)
		require.Nil(t, result)
	})
	t.Run("should return error observers return error", func(t *testing.T) {
		t.Parallel()

		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{
							Address: "observer0",
							ShardId: 0,
						},
					}, nil
				},
				CallPostRestEndPointCalled: func(address string, _ string, _ interface{}, _ interface{}) (int, error) {
					return 0, errors.New("cannot iterate keys")
				},
				ComputeShardIdCalled: func(_ []byte) (uint32, error) {
					return 0, nil
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.IterateKeys("DEADBEEF", 10, [][]byte{[]byte("iterator state")}, common.AccountQueryOptions{})
		require.Error(t, err)
		require.Nil(t, result)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		keyPairs := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		newIteratorState := [][]byte{[]byte("new iterator state")}
		ap, _ := process.NewAccountProcessor(
			&mock.ProcessorStub{
				GetObserversCalled: func(_ uint32, _ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
					return []*data.NodeData{
						{
							Address: "observer0",
							ShardId: 0,
						},
					}, nil
				},

				CallPostRestEndPointCalled: func(address string, path string, iteratorState interface{}, response interface{}) (int, error) {
					assert.Equal(t, "/address/iterate-keys", path)
					iterateKeysResponse := response.(*data.GenericAPIResponse)
					iterateKeysResponse.Data = map[string]interface{}{
						"pairs":         keyPairs,
						"iteratorState": newIteratorState,
					}
					return 0, nil
				},
				ComputeShardIdCalled: func(_ []byte) (uint32, error) {
					return 0, nil
				},
			},
			&mock.PubKeyConverterMock{},
		)

		result, err := ap.IterateKeys("DEADBEEF", 10, [][]byte{[]byte("original iterator state")}, common.AccountQueryOptions{})
		require.NoError(t, err)
		responseMap, ok := result.Data.(map[string]interface{})
		assert.True(t, ok)

		respPairsMap, ok := responseMap["pairs"].(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, keyPairs, respPairsMap)

		respIterState, ok := responseMap["iteratorState"].([][]byte)
		assert.True(t, ok)
		assert.Equal(t, newIteratorState, respIterState)
	})
}
