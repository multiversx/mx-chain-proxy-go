package resultsParser

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"
	vm "github.com/multiversx/mx-chain-vm-common-go"
)

var log = logger.GetOrCreate("resultsParser")

// ResultOutcome encapsulates data contained within the smart contact results.
type ResultOutcome struct {
	ReturnCode    *vm.ReturnCode `json:"returnCode"`
	ReturnMessage string         `json:"returnMessage"`
	Values        [][]byte       `json:"values"`
}

// ParseResultOutcome will try to translate the smart contract results into a ResultOutcome object.
func ParseResultOutcome(tx *transaction.ApiTransactionResult, pubKeyConverter core.PubkeyConverter) (*ResultOutcome, error) {
	outcome := parseOutcomeOnSimpleMoveBalance(tx)
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on simple move balance", tx.Hash)
		return outcome, nil
	}

	outcome = parseOutcomeOnInvalidTransaction(tx)
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on invalid transaction", tx.Hash)
		return outcome, nil
	}

	outcome, err := parseOutcomeOnEasilyFoundResultWithReturnData(tx.SmartContractResults)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on easily found result with return data: %w", err)
	}
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on easily found result with return data", tx.Hash)
		return outcome, nil
	}

	outcome, err = parseOutcomeOnSignalError(tx.Logs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on signal error: %w", err)
	}
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on signal error", tx.Hash)
		return outcome, nil
	}

	outcome, err = parseOutcomeOnTooMuchGasWarning(tx.Logs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on too much gas warning: %w", err)
	}
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on too much gas warning", tx.Hash)
		return outcome, nil
	}

	outcome, err = parseOutcomeOnWriteLogWhereFirstTopicEqualsAddress(tx.Logs, tx.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on write log where first topic equals address: %w", err)
	}
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on write log where first topic equals address", tx.Hash)
		return outcome, nil
	}

	outcome, err = parseOutcomeWithFallbackHeuristics(tx, tx.Receivers[0])
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on fallback heuristics", tx.Hash)
		return outcome, nil
	}
	return nil, nil
}

func parseOutcomeOnSimpleMoveBalance(tx *transaction.ApiTransactionResult) *ResultOutcome {
	noResults := len(tx.SmartContractResults) == 0
	var noLogs bool
	if tx.Logs != nil {
		noLogs = len(tx.Logs.Events) == 0
	} else {
		noLogs = true
	}

	if noResults && noLogs {
		return &ResultOutcome{
			ReturnCode:    nil,
			ReturnMessage: "",
		}
	}

	return nil
}

func parseOutcomeOnInvalidTransaction(tx *transaction.ApiTransactionResult) *ResultOutcome {
	if tx.Status == transaction.TxStatusInvalid {
		if tx.Receipt != nil && tx.Receipt.Data != "" {

			returnCode := vm.OutOfGas
			return &ResultOutcome{
				ReturnCode:    &returnCode,
				ReturnMessage: tx.Receipt.Data,
			}
		}
	}

	// If there's no receipt message, let other heuristics to handle the outcome (most probably, a log with "signalError" is emitted).

	return nil
}

func parseOutcomeOnEasilyFoundResultWithReturnData(scResults []*transaction.ApiSmartContractResult) (*ResultOutcome, error) {
	var scr *transaction.ApiSmartContractResult
	for _, scResult := range scResults {
		if scResult.Nonce != 0 && scResult.Data != "" && scResult.Data[0] == '@' {
			scr = scResult
			break
		}
	}

	if scr == nil {
		return nil, nil
	}

	returnCode, returnDataParts, err := sliceDataFieldInParts(scr.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to slice data field in parts: %w", err)
	}

	returnMessage := returnCode.String()
	if scr.ReturnMessage != "" {
		returnMessage = scr.ReturnMessage
	}

	return &ResultOutcome{
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeOnSignalError(logs *transaction.ApiLogs) (*ResultOutcome, error) {
	event, err := findSingleOrNoneEvent(logs, OnSignalError, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find single or none event: %w", err)
	}

	if event == nil {
		return nil, nil
	}

	returnCode, returnDataParts, err := sliceDataFieldInParts(string(event.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to slice data field in parts: %w", err)
	}
	lastTopic := getLastTopic(event.Topics)

	returnMessage := returnCode.String()
	if lastTopic != nil {
		returnMessage = string(lastTopic)
	}

	return &ResultOutcome{
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
		Values:        returnDataParts,
	}, nil

}

func parseOutcomeOnTooMuchGasWarning(logs *transaction.ApiLogs) (*ResultOutcome, error) {
	event, err := findSingleOrNoneEvent(logs, OnWriteLog, func(e *transaction.Events) *transaction.Events {
		t := findFirstOrNoneTopic(e.Topics, func(topic []byte) []byte {

			if strings.HasPrefix(string(topic), TooMuchGas) {
				return topic
			}
			return nil
		})

		if t != nil {
			return e
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find single or none event: %w", err)
	}

	if event == nil {
		return nil, nil
	}

	returnCode, returnDataParts, err := sliceDataFieldInParts(string(event.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to slice data field in parts: %w", err)
	}

	lastTopic := getLastTopic(event.Topics)

	returnMessage := returnCode.String()
	if lastTopic != nil {
		returnMessage = string(lastTopic)
	}

	return &ResultOutcome{
		ReturnCode:    returnCode,
		ReturnMessage: returnMessage,
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeOnWriteLogWhereFirstTopicEqualsAddress(logs *transaction.ApiLogs, address string) (*ResultOutcome, error) {
	base64Address := base64.StdEncoding.EncodeToString([]byte(address))

	event, err := findSingleOrNoneEvent(logs, OnWriteLog, func(e *transaction.Events) *transaction.Events {
		t := findFirstOrNoneTopic(e.Topics, func(topic []byte) []byte {
			if string(topic) == base64Address {
				return topic
			}
			return nil
		})

		if t != nil {
			return e
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find single or none event: %w", err)
	}

	if event == nil {
		return nil, nil
	}

	returnCode, returnDataParts, err := sliceDataFieldInParts(string(event.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to slice data field in parts: %w", err)
	}

	return &ResultOutcome{
		ReturnCode:    returnCode,
		ReturnMessage: returnCode.String(),
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeWithFallbackHeuristics(tx *transaction.ApiTransactionResult, receiver string) (*ResultOutcome, error) {
	for _, resultItem := range tx.SmartContractResults {
		event, findErr := findSingleOrNoneEvent(resultItem.Logs, OnWriteLog, func(e *transaction.Events) *transaction.Events {
			addressIsSender := e.Address == tx.Sender
			firstTopicIsContract := false
			if e.Topics[0] != nil {
				decodeString, topicDecode := hex.DecodeString(string(e.Topics[0]))
				if topicDecode != nil {
					return nil
				}
				firstTopicIsContract = string(decodeString) == receiver

				if firstTopicIsContract && addressIsSender {
					return e
				}
			}

			return nil
		})
		if findErr != nil {
			return nil, findErr
		}

		if event == nil {
			return nil, nil
		}

		returnCode, returnDataParts, sliceErr := sliceDataFieldInParts(string(event.Data))
		if sliceErr != nil {
			return nil, sliceErr
		}

		return &ResultOutcome{
			ReturnCode:    returnCode,
			ReturnMessage: returnCode.String(),
			Values:        returnDataParts,
		}, nil
	}

	return nil, nil
}

func sliceDataFieldInParts(data string) (*vm.ReturnCode, [][]byte, error) {
	if data == "" {
		return nil, nil, ErrEmptyDataField
	}

	// By default, skip the first part, which is usually empty (e.g. "[empty]@6f6b")
	startingIndex := 1

	// Before trying to parse the hex strings, cut the unwanted parts of the data field, in case of token transfers:
	if strings.HasPrefix(data, "ESDTTransfer") {
		// Skip "ESDTTransfer" (1), token identifier (2), amount (3)
		startingIndex = 3
	} else {
		// TODO: Upon gathering more transaction samples, fix for other kinds of transfers, as well (future PR, as needed).
	}

	// TODO: make this a function that returns a slice of bytes
	parts := stringToBytes(data)
	if len(parts) <= startingIndex {
		return nil, nil, ErrCannotProcessDataField
	}

	returnCodePart := parts[startingIndex]
	returnDataParts := parts[startingIndex+1:]

	if len(returnCodePart) == 0 {
		return nil, nil, ErrNoReturnCode
	}

	returnCode := parseReturnCodeFromHex(returnCodePart)
	if returnCode == 1 {
		return nil, nil, ErrNoReturnCode
	}
	return &returnCode, returnDataParts, nil
}

func stringToBytes(joinedString string) [][]byte {
	splits := strings.Split(joinedString, "@")
	b := make([][]byte, len(splits))
	for i, s := range splits {
		b[i] = []byte(s)
	}

	return b
}

func findSingleOrNoneEvent(
	logs *transaction.ApiLogs,
	identifier string,
	filter func(e *transaction.Events) *transaction.Events,
) (*transaction.Events, error) {
	if logs == nil {
		return nil, nil
	}

	events := findEvents(logs.Events, identifier, filter)

	if len(events) > 1 {
		return nil, ErrFoundMoreThanOneEvent
	}

	if events == nil {
		return nil, nil
	}

	return events[0], nil
}

func findEvents(events []*transaction.Events, identifier string, filter func(e *transaction.Events) *transaction.Events) []*transaction.Events {
	var matches []*transaction.Events
	for _, event := range events {
		if event.Identifier == identifier {
			if filter != nil {
				e := filter(event)

				if e != nil {
					matches = append(matches, e)
				}
				continue
			}
			matches = append(matches, event)
		}
	}

	return matches
}

func findFirstOrNoneTopic(topics [][]byte, filter func(topic []byte) []byte) []byte {
	for _, topic := range topics {
		t := filter(topic)
		if t != nil {
			return t
		}
	}

	return nil
}

func getLastTopic(topics [][]byte) []byte {
	if len(topics) == 0 {
		return nil
	}

	return topics[len(topics)-1]
}

func parseReturnCodeFromHex(hexEncodedBytes []byte) vm.ReturnCode {
	s, err := hex.DecodeString(string(hexEncodedBytes))
	if err != nil {
		return -1
	}

	switch string(s) {
	case "ok":
		return vm.Ok
	case "function not found":
		return vm.FunctionNotFound
	case "wrong signature for function":
		return vm.FunctionWrongSignature
	case "contract not found":
		return vm.ContractNotFound
	case "user error":
		return vm.UserError
	case "out of gas":
		return vm.OutOfGas
	case "account collision":
		return vm.AccountCollision
	case "out of funds":
		return vm.OutOfFunds
	case "call stack overflow":
		return vm.CallStackOverFlow
	case "contract invalid":
		return vm.ContractInvalid
	case "execution failed":
		return vm.ExecutionFailed

	default:
		return -1
	}
}
