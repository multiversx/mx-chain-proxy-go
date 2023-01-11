package database

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/stretchr/testify/require"
)

func TestDatabaseReader(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.multiversx.com"
	user := ""
	password := ""
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	addr := "erd1ewshdn9yv0wx38xgs5cdhvcq4dz0n7tdlgh8wfj9nxugwmyunnyqpkpzal"
	txs, err := reader.GetTransactionsByAddress(addr)
	fmt.Println(txs)
	require.Nil(t, err)
}

func TestDatabaseReader_GetBlockByShardIDAndNonce(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.multiversx.com"
	user := ""
	password := ""
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	block, err := reader.GetAtlasBlockByShardIDAndNonce(core.MetachainShardId, 7720)
	fmt.Println(block)
	require.Nil(t, err)
}
