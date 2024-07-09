package resultsParser

import (
	"github.com/multiversx/mx-chain-vm-common-go/parsers/dataField"
)

type OperationalDataFieldParser interface {
	Parse(dataField []byte, sender, receiver []byte, numOfShards uint32) *datafield.ResponseParseData
}
