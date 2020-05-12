package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseReader(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.elrond.com"
	user := "basic_auth_username"
	password := "basic_auth_password"
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	addr := "erd1ewshdn9yv0wx38xgs5cdhvcq4dz0n7tdlgh8wfj9nxugwmyunnyqpkpzal"
	txs, err := reader.GetTransactionsByAddress(addr)
	fmt.Println(txs)
	require.Nil(t, err)
}

func TestDatabaseReader_GetLatestBlockHeight(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.elrond.com"
	user := "basic_auth_username"
	password := "basic_auth_password"
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	blockHeight, err := reader.GetLatestBlockHeight()
	fmt.Println(blockHeight)
	require.Nil(t, err)
}

func TestDatabaseReader_GetBlock(t *testing.T) {
	t.Skip("this test queries Elastic Search")

	url := "https://elastic-aws.elrond.com"
	user := "basic_auth_username"
	password := "basic_auth_password"
	reader, err := NewElasticSearchConnector(url, user, password)
	require.Nil(t, err)

	block, err := reader.GetBlockByNonce(7561)
	fmt.Println(block)
	require.Nil(t, err)
}
