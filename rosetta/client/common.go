package client

const (
	MetachainID           = 4294967295
	RoundDurationInSecond = int64(6)
	ElrondGenesisTime     = int64(1596117600)
)

type NetworkConfig struct {
	ChainID        string `json:"erd_chain_id"`
	Denomination   uint64 `json:"erd_denomination"`
	GasPerDataByte uint64 `json:"erd_gas_per_data_byte"`
	ClientVersion  string `json:"erd_latest_tag_software_version"`
	MinGasPrice    uint64 `json:"erd_min_gas_price"`
	MinGasLimit    uint64 `json:"erd_min_gas_limit"`
	MinTxVersion   uint32 `json:"erd_min_transaction_version"`
}

type NetworkStatus struct {
	CurrentNonce uint64 `json:"erd_nonce"`
	CurrentEpoch uint64 ` json:"erd_epoch_number"`
}

type BlockData struct {
	Nonce         uint64
	Hash          string
	PrevBlockHash string
	Timestamp     int64
}

func CalculateBlockTimestampUnix(round uint64) int64 {
	return (ElrondGenesisTime + int64(round)*RoundDurationInSecond) * 1000
}
