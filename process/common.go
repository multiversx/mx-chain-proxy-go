package process

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

func prepareQuerySc(dataForQuery []byte) (*process.SCQuery, error) {
	var query process.SCQuery

	var p []interface{}
	err := json.Unmarshal(dataForQuery, &p)
	if err != nil {
		return &query, err
	}

	params := make(map[string]string)
	v := reflect.ValueOf(p[0])
	for _, key := range v.MapKeys() {
		params[fmt.Sprint(key.Interface())] = fmt.Sprint(v.MapIndex(key).Interface())
	}

	argumentParser := vmcommon.NewAtArgumentParser()

	err = argumentParser.ParseData(params["data"])
	if err != nil {
		return &process.SCQuery{}, err
	}

	function, err := argumentParser.GetFunction()
	if err != nil {
		return &process.SCQuery{}, err
	}

	arguments, err := argumentParser.GetArguments()
	if err != nil {
		return &process.SCQuery{}, err
	}

	query.ScAddress, _ = hex.DecodeString(removeHexSuffix(params["to"]))
	query.FuncName = function
	query.Arguments = arguments

	return &query, nil
}

func prepareTx(dataForTx []byte) (*data.Transaction, error) {
	var tx data.Transaction

	var params []map[string]string
	err := json.Unmarshal(dataForTx, &params)
	if err != nil {
		return &data.Transaction{}, err
	}

	tx.Nonce = hex2int(params[0]["nonce"])
	tx.Value = fmt.Sprintf("%d", hex2int(params[0]["value"]))
	tx.Sender = removeHexSuffix(params[0]["from"])
	tx.Receiver = removeHexSuffix(params[0]["to"])
	tx.Data = []byte(removeHexSuffix(params[0]["data"]))

	tx.GasPrice, _ = strconv.ParseUint(params[0]["gasPrice"], 10, 64)
	tx.GasLimit, _ = strconv.ParseUint(params[0]["gasLimit"], 10, 64)

	tx.Signature = params[0]["signature"]

	return &tx, nil
}

func prepareSignedTx(dataForTx []byte) (*data.Transaction, error) {
	var tx data.Transaction

	var params []map[string]string
	err := json.Unmarshal(dataForTx, &params)
	if err != nil {
		return &data.Transaction{}, err
	}

	tx.Nonce = hex2int(params[0]["nonce"])
	tx.Value = fmt.Sprintf("%d", hex2int(params[0]["value"]))
	tx.Sender = removeHexSuffix(params[0]["sender"])
	tx.Receiver = removeHexSuffix(params[0]["receiver"])
	tx.Data = []byte(removeHexSuffix(params[0]["data"]))

	tx.GasPrice, _ = strconv.ParseUint(params[0]["gasPrice"], 10, 64)
	tx.GasLimit, _ = strconv.ParseUint(params[0]["gasLimit"], 10, 64)

	tx.Signature = params[0]["signature"]

	return &tx, nil
}

func removeHexSuffix(hexStr string) string {
	return strings.Replace(hexStr, "0x", "", -1)
}

func int2hex(value uint64) string {
	s := fmt.Sprintf("%0x", value)
	return "0x" + s
}

func hex2int(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := removeHexSuffix(hexStr)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return result
}
