package mock

// ProofProcessorStub -
type ProofProcessorStub struct {
	GetProofCalled    func([]byte, []byte) ([][]byte, error)
	VerifyProofCalled func([]byte, []byte, [][]byte) (bool, error)
}

// GetProof -
func (pp *ProofProcessorStub) GetProof(rootHash []byte, address []byte) ([][]byte, error) {
	if pp.GetProofCalled != nil {
		return pp.GetProofCalled(rootHash, address)
	}

	return nil, nil
}

// VerifyProof -
func (pp *ProofProcessorStub) VerifyProof(rootHash []byte, address []byte, proof [][]byte) (bool, error) {
	if pp.VerifyProofCalled != nil {
		return pp.VerifyProofCalled(rootHash, address, proof)
	}

	return false, nil
}
