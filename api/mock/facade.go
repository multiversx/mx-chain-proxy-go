package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Facade is the mock implementation of a node's router handler
type Facade struct {
	GetAccountHandler                 func(address string) (*data.Account, error)
	SendTransactionHandler            func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler   func(txs []*data.Transaction) (uint64, error)
	SendUserFundsCalled               func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler             func(query *process.SCQuery) (*vmcommon.VMOutput, error)
	GetHeartbeatDataHandler           func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler        func() (map[string]*data.ValidatorApiResponse, error)
	SendTransactionCostRequestHandler func(tx *data.Transaction) (string, error)
	GetNodeStatusDataHandler          func(shardId string) (map[string]interface{}, error)
}

// ValidatorStatistics is the mock implementation of a handler's ValidatorStatistics method
func (f *Facade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return f.ValidatorStatisticsHandler()
}

// GetNodeStatusData --
func (f *Facade) GetNodeStatusData(shardId string) (map[string]interface{}, error) {
	return f.GetNodeStatusDataHandler(shardId)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SendMultipleTransactions is the mock implementation of a handler's SendMultipleTransactions method
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// SendTransactionCostRequest --
func (f *Facade) SendTransactionCostRequest(tx *data.Transaction) (string, error) {
	return f.SendTransactionCostRequestHandler(tx)
}

// SendUserFunds is the mock implementation of a handler's SendUserFunds method
func (f *Facade) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery is a mock implementation.
func (f *Facade) ExecuteSCQuery(query *process.SCQuery) (*vmcommon.VMOutput, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData is the mock implementation of a handler's GetHeartbeatData method
func (f *Facade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
