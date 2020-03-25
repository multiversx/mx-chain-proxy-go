package data

// Account defines the data structure for an account
type Account struct {
	Address  string `json:"address"`
	Nonce    uint64 `json:"nonce"`
	Balance  string `json:"balance"`
	Code     string `json:"code"`
	CodeHash []byte `json:"codeHash"`
	RootHash []byte `json:"rootHash"`
}

// ValidatorApiResponse represents the data which is fetched from each validator for returning it in API call
type ValidatorApiResponse struct {
	NrLeaderSuccess    uint32  `json:"nrLeaderSuccess"`
	NrLeaderFailure    uint32  `json:"nrLeaderFailure"`
	NrValidatorSuccess uint32  `json:"nrValidatorSuccess"`
	NrValidatorFailure uint32  `json:"nrValidatorFailure"`
	Rating             float32 `json:"rating"`
	TempRating         float32 `json:"tempRating"`
}

// ResponseAccount defines a wrapped account that the node respond with
type ResponseAccount struct {
	AccountData Account `json:"account"`
}
