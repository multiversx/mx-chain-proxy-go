package resultsParser

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("resultsParser")

const (
	tooMuchGas  = "@too much gas"
	atSeparator = "@"
)

// ResultOutcome encapsulates data contained within the smart contact results.
type ResultOutcome struct {
	ReturnCode    string   `json:"returnCode"`
	ReturnMessage string   `json:"returnMessage"`
	Values        [][]byte `json:"values"`
}

// ParseResultOutcome will try to translate the smart contract results or logs into a ResultOutcome object.
func ParseResultOutcome(tx *transaction.ApiTransactionResult) (*ResultOutcome, error) {
	outcome := parseOutcomeOnSimpleMoveBalance(tx)
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on simple move balance", tx.Hash)
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

	outcome, err = parseOutcomeWithFallbackHeuristics(tx)
	if outcome != nil {
		log.Trace("txHash [%s] result outcome on fallback heuristics", tx.Hash)
		return outcome, nil
	}
	return nil, nil
}

func parseOutcomeOnSimpleMoveBalance(tx *transaction.ApiTransactionResult) *ResultOutcome {
	noResults := len(tx.SmartContractResults) == 0
	noLogs := tx.Logs == nil || len(tx.Logs.Events) == 0

	if noResults && noLogs {
		return &ResultOutcome{
			ReturnCode:    "",
			ReturnMessage: "",
		}
	}

	return nil
}

func parseOutcomeOnEasilyFoundResultWithReturnData(scResults []*transaction.ApiSmartContractResult) (*ResultOutcome, error) {
	var scr *transaction.ApiSmartContractResult
	for _, scResult := range scResults {
		//TODO: check whether scResult.Nonce should be the tx.Nonce + 1
		if scResult.Nonce != 0 && strings.HasPrefix(scResult.Data, atSeparator) {
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

	returnMessage := string(returnCode)
	if scr.ReturnMessage != "" {
		returnMessage = scr.ReturnMessage
	}

	return &ResultOutcome{
		ReturnCode:    string(returnCode),
		ReturnMessage: returnMessage,
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeOnTooMuchGasWarning(logs *transaction.ApiLogs) (*ResultOutcome, error) {
	event, err := findSingleOrNoneEvent(logs, core.WriteLogIdentifier, func(e *transaction.Events) *transaction.Events {
		t := findFirstOrNoneTopic(e.Topics, func(topic []byte) []byte {

			if strings.HasPrefix(string(topic), tooMuchGas) {
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

	returnMessage := string(returnCode)
	if lastTopic != nil {
		returnMessage = string(lastTopic)
	}

	return &ResultOutcome{
		ReturnCode:    string(returnCode),
		ReturnMessage: returnMessage,
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeOnWriteLogWhereFirstTopicEqualsAddress(logs *transaction.ApiLogs, address string) (*ResultOutcome, error) {
	base64Address := base64.StdEncoding.EncodeToString([]byte(address))

	event, err := findSingleOrNoneEvent(logs, core.WriteLogIdentifier, func(e *transaction.Events) *transaction.Events {
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
		ReturnCode:    string(returnCode),
		ReturnMessage: string(returnCode),
		Values:        returnDataParts,
	}, nil
}

func parseOutcomeWithFallbackHeuristics(tx *transaction.ApiTransactionResult) (*ResultOutcome, error) {
	if len(tx.Receivers) == 0 {
		return nil, nil
	}

	for _, resultItem := range tx.SmartContractResults {
		event, findErr := findSingleOrNoneEvent(resultItem.Logs, core.WriteLogIdentifier, func(e *transaction.Events) *transaction.Events {
			addressIsSender := e.Address == tx.Sender
			firstTopicIsContract := false
			if e.Topics[0] != nil {
				decodeString, topicDecode := hex.DecodeString(string(e.Topics[0]))
				if topicDecode != nil {
					return nil
				}
				firstTopicIsContract = string(decodeString) == tx.Receivers[0]

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
			ReturnCode:    string(returnCode),
			ReturnMessage: string(returnCode),
			Values:        returnDataParts,
		}, nil
	}

	return nil, nil
}

func sliceDataFieldInParts(data string) ([]byte, [][]byte, error) {
	if data == "" {
		return nil, nil, ErrEmptyDataField
	}

	parts := stringToBytes(data)

	if len(parts) == 0 || string(parts[0]) != "" {
		return nil, nil, ErrCannotProcessDataField
	}

	startingIndex := 1

	returnCodePart := parts[startingIndex]
	returnDataParts := parts[startingIndex+1:]

	if len(returnCodePart) == 0 {
		return nil, nil, ErrNoReturnCode
	}

	returnCode, err := hex.DecodeString(string(returnCodePart))
	if err != nil {
		return nil, nil, ErrCannotProcessDataField
	}
	return returnCode, returnDataParts, nil
}

func stringToBytes(joinedString string) [][]byte {
	splits := strings.Split(joinedString, atSeparator)
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
