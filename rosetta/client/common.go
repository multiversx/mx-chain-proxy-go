package client

const (
	MetachainID = 4294967295
)

type NetworkConfig struct {
	ChainID        string `json:"erd_chain_id"`
	Denomination   uint64 `json:"erd_denomination"`
	GasPerDataByte uint64 `json:"erd_gas_per_data_byte"`
	ClientVersion  string `json:"erd_latest_tag_software_version"`
	MinGasPrice    uint64 `json:"erd_min_gas_price"`
	MinGasLimit    uint64 `json:"erd_min_gas_limit"`
	MinTxVersion   uint32 `json:"erd_min_transaction_version"`
	StartTime      uint64 `json:"erd_start_time"`
	RoundDuration  uint64 `json:"erd_round_duration"`
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
