package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseReader(t *testing.T) {
	//t.Skip("this a manual tests that run only with a valid elasticseach database")

	url := "https://elastic-aws.elrond.com"
	user := "basic_auth_username"
	password := "basic_auth_password"
	reader, _ := NewDatabaseReader(url, user, password)

	addr := "erd1ewshdn9yv0wx38xgs5cdhvcq4dz0n7tdlgh8wfj9nxugwmyunnyqpkpzal"
	txs, err := reader.GetTransactionsByAddress(addr)
	fmt.Println(txs)
	assert.Nil(t, err)
}

func TestDatabaseReader_GetLatestBlockHeight(t *testing.T) {
	t.Skip("this a manual tests that run only with a valid elasticseach database")

	url := "https://elastic-aws.elrond.com"
	user := "basic_auth_username"
	password := "basic_auth_password"
	reader, _ := NewDatabaseReader(url, user, password)

	blockHeight, err := reader.GetLatestBlockHeight()
	fmt.Println(blockHeight)
	assert.Nil(t, err)
}

func TestDatabaseReader_GetBlock(t *testing.T) {
	t.Skip("this a manual tests that run only with a valid elasticseach database")

	url := "https://elastic-aws.elrond.com"
	user := "basic_auth_username"
	password := "basic_auth_password"
	reader, _ := NewDatabaseReader(url, user, password)

	block, err := reader.GetBlockByNonce(7561)
	fmt.Println(block)
	assert.Nil(t, err)
}
