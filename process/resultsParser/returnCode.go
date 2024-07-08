package resultsParser

import (
	"bytes"
	"encoding/hex"
)

// ReturnCode is an enum that describes some common smart contract results.
type ReturnCode string

const (
	None       ReturnCode = ""
	Ok         ReturnCode = "ok"
	UserError  ReturnCode = "user error"
	OutOfFunds ReturnCode = "out of funds"
)

func (rc ReturnCode) String() string {
	return string(rc)
}

func fromBuffer(bytes bytes.Buffer) ReturnCode {
	text := bytes.String()
	decodeString, _ := hex.DecodeString(text)

	return ReturnCode(decodeString)
}
