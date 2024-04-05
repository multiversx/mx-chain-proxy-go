package facade_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/facade"
	"github.com/multiversx/mx-chain-proxy-go/facade/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Nonce uint64
	Hash  string
}

var publicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "erd")

func TestNewProxyFacade_NilActionsProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		nil,
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilActionsProcessor, err)
}

func TestNewProxyFacade_NilAccountProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		nil,
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilAccountProcessor, err)
}

func TestNewProxyFacade_NilTransactionProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		nil,
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilTransactionProcessor, err)
}

func TestNewProxyFacade_NilGetValuesProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		nil,
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilSCQueryService, err)
}

func TestNewProxyFacade_NilNodeGroupProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		nil,
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilNodeGroupProcessor, err)
}

func TestNewProxyFacade_NilValStatsProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		nil,
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilValidatorStatisticsProcessor, err)
}

func TestNewProxyFacade_NilFaucetProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		nil,
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilFaucetProcessor, err)
}

func TestNewProxyFacade_NilNodeProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		nil,
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilNodeStatusProcessor, err)
}

func TestNewProxyFacade_NilBlocksProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		nil,
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilBlocksProcessor, err)
}

func TestNewProxyFacade_NilProofProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		nil,
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilProofProcessor, err)
}

func TestNewProxyFacade_NilStatusProcessorShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		nil,
		&mock.AboutInfoProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilStatusProcessor, err)
}

func TestNewProxyFacade_NilAboutInfoProcessorShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		nil,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilAboutInfoProcessor, err)
}

func TestNewProxyFacade_ShouldWork(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	assert.NotNil(t, epf)
	assert.Nil(t, err)
}

func TestNewProxyFacade_GetBlocksByRound(t *testing.T) {
	t.Parallel()

	expectedResponse := &data.BlocksApiResponse{
		Data: data.BlocksApiResponsePayload{
			Blocks: []*api.Block{
				{
					Round: 4,
					Hash:  "hash1",
				},
				{
					Round: 4,
					Hash:  "hash2",
				},
			},
		},
	}

	errGetBlockByRound := errors.New("could not get block by round")
	epf, err := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{
			GetBlocksByRoundCalled: func(round uint64, _ common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
				if round == 4 {
					return expectedResponse, nil
				}
				return nil, errGetBlockByRound
			},
		},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)
	require.NoError(t, err)

	ret, err := epf.GetBlocksByRound(3, common.BlockQueryOptions{WithTransactions: true})
	require.Equal(t, errGetBlockByRound, err)
	require.Nil(t, ret)

	ret, err = epf.GetBlocksByRound(4, common.BlockQueryOptions{WithTransactions: true})
	require.Nil(t, err)
	require.Equal(t, expectedResponse, ret)
}

func TestProxyFacade_GetAccount(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string, options common.AccountQueryOptions) (account *data.AccountModel, e error) {
				wasCalled = true
				return &data.AccountModel{}, nil
			},
		},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	_, _ = epf.GetAccount("", common.AccountQueryOptions{})

	assert.True(t, wasCalled)
}

func TestProxyFacade_SendTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true

				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	_, _, _ = epf.SendTransaction(&data.Transaction{})

	assert.True(t, wasCalled)
}

func TestProxyFacade_SimulateTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SimulateTransactionCalled: func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
				wasCalled = true
				return nil, nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	_, _ = epf.SimulateTransaction(&data.Transaction{}, false)

	assert.True(t, wasCalled)
}

func TestProxyFacade_SendUserFunds(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
				return &data.AccountModel{
					Account: data.Account{
						Nonce: uint64(0),
					},
				}, nil
			},
		},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true
				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{
			SenderDetailsFromPemCalled: func(receiver string) (crypto.PrivateKey, string, error) {
				return getPrivKey(), "rcvr", nil
			},
			GenerateTxForSendUserFundsCalled: func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64, receiver string, value *big.Int, config *data.NetworkConfig) (*data.Transaction, error) {
				return &data.Transaction{}, nil
			},
		},
		&mock.NodeStatusProcessorStub{
			GetConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
				return &data.GenericAPIResponse{
					Data: map[string]interface{}{
						"config": map[string]interface{}{
							"erd_chain_id":                "chainID",
							"erd_min_transaction_version": 1.0,
						},
					},
				}, nil
			},
		},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	_ = epf.SendUserFunds("", big.NewInt(0))

	assert.True(t, wasCalled)
}

func TestProxyFacade_GetDataValue(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{
			ExecuteQueryCalled: func(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
				wasCalled = true
				return &vm.VMOutputApi{}, data.BlockInfo{}, nil
			},
		},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	_, _, _ = epf.ExecuteSCQuery(nil)

	assert.True(t, wasCalled)
}

func TestProxyFacade_GetHeartbeatData(t *testing.T) {
	t.Parallel()

	expectedResults := &data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				ReceivedShardID: 0,
				ComputedShardID: 1,
			},
		},
	}
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{
			GetHeartbeatDataCalled: func() (*data.HeartbeatResponse, error) {
				return expectedResults, nil
			},
		},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, _ := epf.GetHeartbeatData()

	assert.Equal(t, expectedResults, actualResult)
}

func TestProxyFacade_ReloadObservers(t *testing.T) {
	t.Parallel()

	expectedResult := data.NodesReloadResponse{
		Description: "abc",
		Error:       "bca",
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{
			ReloadObserversCalled: func() data.NodesReloadResponse {
				return expectedResult
			},
		},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult := epf.ReloadObservers()

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_ReloadFullHistoryObservers(t *testing.T) {
	t.Parallel()

	expectedResult := data.NodesReloadResponse{
		Description: "abc",
		Error:       "bca",
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{
			ReloadFullHistoryObserversCalled: func() data.NodesReloadResponse {
				return expectedResult
			},
		},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult := epf.ReloadFullHistoryObservers()

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_GetBlockByHash(t *testing.T) {
	t.Parallel()

	expectedResult := &data.BlockApiResponse{
		Data: data.BlockApiResponsePayload{
			Block: api.Block{
				Nonce: 10,
				Round: 10,
			},
		},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetBlockByHashCalled: func(_ uint32, _ string, _ common.BlockQueryOptions) (*data.BlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetBlockByHash(0, "aaaa", common.BlockQueryOptions{})
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_GetBlockByNonce(t *testing.T) {
	t.Parallel()

	expectedResult := &data.BlockApiResponse{
		Data: data.BlockApiResponsePayload{
			Block: api.Block{
				Nonce: 10,
				Round: 10,
			},
		},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetBlockByNonceCalled: func(_ uint32, _ uint64, _ common.BlockQueryOptions) (*data.BlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetBlockByNonce(0, 10, common.BlockQueryOptions{})
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

// Internal Blocks

func TestProxyFacade_GetInternalBlockByHash(t *testing.T) {
	t.Parallel()

	expectedResult := &data.InternalBlockApiResponse{
		Data: data.InternalBlockApiResponsePayload{
			Block: &testStruct{
				Nonce: 10,
				Hash:  "aaaa",
			},
		},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetInternalBlockByHashCalled: func(_ uint32, _ string, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetInternalBlockByHash(0, "aaaa", common.Internal)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_GetInternalBlockByNonce(t *testing.T) {
	t.Parallel()

	expectedResult := &data.InternalBlockApiResponse{
		Data: data.InternalBlockApiResponsePayload{
			Block: &testStruct{
				Nonce: 10,
				Hash:  "aaaa",
			},
		},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetInternalBlockByNonceCalled: func(_ uint32, _ uint64, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetInternalBlockByNonce(0, 10, common.Internal)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_GetInternalMiniBlockByHash(t *testing.T) {
	t.Parallel()

	expectedResult := &data.InternalMiniBlockApiResponse{
		Data: data.InternalMiniBlockApiResponsePayload{
			MiniBlock: &testStruct{
				Nonce: 10,
				Hash:  "aaaa",
			},
		},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetInternalMiniBlockByHashCalled: func(_ uint32, _ string, epoch uint32, _ common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetInternalMiniBlockByHash(0, "aaaa", 1, common.Internal)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_GetRatingsConfig(t *testing.T) {
	t.Parallel()

	expectedResult := &data.GenericAPIResponse{
		Data: &testStruct{
			Nonce: 0,
			Hash:  "aaaa",
		},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{
			GetRatingsConfigCalled: func() (*data.GenericAPIResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetRatingsConfig()
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestProxyFacade_GetTransactionsPool(t *testing.T) {
	t.Parallel()

	providedNonce := uint64(5)
	providedTx := data.WrappedTransaction{
		TxFields: map[string]interface{}{
			"sender": "sender",
			"nonce":  providedNonce,
			"hash":   "hash",
		},
	}
	providedNonceGap := data.NonceGap{
		From: 6,
		To:   20,
	}
	expectedTxPool := &data.TransactionsPool{
		RegularTransactions: []data.WrappedTransaction{providedTx},
	}
	expectedTxPoolForSender := &data.TransactionsPoolForSender{
		Transactions: []data.WrappedTransaction{providedTx},
	}
	expectedNonceGaps := &data.TransactionsPoolNonceGaps{
		Gaps: []data.NonceGap{providedNonceGap},
	}

	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			GetTransactionsPoolCalled: func(fields string) (*data.TransactionsPool, error) {
				return expectedTxPool, nil
			},
			GetTransactionsPoolForShardCalled: func(shardID uint32, fields string) (*data.TransactionsPool, error) {
				return expectedTxPool, nil
			},
			GetTransactionsPoolForSenderCalled: func(sender, fields string) (*data.TransactionsPoolForSender, error) {
				return expectedTxPoolForSender, nil
			},
			GetLastPoolNonceForSenderCalled: func(sender string) (uint64, error) {
				return providedNonce, nil
			},
			GetTransactionsPoolNonceGapsForSenderCalled: func(sender string) (*data.TransactionsPoolNonceGaps, error) {
				return expectedNonceGaps, nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualTxPool, err := epf.GetTransactionsPool("")
	require.Nil(t, err)
	assert.Equal(t, expectedTxPool, actualTxPool)

	actualTxPool, err = epf.GetTransactionsPoolForShard(0, "")
	require.Nil(t, err)
	assert.Equal(t, expectedTxPool, actualTxPool)

	actualTxPoolForSender, err := epf.GetTransactionsPoolForSender("", "")
	require.Nil(t, err)
	assert.Equal(t, expectedTxPoolForSender, actualTxPoolForSender)

	actualNonce, err := epf.GetLastPoolNonceForSender("")
	require.Nil(t, err)
	assert.Equal(t, providedNonce, actualNonce)

	actualNonceGaps, err := epf.GetTransactionsPoolNonceGapsForSender("")
	require.Nil(t, err)
	assert.Equal(t, expectedNonceGaps, actualNonceGaps)
}

func TestProxyFacade_GetGasConfigs(t *testing.T) {
	t.Parallel()

	expectedResult := &data.GenericAPIResponse{
		Data: &testStruct{
			Hash: "aaaa",
		},
	}

	wasCalled := false
	epf, _ := facade.NewProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.NodeGroupProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{
			GetGasConfigsCalled: func() (*data.GenericAPIResponse, error) {
				wasCalled = true
				return expectedResult, nil
			},
		},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
		&mock.StatusProcessorStub{},
		&mock.AboutInfoProcessorStub{},
	)

	actualResult, err := epf.GetGasConfigs()
	require.Nil(t, err)

	assert.True(t, wasCalled)
	assert.Equal(t, expectedResult, actualResult)
}

func getPrivKey() crypto.PrivateKey {
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	sk, _ := keyGen.GeneratePair()

	return sk
}
