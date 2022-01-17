package facade_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/facade/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Nonce uint64
	Hash  string
}

var publicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32)

func TestNewElrondProxyFacade_NilActionsProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		nil,
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilActionsProcessor, err)
}

func TestNewElrondProxyFacade_NilAccountProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		nil,
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilAccountProcessor, err)
}

func TestNewElrondProxyFacade_NilTransactionProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		nil,
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilTransactionProcessor, err)
}

func TestNewElrondProxyFacade_NilGetValuesProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		nil,
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilSCQueryService, err)
}

func TestNewElrondProxyFacade_NilHeartbeatProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
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
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilHeartbeatProcessor, err)
}

func TestNewElrondProxyFacade_NilValStatsProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		nil,
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilValidatorStatisticsProcessor, err)
}

func TestNewElrondProxyFacade_NilFaucetProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		nil,
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilFaucetProcessor, err)
}

func TestNewElrondProxyFacade_NilNodeProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		nil,
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilNodeStatusProcessor, err)
}

func TestNewElrondProxyFacade_NilBlocksProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		nil,
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilBlocksProcessor, err)
}

func TestNewElrondProxyFacade_NilProofProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		nil,
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilProofProcessor, err)
}

func TestNewElrondProxyFacade_ShouldWork(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	assert.NotNil(t, epf)
	assert.Nil(t, err)
}

func TestNewElrondProxyFacade_GetBlocksByRound(t *testing.T) {
	t.Parallel()

	expectedResponse := &data.BlocksApiResponse{
		Data: data.BlocksApiResponsePayload{
			Blocks: []*data.Block{
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
	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{
			GetBlocksByRoundCalled: func(round uint64, _ bool) (*data.BlocksApiResponse, error) {
				if round == 4 {
					return expectedResponse, nil
				}
				return nil, errGetBlockByRound
			},
		},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	ret, err := epf.GetBlocksByRound(3, true)
	require.Equal(t, errGetBlockByRound, err)
	require.Nil(t, ret)

	ret, err = epf.GetBlocksByRound(4, true)
	require.Nil(t, err)
	require.Equal(t, expectedResponse, ret)
}

func TestElrondProxyFacade_GetAccount(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (account *data.Account, e error) {
				wasCalled = true
				return &data.Account{}, nil
			},
		},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	_, _ = epf.GetAccount("")

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true

				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	_, _, _ = epf.SendTransaction(&data.Transaction{})

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SimulateTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SimulateTransactionCalled: func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
				wasCalled = true
				return nil, nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	_, _ = epf.SimulateTransaction(&data.Transaction{}, false)

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendUserFunds(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (*data.Account, error) {
				return &data.Account{
					Nonce: uint64(0),
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
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{
			SenderDetailsFromPemCalled: func(receiver string) (crypto.PrivateKey, string, error) {
				return getPrivKey(), "rcvr", nil
			},
			GenerateTxForSendUserFundsCalled: func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64, receiver string, value *big.Int, chainID string, version uint32) (*data.Transaction, error) {
				return &data.Transaction{}, nil
			},
		},
		&mock.NodeStatusProcessorStub{
			GetConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
				return &data.GenericAPIResponse{
					Data: map[string]interface{}{
						"config": map[string]interface{}{
							core.MetricChainId:               "chainID",
							core.MetricMinTransactionVersion: 1.0,
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
	)

	_ = epf.SendUserFunds("", big.NewInt(0))

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetDataValue(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{
			ExecuteQueryCalled: func(query *data.SCQuery) (*vm.VMOutputApi, error) {
				wasCalled = true
				return &vm.VMOutputApi{}, nil
			},
		},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	_, _ = epf.ExecuteSCQuery(nil)

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetHeartbeatData(t *testing.T) {
	t.Parallel()

	expectedResults := &data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				ReceivedShardID: 0,
				ComputedShardID: 1,
			},
		},
	}
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{
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
	)

	actualResult, _ := epf.GetHeartbeatData()

	assert.Equal(t, expectedResults, actualResult)
}

func TestElrondProxyFacade_ReloadObservers(t *testing.T) {
	t.Parallel()

	expectedResult := data.NodesReloadResponse{
		Description: "abc",
		Error:       "bca",
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{
			ReloadObserversCalled: func() data.NodesReloadResponse {
				return expectedResult
			},
		},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult := epf.ReloadObservers()

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_ReloadFullHistoryObservers(t *testing.T) {
	t.Parallel()

	expectedResult := data.NodesReloadResponse{
		Description: "abc",
		Error:       "bca",
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{
			ReloadFullHistoryObserversCalled: func() data.NodesReloadResponse {
				return expectedResult
			},
		},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult := epf.ReloadFullHistoryObservers()

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_GetBlockByHash(t *testing.T) {
	t.Parallel()

	expectedResult := &data.BlockApiResponse{
		Data: data.BlockApiResponsePayload{
			Block: data.Block{
				Nonce: 10,
				Round: 10,
			},
		},
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetBlockByHashCalled: func(shardId uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult, err := epf.GetBlockByHash(0, "aaaa", false)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_GetBlockByNonce(t *testing.T) {
	t.Parallel()

	expectedResult := &data.BlockApiResponse{
		Data: data.BlockApiResponsePayload{
			Block: data.Block{
				Nonce: 10,
				Round: 10,
			},
		},
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetBlockByNonceCalled: func(shardId uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult, err := epf.GetBlockByNonce(0, 10, false)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

// Internal Blocks

func TestElrondProxyFacade_GetInternalBlockByHash(t *testing.T) {
	t.Parallel()

	expectedResult := &data.InternalBlockApiResponse{
		Data: data.InternalBlockApiResponsePayload{
			Block: &testStruct{
				Nonce: 10,
				Hash:  "aaaa",
			},
		},
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetInternalBlockByHashCalled: func(shardId uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult, err := epf.GetInternalBlockByHash(0, "aaaa", common.Internal)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_GetInternalBlockByNonce(t *testing.T) {
	t.Parallel()

	expectedResult := &data.InternalBlockApiResponse{
		Data: data.InternalBlockApiResponsePayload{
			Block: &testStruct{
				Nonce: 10,
				Hash:  "aaaa",
			},
		},
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetInternalBlockByNonceCalled: func(shardId uint32, nonce uint64, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult, err := epf.GetInternalBlockByNonce(0, 10, common.Internal)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_GetInternalMiniBlockByHash(t *testing.T) {
	t.Parallel()

	expectedResult := &data.InternalBlockApiResponse{
		Data: data.InternalBlockApiResponsePayload{
			Block: &testStruct{
				Nonce: 10,
				Hash:  "aaaa",
			},
		},
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{
			GetInternalMiniBlockByHashCalled: func(shardId uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
				return expectedResult, nil
			},
		},
		&mock.BlocksProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
		&mock.ESDTSuppliesProcessorStub{},
	)

	actualResult, err := epf.GetInternalMiniBlockByHash(0, "aaaa", common.Internal)
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_GetRatingsConfig(t *testing.T) {
	t.Parallel()

	expectedResult := &data.GenericAPIResponse{
		Data: &testStruct{
			Nonce: 0,
			Hash:  "aaaa",
		},
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
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
	)

	actualResult, err := epf.GetRatingsConfig()
	require.Nil(t, err)

	assert.Equal(t, expectedResult, actualResult)
}

func getPrivKey() crypto.PrivateKey {
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	sk, _ := keyGen.GeneratePair()

	return sk
}
