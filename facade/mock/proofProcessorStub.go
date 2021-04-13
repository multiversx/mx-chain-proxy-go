package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// ProofProcessorStub -
type ProofProcessorStub struct {
	GetProofCalled                func([]byte, []byte) (*data.GenericAPIResponse, error)
	GetProofCurrentRootHashCalled func([]byte) (*data.GenericAPIResponse, error)
	VerifyProofCalled             func([]byte, []byte, []string) (*data.GenericAPIResponse, error)
}

// GetProof -
func (pp *ProofProcessorStub) GetProof(rootHash []byte, address []byte) (*data.GenericAPIResponse, error) {
	if pp.GetProofCalled != nil {
		return pp.GetProofCalled(rootHash, address)
	}

	return nil, nil
}

// GetProofCurrentRootHash -
func (pp *ProofProcessorStub) GetProofCurrentRootHash(address []byte) (*data.GenericAPIResponse, error) {
	if pp.GetProofCurrentRootHashCalled != nil {
		return pp.GetProofCurrentRootHashCalled(address)
	}

	return nil, nil
}

// VerifyProof -
func (pp *ProofProcessorStub) VerifyProof(rootHash []byte, address []byte, proof []string) (*data.GenericAPIResponse, error) {
	if pp.VerifyProofCalled != nil {
		return pp.VerifyProofCalled(rootHash, address, proof)
	}

	return nil, nil
}
