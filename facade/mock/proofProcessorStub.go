package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// ProofProcessorStub -
type ProofProcessorStub struct {
	GetProofCalled                func(string, string) (*data.GenericAPIResponse, error)
	GetProofDataTrieCalled        func(string, string, string) (*data.GenericAPIResponse, error)
	GetProofCurrentRootHashCalled func(string) (*data.GenericAPIResponse, error)
	VerifyProofCalled             func(string, string, []string) (*data.GenericAPIResponse, error)
}

// GetProof -
func (pp *ProofProcessorStub) GetProof(rootHash string, address string) (*data.GenericAPIResponse, error) {
	if pp.GetProofCalled != nil {
		return pp.GetProofCalled(rootHash, address)
	}

	return nil, nil
}

// GetProofDataTrie -
func (pp *ProofProcessorStub) GetProofDataTrie(rootHash string, address string, key string) (*data.GenericAPIResponse, error) {
	if pp.GetProofDataTrieCalled != nil {
		return pp.GetProofDataTrieCalled(rootHash, address, key)
	}

	return nil, nil
}

// GetProofCurrentRootHash -
func (pp *ProofProcessorStub) GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error) {
	if pp.GetProofCurrentRootHashCalled != nil {
		return pp.GetProofCurrentRootHashCalled(address)
	}

	return nil, nil
}

// VerifyProof -
func (pp *ProofProcessorStub) VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error) {
	if pp.VerifyProofCalled != nil {
		return pp.VerifyProofCalled(rootHash, address, proof)
	}

	return nil, nil
}
