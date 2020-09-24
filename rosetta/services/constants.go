package services

const (
	NumBlocksToGet = uint64(200)

	// GenesisBlockHash is const that will keep genesis block hash in hex format
	GenesisBlockHash = "cd229e4ad2753708e4bab01d7f249affe29441829524c9529e84d51b6d12f2a7"

	// RosettaVersion -
	RosettaVersion = "1.4.0"
	// ElrondBlockchainName is the name of the Elrond blockchain
	ElrondBlockchainName = "Elrond"

	// OpStatusOK is the operation status for successful operations.
	OpStatusSuccess = "Success"
	// OpStatusFailed is the operation status for failed operations.
	OpStatusFailed = "Failed"

	opTransfer = "Transfer"
	opFee      = "Fee"
	opReward   = "Reward"
	opScResult = "SmartContractResult"
)
