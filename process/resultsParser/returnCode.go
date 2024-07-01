package resultsParser

import (
	"bytes"
	"encoding/hex"
)

type ReturnCode string

const (
	None                   ReturnCode = ""
	Ok                     ReturnCode = "ok"
	FunctionNotFound       ReturnCode = "function not found"
	FunctionWrongSignature ReturnCode = "wrong signature for function"
	ContractNotFound       ReturnCode = "contract not found"
	UserError              ReturnCode = "user error"
	OutOfGas               ReturnCode = "out of gas"
	AccountCollision       ReturnCode = "account collision"
	OutOfFunds             ReturnCode = "out of funds"
	CallStackOverFlow      ReturnCode = "call stack overflow"
	ContractInvalid        ReturnCode = "contract invalid"
	ExecutionFailed        ReturnCode = "execution failed"
	Unknown                ReturnCode = "unknown"
)

func (rc ReturnCode) String() string {
	return string(rc)
}

func fromBuffer(bytes bytes.Buffer) ReturnCode {
	text := bytes.String()
	decodeString, _ := hex.DecodeString(text)

	return ReturnCode(decodeString)
}
