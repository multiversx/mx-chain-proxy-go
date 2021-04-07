package data

// GetProofResponse represents the response for the GetProof endpoint
type GetProofResponse struct {
	Data  [][]byte `json:"data"`
	Error string   `json:"error"`
	Code  string   `json:"code"`
}

// VerifyProofRequest represents the response for the VerifyProof endpoint
type VerifyProofResponse struct {
	Data  bool   `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// VerifyProofRequest represents the parameters needed to verify a Merkle proof
type VerifyProofRequest struct {
	RootHash []byte   `json:"roothash"`
	Address  []byte   `json:"address"`
	Proof    [][]byte `json:"proof"`
}
