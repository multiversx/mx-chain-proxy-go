package process

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type ProofProcessor struct {
	proc            Processor
	pubKeyConverter core.PubkeyConverter
}

func NewProofProcessor(proc Processor, pubKeyConverter core.PubkeyConverter) (*ProofProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &ProofProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// GetProof sends the request to the right observer and then replies with the returned answer
func (pp *ProofProcessor) GetProof(rootHash string, address string) (*data.GenericAPIResponse, error) {
	observers, err := pp.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	responseGetProof := data.GenericAPIResponse{}
	getProofEndpoint := "/proof/root-hash/" + rootHash + "/address/" + address
	for _, observer := range observers {

		respCode, err := pp.proc.CallGetRestEndPoint(observer.Address, getProofEndpoint, &responseGetProof)

		if responseGetProof.Error != "" {
			return nil, errors.New(responseGetProof.Error)
		}

		if err != nil {
			log.Error("GetProof request",
				"observer", observer.Address,
				"address", address,
				"error", err.Error(),
			)

			continue
		}

		if respCode == http.StatusOK {
			log.Info("GetProof request",
				"address", address,
				"rootHash", rootHash,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode,
			)

			return &responseGetProof, nil
		}
	}

	return nil, WrapObserversError(responseGetProof.Error)
}

// GetProofDataTrie sends the request to the right observer and then replies with the returned answer
func (pp *ProofProcessor) GetProofDataTrie(rootHash string, address string, key string) (*data.GenericAPIResponse, error) {
	observers, err := pp.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	responseGetProof := data.GenericAPIResponse{}
	getProofDataTrieEndpoint := fmt.Sprintf("/proof/root-hash/%s/address/%s/key/%s", rootHash, address, key)
	for _, observer := range observers {

		respCode, err := pp.proc.CallGetRestEndPoint(observer.Address, getProofDataTrieEndpoint, &responseGetProof)

		if responseGetProof.Error != "" {
			return nil, errors.New(responseGetProof.Error)
		}

		if err != nil {
			log.Error("GetProofDataTrie request",
				"observer", observer.Address,
				"address", address,
				"error", err.Error(),
			)

			continue
		}

		if respCode == http.StatusOK {
			log.Info("GetProofDataTrie request",
				"address", address,
				"rootHash", rootHash,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode,
			)

			return &responseGetProof, nil
		}
	}

	return nil, WrapObserversError(responseGetProof.Error)
}

// GetProofCurrentRootHash sends the request to the right observer and then replies with the returned answer
func (pp *ProofProcessor) GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error) {
	observers, err := pp.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	responseGetProof := data.GenericAPIResponse{}
	getProofEndpoint := "/proof/address/" + address
	for _, observer := range observers {

		respCode, err := pp.proc.CallGetRestEndPoint(observer.Address, getProofEndpoint, &responseGetProof)

		if responseGetProof.Error != "" {
			return nil, errors.New(responseGetProof.Error)
		}

		if err != nil {
			log.Error("GetProofCurrentRootHash request",
				"observer", observer.Address,
				"address", address,
				"error", err.Error(),
			)

			continue
		}

		if respCode == http.StatusOK {
			log.Info("GetProof request",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode,
			)

			return &responseGetProof, nil
		}
	}

	return nil, WrapObserversError(responseGetProof.Error)
}

// VerifyProof sends the request to the right observer and then replies with the returned answer
func (pp *ProofProcessor) VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error) {
	observers, err := pp.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	verifyProofEndpoint := "/proof/verify"
	requestParams := data.VerifyProofRequest{
		RootHash: rootHash,
		Address:  address,
		Proof:    proof,
	}
	responseVerifyProof := data.GenericAPIResponse{}
	for _, observer := range observers {

		respCode, err := pp.proc.CallPostRestEndPoint(observer.Address, verifyProofEndpoint, requestParams, &responseVerifyProof)

		if responseVerifyProof.Error != "" {
			return nil, errors.New(responseVerifyProof.Error)
		}

		if err != nil {
			log.Error("VerifyProof request",
				"observer", observer.Address,
				"address", address,
				"error", err.Error(),
			)

			continue
		}

		if respCode == http.StatusOK {
			log.Info("VerifyProof request",
				"address", address,
				"rootHash", rootHash,
				"proof", proof,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode,
			)

			return &responseVerifyProof, nil
		}
	}

	return nil, WrapObserversError(responseVerifyProof.Error)
}

func (pp *ProofProcessor) getObserversForAddress(address string) ([]*data.NodeData, error) {
	addressBytes, err := pp.pubKeyConverter.Decode(address)
	if err != nil {
		return nil, err
	}

	shardID, err := pp.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	observers, err := pp.proc.GetObservers(shardID, data.AvailabilityAll)
	if err != nil {
		return nil, err
	}

	return observers, nil
}
