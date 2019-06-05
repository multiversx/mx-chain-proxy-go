package data

// Account defines the data structure for an account
type Account struct {
	Address  string `json:"address"`
	Nonce    uint64 `json:"nonce"`
	Balance  string `json:"balance"`
	CodeHash []byte `json:"codeHash"`
	RootHash []byte `json:"rootHash"`
}

// ResponseAccount defines a wrapped account that the node respond with
type ResponseAccount struct {
	Account Account `json:"account"`
}
