package data

// Account defines the data structure for an account
type Account struct {
	Address  string `json:"address"`
	Nonce    uint64 `json:"nonce"`
	Balance  string `json:"balance"`
	CodeHash []byte `json:"codeHash"`
	RootHash []byte `json:"rootHash"`
}
