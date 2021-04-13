package data

// VerifyProofRequest represents the parameters needed to verify a Merkle proof
type VerifyProofRequest struct {
	RootHash []byte   `json:"roothash"`
	Address  []byte   `json:"address"`
	Proof    []string `json:"proof"`
}
