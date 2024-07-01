package resultsParser

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/data/transaction"

	"github.com/stretchr/testify/require"
)

var testPubkeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "erd")

func TestResultsParser_ParseUntypedOutcome(t *testing.T) {
	t.Parallel()

	t.Run("should parse contract outcome, on easily found result with return data", func(t *testing.T) {
		t.Parallel()

		transactionResult := &transaction.ApiTransactionResult{
			SmartContractResults: []*transaction.ApiSmartContractResult{
				{
					Nonce:         42,
					Data:          "@6f6b@03",
					ReturnMessage: "foobar",
				},
			},
		}

		bundle, _ := ParseResultOutcome(transactionResult, testPubkeyConverter)
		require.Equal(t, Ok, bundle.ReturnCode)
		require.Equal(t, bundle.Values, []*bytes.Buffer{bytes.NewBuffer([]byte("03"))})
	})

	t.Run("should parse contract outcome, on signal error", func(t *testing.T) {
		t.Parallel()

		transactionResult := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Address: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
				Events: []*transaction.Events{
					{
						Identifier: "signalError",
						Topics: [][]byte{
							[]byte("something happened"),
						},
						Data: []byte("@75736572206572726f72@07"),
					},
				},
			},
		}

		outcome, _ := ParseResultOutcome(transactionResult, testPubkeyConverter)
		require.Equal(t, UserError, outcome.ReturnCode)
		require.Equal(t, outcome.Values, []*bytes.Buffer{bytes.NewBuffer([]byte("07"))})
	})

	t.Run("should parse contract outcome, on too much gas warning", func(t *testing.T) {
		t.Parallel()

		transactionResult := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Address: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
				Events: []*transaction.Events{
					{
						Identifier: "writeLog",
						Topics: [][]byte{
							[]byte("QHRvbyBtdWNoIGdhcyBwcm92aWRlZCBmb3IgcHJvY2Vzc2luZzogZ2FzIHByb3ZpZGVkID0gNTk2Mzg0NTAwLCBnYXMgdXNlZCA9IDczMzAxMA=="),
						},
						Data: []byte("@6f6b"),
					},
				},
			},
		}

		outcome, _ := ParseResultOutcome(transactionResult, testPubkeyConverter)
		require.Equal(t, Ok, outcome.ReturnCode)
		require.Empty(t, outcome.Values)
	})

	t.Run("should parse contract outcome, on write log where first topic equals address", func(t *testing.T) {
		t.Parallel()

		transactionResult := &transaction.ApiTransactionResult{
			Sender: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
			Logs: &transaction.ApiLogs{
				Events: []*transaction.Events{
					{
						Identifier: "writeLog",
						Topics: [][]byte{
							[]byte("ZXJkMXF5dTV3dGhsZHpyOHd4NWM5dWNnOGtqYWdnMGpmczUzczhucjN6cHozaHlwZWZzZGQ4c3N5Y3I2dGg="),
						},
						Data: []byte("@6f6b="),
					},
				},
			},
		}

		outcome, _ := ParseResultOutcome(transactionResult, testPubkeyConverter)
		require.Equal(t, Ok, outcome.ReturnCode)
		require.Empty(t, outcome.Values)
	})
}

func TestResultsParser_RealWorld(t *testing.T) {
	//TODO: do commit the skip
	t.Skip()

	filePath := "transactions.json"

	txs, err := readJSONFromFile(filePath)
	if err != nil {
		fmt.Printf("Error reading from file: %v\n", err)
		return
	}

	var nilOutcomes []*transaction.ApiTransactionResult
	for i, tx := range txs {
		outcome, err := ParseResultOutcome(tx, testPubkeyConverter)
		if err != nil {
			panic(fmt.Errorf("error parsing results: %v %d\n", err, i))
		}

		if outcome == nil {
			nilOutcomes = append(nilOutcomes, tx)
		}

	}

	for _, tx := range nilOutcomes {
		//o_, _ := ParseResultOutcome(tx)

		for _, e := range tx.Logs.Events {

			for _, tt := range e.Topics {
				decodeString, err := base64.StdEncoding.DecodeString(string(tt))
				if err != nil {
					fmt.Printf("Error decoding base64 string: %v\n", err)
					continue
				}

				fmt.Println(decodeString)
			}
		}

		//fmt.Println(o)
	}
}

func readJSONFromFile(filePath string) ([]*transaction.ApiTransactionResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var txs []*transaction.ApiTransactionResult
	if err := json.Unmarshal(byteValue, &txs); err != nil {
		return nil, err
	}

	return txs, nil
}

func TestA(t *testing.T) {
	t.Skip()
	c := http.Client{}

	req, _ := http.NewRequest("GET", "https://gateway.multiversx.com/transaction/393db73fde175727009c50629220e4be6e36618037a9e163757eab34934876be?withResults=true", nil)

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	response := struct {
		Data struct {
			Transaction *transaction.ApiTransactionResult `json:"transaction"`
		} `json:"data"`
	}{}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)

}
