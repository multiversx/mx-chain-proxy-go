package data

// ApiBlock represents the structure of a block as it is sent from API
type ApiBlock struct {
	Nonce        uint64           `form:"nonce" json:"nonce"`
	Hash         string           `form:"hash" json:"hash"`
	Transactions []ApiTransaction `form:"transactions" json:"transactions,omitempty"`
}
