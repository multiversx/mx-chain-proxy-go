package transactionDecoder

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/stretchr/testify/require"
)

var testPubkeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "erd")

func Test_GetTransactionMetadata(t *testing.T) {
	t.Parallel()

	t.Run("NFT Smart contract call", func(t *testing.T) {
		t.Parallel()

		tx := TransactionToDecode{
			Sender:   "erd18w6yj09l9jwlpj5cjqq9eccfgulkympv7d4rj6vq4u49j8fpwzwsvx7e85",
			Receiver: "erd18w6yj09l9jwlpj5cjqq9eccfgulkympv7d4rj6vq4u49j8fpwzwsvx7e85",
			Data:     "RVNEVE5GVFRyYW5zZmVyQDRjNGI0ZDQ1NTgyZDYxNjE2MjM5MzEzMEAyZmI0ZTlAZTQwZjE2OTk3MTY1NWU2YmIwNGNAMDAwMDAwMDAwMDAwMDAwMDA1MDBkZjNiZWJlMWFmYTEwYzQwOTI1ZTgzM2MxNGE0NjBlMTBhODQ5ZjUwYTQ2OEA3Mzc3NjE3MDVmNmM2YjZkNjU3ODVmNzQ2ZjVmNjU2NzZjNjRAMGIzNzdmMjYxYzNjNzE5MUA=",
			Value:    "0",
		}

		metadata, err := GetTransactionMetadata(tx, testPubkeyConverter)
		require.NoError(t, err)
		require.Equal(t, tx.Sender, metadata.Sender)
		require.Equal(t, "erd1qqqqqqqqqqqqqpgqmua7hcd05yxypyj7sv7pffrquy9gf86s535qxct34s", metadata.Receiver)
		value, _ := big.NewInt(0).SetString("1076977887712805212893260", 10)
		require.Equal(t, value, metadata.Value)
		require.Equal(t, "swap_lkmex_to_egld", metadata.FunctionName)
		require.Equal(t, []string{"0b377f261c3c7191", ""}, metadata.FunctionArgs)
		value1, _ := big.NewInt(0).SetString("1076977887712805212893260", 10)
		require.Equal(t, []TransactionMetadataTransfer{
			{
				Value: value1,
				Properties: TokenTransferProperties{
					Identifier: "LKMEX-aab910-2fb4e9",
					Collection: "LKMEX-aab910",
				},
			},
		}, metadata.Transfers)

		fmt.Println(metadata)
	})

	t.Run("ESDT Transfer", func(t *testing.T) {
		t.Parallel()

		tx := TransactionToDecode{
			Sender:   "erd1jvc6nyyl73q2yardw7dj8235h5zqaum4qyc8wlgs6aa26seysuvsrp48x2",
			Receiver: "erd1flqg2zf3knya94lcupscdwmrud029mes8a85r202rvwpzjyk5tjqxt8dxu",
			Data:     "RVNEVFRyYW5zZmVyQDUwNGM0MTU0NDEyZDM5NjI2MTM2NjMzM0AwMTJhMDVmMjAw",
			Value:    "0",
		}

		metadata, err := GetTransactionMetadata(tx, testPubkeyConverter)
		require.NoError(t, err)
		require.Equal(t, tx.Sender, metadata.Sender)
		require.Equal(t, tx.Receiver, metadata.Receiver)
		value, _ := big.NewInt(0).SetString("5000000000", 10)
		require.Equal(t, value, metadata.Value)
		require.Equal(t, []TransactionMetadataTransfer{{Value: value, Properties: TokenTransferProperties{Identifier: "PLATA-9ba6c3"}}}, metadata.Transfers)
	})

	t.Run("MultiESDTNFTTransfer fungible (with 00 nonce) + meta", func(t *testing.T) {
		t.Parallel()

		tx := TransactionToDecode{
			Sender:   "erd1lkrrrn3ws9sp854kdpzer9f77eglqpeet3e3k3uxvqxw9p3eq6xqxj43r9",
			Receiver: "erd1lkrrrn3ws9sp854kdpzer9f77eglqpeet3e3k3uxvqxw9p3eq6xqxj43r9",
			Data:     "TXVsdGlFU0RUTkZUVHJhbnNmZXJAMDAwMDAwMDAwMDAwMDAwMDA1MDBkZjNiZWJlMWFmYTEwYzQwOTI1ZTgzM2MxNGE0NjBlMTBhODQ5ZjUwYTQ2OEAwMkA0YzRiNGQ0NTU4MmQ2MTYxNjIzOTMxMzBAMmZlM2IwQDA5Yjk5YTZkYjMwMDI3ZTRmM2VjQDU1NTM0NDQzMmQzMzM1MzA2MzM0NjVAMDBAMDEyNjMwZTlhMjlmMmY5MzgxNDQ5MUA3MDYxNzk1ZjZkNjU3NDYxNWY2MTZlNjQ1ZjY2NzU2ZTY3Njk2MjZjNjVAMGVkZTY0MzExYjhkMDFiNUA=",
			Value:    "0",
		}

		metadata, err := GetTransactionMetadata(tx, testPubkeyConverter)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, tx.Sender, metadata.Sender)
		require.Equal(t, "erd1qqqqqqqqqqqqqpgqmua7hcd05yxypyj7sv7pffrquy9gf86s535qxct34s", metadata.Receiver)
		require.Equal(t, tx.Value, metadata.Value.String())
		require.Equal(t, "pay_meta_and_fungible", metadata.FunctionName)
		require.Equal(t, []string{"0ede64311b8d01b5", ""}, metadata.FunctionArgs)
		value1, _ := big.NewInt(0).SetString("45925073746530627023852", 10)
		value2, _ := big.NewInt(0).SetString("1389278024872597502641297", 10)
		require.Equal(t, []TransactionMetadataTransfer{
			{
				Value: value1,
				Properties: TokenTransferProperties{
					Identifier: "LKMEX-aab910-2fe3b0",
					Collection: "LKMEX-aab910",
				},
			},
			{
				Value: value2,
				Properties: TokenTransferProperties{
					Token: "USDC-350c4e",
				},
			},
		}, metadata.Transfers)
	})

	t.Run("MultiESDTNFTTransfer fungibles (00 and missing nonce)", func(t *testing.T) {
		t.Parallel()

		tx := TransactionToDecode{
			Sender:   "erd1lkrrrn3ws9sp854kdpzer9f77eglqpeet3e3k3uxvqxw9p3eq6xqxj43r9",
			Receiver: "erd1lkrrrn3ws9sp854kdpzer9f77eglqpeet3e3k3uxvqxw9p3eq6xqxj43r9",
			Data:     "TXVsdGlFU0RUTkZUVHJhbnNmZXJAMDAwMDAwMDAwMDAwMDAwMDA1MDBkZjNiZWJlMWFmYTEwYzQwOTI1ZTgzM2MxNGE0NjBlMTBhODQ5ZjUwYTQ2OEAwMkA1MjQ5NDQ0NTJkMzAzNTYyMzE2MjYyQDAwQDA5Yjk5YTZkYjMwMDI3ZTRmM2VjQDU1NTM0NDQzMmQzMzM1MzA2MzM0NjVAQDAxMjYzMGU5YTI5ZjJmOTM4MTQ0OTE=",
			Value:    "0",
		}

		metadata, err := GetTransactionMetadata(tx, testPubkeyConverter)
		require.NoError(t, err)
		require.Equal(t, tx.Sender, metadata.Sender)
		require.Equal(t, "erd1qqqqqqqqqqqqqpgqmua7hcd05yxypyj7sv7pffrquy9gf86s535qxct34s", metadata.Receiver)
		require.Equal(t, tx.Value, metadata.Value.String())
		value1, _ := big.NewInt(0).SetString("45925073746530627023852", 10)
		value2, _ := big.NewInt(0).SetString("1389278024872597502641297", 10)
		require.Equal(t, []TransactionMetadataTransfer{
			{
				Value: value1,
				Properties: TokenTransferProperties{
					Token: "RIDE-05b1bb",
				},
			},
			{
				Value: value2,
				Properties: TokenTransferProperties{
					Token: "USDC-350c4e",
				},
			},
		}, metadata.Transfers)
	})

	t.Run("Relayed transaction v1", func(t *testing.T) {
		t.Parallel()

		tx := TransactionToDecode{
			Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
			Receiver: "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx",
			Data:     "cmVsYXllZFR4QDdiMjI2ZTZmNmU2MzY1MjIzYTMxMzkzODJjMjI3MzY1NmU2NDY1NzIyMjNhMjI2NzQ1NmU1NzRmNjU1NzZkNmQ0MTMwNjMzMDZhNmI3MTc2NGQzNTQyNDE3MDdhNjE2NDRiNDY1NzRlNTM0ZjY5NDE3NjQzNTc1MTYzNzc2ZDQ3NTA2NzNkMjIyYzIyNzI2NTYzNjU2OTc2NjU3MjIyM2EyMjQxNDE0MTQxNDE0MTQxNDE0MTQxNDE0NjQxNDIzNDc1NTk1MjcxNjMzNDY1NDQ0OTM0Nzk2NzM4N2E0ODc3NjI0NDMwNWE2ODZiNTg0MjM1NzAzMTc3M2QyMjJjMjI3NjYxNmM3NTY1MjIzYTMwMmMyMjY3NjE3MzUwNzI2OTYzNjUyMjNhMzEzMDMwMzAzMDMwMzAzMDMwMzAyYzIyNjc2MTczNGM2OTZkNjk3NDIyM2EzNjMwMzAzMDMwMzAzMDMwMmMyMjY0NjE3NDYxMjIzYTIyNTk1NzUyNmIyMjJjMjI3MzY5Njc2ZTYxNzQ3NTcyNjUyMjNhMjI0ZTMwNzIzMTcwNmYzNzZiNzY0ZjU0NGI0OTQ3NDcyZjc1NmI2NzcyMzg1YTYyNTc2NDU4NjczMTY2NTEzMDc2NmQ3NTYyMzU3OTM0NGY3MzUzNDE3MTM0N2EyZjU5Mzc2YzQ2NTI3OTU3NzM2NzM0NGUyYjZmNGE2OTQ5NDk1Nzc3N2E2YjZkNmM2YTQ5NDE3MjZkNjkzMTY5NTg0ODU0NzkzNDRiNjc0MTQxM2QzZDIyMmMyMjYzNjg2MTY5NmU0OTQ0MjIzYTIyNTY0MTNkM2QyMjJjMjI3NjY1NzI3MzY5NmY2ZTIyM2EzMTdk",
			Value:    "0",
		}

		metadata, err := GetTransactionMetadata(tx, testPubkeyConverter)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx", metadata.Sender)
		require.Equal(t, "erd1qqqqqqqqqqqqqpgqrchxzx5uu8sv3ceg8nx8cxc0gesezure5awqn46gtd", metadata.Receiver)
		require.Equal(t, tx.Value, metadata.Value.String())
		require.Equal(t, "add", metadata.FunctionName)
		require.Empty(t, metadata.FunctionArgs)
	})

	t.Run("Relayed transaction v2", func(t *testing.T) {
		t.Parallel()

		tx := TransactionToDecode{
			Sender:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
			Receiver: "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx",
			Data:     "cmVsYXllZFR4VjJAMDAwMDAwMDAwMDAwMDAwMDA1MDAxZTJlNjExYTljZTFlMGM4ZTMyODNjY2M3YzFiMGY0NjYxOTE3MDc5YTc1Y0AwZkA2MTY0NjRAOWFiZDEzZjRmNTNmM2YyMzU5Nzc0NGQ2NWZjNWQzNTFiYjY3NzNlMDVhOTU0YjQxOWMwOGQxODU5M2QxYzY5MjYyNzlhNGQxNjE0NGQzZjg2NmE1NDg3ODAzMTQyZmNmZjBlYWI2YWQ1ODgyMDk5NjlhY2I3YWJlZDIxMDIwMGI=",
			Value:    "0",
		}

		metadata, err := GetTransactionMetadata(tx, testPubkeyConverter)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx", metadata.Sender)
		require.Equal(t, "erd1qqqqqqqqqqqqqpgqrchxzx5uu8sv3ceg8nx8cxc0gesezure5awqn46gtd", metadata.Receiver)
		require.Equal(t, tx.Value, metadata.Value.String())
		require.Equal(t, "add", metadata.FunctionName)
		require.Empty(t, metadata.FunctionArgs)
	})
}
