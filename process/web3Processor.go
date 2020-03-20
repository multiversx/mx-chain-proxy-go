package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"strconv"
)

// Web3Processor --
type ProcessorWeb3 struct {
	nodeStatusProc  *NodeStatusProcessor
	transactionProc *TransactionProcessor
	accountProc     *AccountProcessor

	methods map[string]func(r data.RequestBodyWeb3) (data.ResponseWeb3, error)
}

// NewWeb3Processor will create a new web3 processor object
func NewWeb3Processor(nodeStatusProc *NodeStatusProcessor, txProc *TransactionProcessor, accountProc *AccountProcessor) (*ProcessorWeb3, error) {
	procWeb3 := &ProcessorWeb3{
		nodeStatusProc:  nodeStatusProc,
		transactionProc: txProc,
		accountProc:     accountProc,
	}

	procWeb3.initMethodsMap()

	return procWeb3, nil
}

func (pw *ProcessorWeb3) initMethodsMap() {
	methods := map[string]func(r data.RequestBodyWeb3) (data.ResponseWeb3, error){
		"eth_gasPrice":            pw.gasPrice,
		"eth_chainId":             pw.chainId,
		"eth_blockNumber":         pw.blockNumber,
		"eth_estimateGas":         pw.estimateGas,
		"eth_getTransactionCount": pw.transactionCount,
		"eth_sendRawTransaction":  pw.sendRawTransaction,
	}

	pw.methods = methods
}

// PrepareDataForRequest will request data from observer at prepare data in a proper format
func (pw *ProcessorWeb3) PrepareDataForRequest(requestBody data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	method, ok := pw.methods[requestBody.FuncName]
	if !ok {
		return data.ResponseWeb3{}, errors.New("invalid function")
	}

	return method(requestBody)
}

func (pw *ProcessorWeb3) sendRawTransaction(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	tx, err := prepareSignedTx(r.Params)
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	_, hash, err := pw.transactionProc.SendTransaction(tx)
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	return data.ResponseWeb3{
		JsonRpc: r.JsonRpc,
		Id:      r.Id,
		Result:  "0x" + hash,
	}, nil
}

func (pw *ProcessorWeb3) estimateGas(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	tx, err := prepareTx(r.Params)
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	cost, err := pw.transactionProc.SendTransactionCostRequest(tx)
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	costInt, _ := strconv.Atoi(cost)

	return data.ResponseWeb3{
		JsonRpc: r.JsonRpc,
		Id:      r.Id,
		Result:  int2hex(uint64(costInt)),
	}, nil
}

func (pw *ProcessorWeb3) gasPrice(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	nodeStatus, err := pw.nodeStatusProc.GetNodeStatusData("0")
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	nodeStatusMap := nodeStatus["details"].(map[string]interface{})

	gasPrice := nodeStatusMap[core.MetricMinGasPrice].(float64)

	return data.ResponseWeb3{
		JsonRpc: r.JsonRpc,
		Id:      r.Id,
		Result:  fmt.Sprintf("%d", uint64(gasPrice)),
	}, nil
}

func (pw *ProcessorWeb3) chainId(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	nodeStatus, err := pw.nodeStatusProc.GetNodeStatusData("0")
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	// TODO chain id maybe should be a number
	nodeStatusMap := nodeStatus["details"].(map[string]interface{})
	_ = nodeStatusMap[core.MetricChainId].(string)
	chainIdHex := "0x100"

	return data.ResponseWeb3{
		JsonRpc: r.JsonRpc,
		Id:      r.Id,
		Result:  chainIdHex,
	}, nil
}

func (pw *ProcessorWeb3) blockNumber(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	nodeStatus, err := pw.nodeStatusProc.GetNodeStatusData("0")
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	nodeStatusMap := nodeStatus["details"].(map[string]interface{})

	highestBlockNumber := nodeStatusMap[core.MetricProbableHighestNonce].(float64)

	return data.ResponseWeb3{
		JsonRpc: r.JsonRpc,
		Id:      r.Id,
		Result:  uint64(highestBlockNumber),
	}, nil
}

func (pw *ProcessorWeb3) transactionCount(r data.RequestBodyWeb3) (data.ResponseWeb3, error) {
	var params []string
	err := json.Unmarshal(r.Params, &params)
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	address := removeHexSuffix(params[0])

	account, err := pw.accountProc.GetAccount(address)
	if err != nil {
		return data.ResponseWeb3{}, err
	}

	return data.ResponseWeb3{
		JsonRpc: r.JsonRpc,
		Id:      r.Id,
		Result:  int2hex(account.Nonce),
	}, nil
}
