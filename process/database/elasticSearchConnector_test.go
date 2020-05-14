package database

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/stretchr/testify/require"
)

func TestDatabaseReader(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.elrond.com"
	user := ""
	password := ""
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	addr := "erd1ewshdn9yv0wx38xgs5cdhvcq4dz0n7tdlgh8wfj9nxugwmyunnyqpkpzal"
	txs, err := reader.GetTransactionsByAddress(addr)
	fmt.Println(txs)
	require.Nil(t, err)
}

func TestDatabaseReader_GetBlock(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.elrond.com"
	user := ""
	password := ""
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	block, err := reader.GetBlockByShardIDAndNonce(core.MetachainShardId, 7720)
	fmt.Println(block)
	require.Nil(t, err)
}
