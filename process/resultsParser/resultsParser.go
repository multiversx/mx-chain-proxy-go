package resultsParser

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"

	"github.com/multiversx/mx-chain-proxy-go/process/resultsParser/transactionDecoder"
)

var log = logger.GetOrCreate("api/gin")

// TODO: add full comment
// ResultOutcome -
type ResultOutcome struct {
	ReturnCode    ReturnCode      `json:"returnCode"`
	ReturnMessage string          `json:"returnMessage"`
	Values        []*bytes.Buffer `json:"values"`
}

// TODO: add full comment
// ParseResultOutcome -
func ParseResultOutcome(tx *transaction.ApiTransactionResult, pubKeyConverter core.PubkeyConverter) (*ResultOutcome, error) {
	metadata, err := transactionDecoder.GetTransactionMetadata(transactionDecoder.TransactionToDecode{
		Sender:   tx.Sender,
		Receiver: tx.Receiver,
		Data:     string(tx.Data),
		Value:    tx.Value,
	}, pubKeyConverter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction metadata: %w", err)
	}

	outcome := parseOutcomeOnSimpleMoveBalance(tx)
	if outcome != nil {
		log.Trace("result outcome on simple move balance")
		return outcome, nil
	}

	outcome = parseOutcomeOnInvalidTransaction(tx)
	if outcome != nil {
		log.Trace("result outcome on invalid transaction")
		return outcome, nil
	}

	outcome, err = parseOutcomeOnEasilyFoundResultWithReturnData(tx.SmartContractResults)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on easily found result with return data: %w", err)
	}
	if outcome != nil {
		log.Trace("result outcome on easily found result with return data")
		return outcome, nil
	}

	outcome, err = parseOutcomeOnSignalError(tx.Logs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on signal error: %w", err)
	}
	if outcome != nil {
		log.Trace("result outcome on signal error")
		return outcome, nil
	}

	outcome, err = parseOutcomeOnTooMuchGasWarning(tx.Logs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on too much gas warning: %w", err)
	}
	if outcome != nil {
		log.Trace("result outcome on too much gas warning")
		return outcome, nil
	}

	outcome, err = parseOutcomeOnWriteLogWhereFirstTopicEqualsAddress(tx.Logs, tx.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to parse outcome on write log where first topic equals address: %w", err)
	}
	if outcome != nil {
		log.Trace("on writelog with topics[0] == tx.sender")
		return outcome, nil
	}

	outcome, err = parseOutcomeWithFallbackHeuristics(tx, *metadata)
	if outcome != nil {
		log.Trace("result outcome on fallback heuristics")
		panic("whatever")
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
			ReturnCode:    None,
			ReturnMessage: None.String(),
		}
	}

	return nil
}

func parseOutcomeOnInvalidTransaction(tx *transaction.ApiTransactionResult) *ResultOutcome {
	if tx.Status == transaction.TxStatusInvalid {

		if tx.Receipt != nil && tx.Receipt.Data != "" {
			return &ResultOutcome{
				ReturnCode:    OutOfFunds,
				ReturnMessage: tx.Receipt.Data,
			}
		}
	}

	// If there's no receipt message, let other heuristics to handle the outcome (most probably, a log with "signalError" is emitted).

	return nil
}

func parseOutcomeOnEasilyFoundResultWithReturnData(results []*transaction.ApiSmartContractResult) (*ResultOutcome, error) {
	var r *transaction.ApiSmartContractResult
	for _, result := range results {
		if result.Nonce != 0 && result.Data != "" && result.Data[0] == '@' {
			r = result
			break
		}
	}

	if r == nil {
		return nil, nil
	}

	returnCode, returnDataParts, err := sliceDataFieldInParts(r.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to slice data field in parts: %w", err)
	}

	returnMessage := returnCode.String()
	if r.ReturnMessage != "" {
		returnMessage = r.ReturnMessage
	}

	return &ResultOutcome{
		ReturnCode:    *returnCode,
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
		ReturnCode:    *returnCode,
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
		ReturnCode:    *returnCode,
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
		ReturnCode:    *returnCode,
		ReturnMessage: returnCode.String(),
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeWithFallbackHeuristics(tx *transaction.ApiTransactionResult, metadata transactionDecoder.TransactionMetadata) (*ResultOutcome, error) {
	for _, resultItem := range tx.SmartContractResults {
		event, findErr := findSingleOrNoneEvent(resultItem.Logs, OnWriteLog, func(e *transaction.Events) *transaction.Events {
			addressIsSender := e.Address == tx.Sender
			firstTopicIsContract := false
			if e.Topics[0] != nil {
				decodeString, topicDecode := hex.DecodeString(string(e.Topics[0]))
				if topicDecode != nil {
					return nil
				}
				firstTopicIsContract = string(decodeString) == metadata.Receiver

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
			ReturnCode:    *returnCode,
			ReturnMessage: returnCode.String(),
			Values:        returnDataParts,
		}, nil
	}

	return nil, nil
}

func sliceDataFieldInParts(data string) (*ReturnCode, []*bytes.Buffer, error) {
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
	parts := stringToBuffers(data)
	returnCodePart := parts[startingIndex]
	returnDataParts := parts[startingIndex+1:]

	if returnCodePart.Len() == 0 {
		return nil, nil, errors.New("no return code")
	}

	returnCode := fromBuffer(*returnCodePart)
	return &returnCode, returnDataParts, nil
}

func stringToBuffers(joinedString string) []*bytes.Buffer {
	splits := strings.Split(joinedString, "@")
	b := make([]*bytes.Buffer, len(splits))
	for i, s := range splits {
		bufferString := bytes.NewBufferString(s)
		b[i] = bufferString
	}

	return b
}

func findSingleOrNoneEvent(
	logs *transaction.ApiLogs,
	identifier string,
	predicate func(e *transaction.Events) *transaction.Events,
) (*transaction.Events, error) {
	if logs == nil {
		return nil, nil
	}

	events := findEvents(logs.Events, identifier, predicate)

	if len(events) > 1 {
		return nil, errors.New("found more than one event")
	}

	if events == nil {
		return nil, nil
	}

	return events[0], nil
}

func findEvents(events []*transaction.Events, identifier string, predicate func(e *transaction.Events) *transaction.Events) []*transaction.Events {
	var matches []*transaction.Events
	for _, event := range events {
		if event.Identifier == identifier {
			if predicate != nil {
				e := predicate(event)

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

func findFirstOrNoneTopic(topics [][]byte, predicate func(topic []byte) []byte) []byte {
	for _, topic := range topics {
		t := predicate(topic)
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
