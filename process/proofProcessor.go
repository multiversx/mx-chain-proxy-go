package process

import (
	"errors"
	"net/http"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
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
func (pp *ProofProcessor) GetProof(rootHash []byte, address []byte) (*data.GenericAPIResponse, error) {
	observers, err := pp.getObserversForAddress(string(address))
	if err != nil {
		return nil, err
	}

	getProofEndpoint := "/proof/root-hash/" + string(rootHash) + "/address/" + string(address)
	for _, observer := range observers {
		responseGetProof := &data.GenericAPIResponse{}

		respCode, err := pp.proc.CallGetRestEndPoint(observer.Address, getProofEndpoint, responseGetProof)

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

			return responseGetProof, nil
		}
	}

	return nil, ErrSendingRequest
}

// GetProofCurrentRootHash sends the request to the right observer and then replies with the returned answer
func (pp *ProofProcessor)GetProofCurrentRootHash(address []byte) (*data.GenericAPIResponse, error){
	observers, err := pp.getObserversForAddress(string(address))
	if err != nil {
		return nil, err
	}

	getProofEndpoint := "/proof/address/" + string(address)
	for _, observer := range observers {
		responseGetProof := &data.GenericAPIResponse{}

		respCode, err := pp.proc.CallGetRestEndPoint(observer.Address, getProofEndpoint, responseGetProof)

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

			return responseGetProof, nil
		}
	}

	return nil, ErrSendingRequest
}

// VerifyProof sends the request to the right observer and then replies with the returned answer
func (pp *ProofProcessor) VerifyProof(rootHash []byte, address []byte, proof []string) (*data.GenericAPIResponse, error) {
	observers, err := pp.getObserversForAddress(string(address))
	if err != nil {
		return nil, err
	}

	verifyProofEndpoint := "/proof/verify"
	requestParams := data.VerifyProofRequest{
		RootHash: rootHash,
		Address:  address,
		Proof:    proof,
	}
	for _, observer := range observers {
		responseVerifyProof := &data.GenericAPIResponse{}

		respCode, err := pp.proc.CallPostRestEndPoint(observer.Address, verifyProofEndpoint, requestParams, responseVerifyProof)

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

			return responseVerifyProof, nil
		}
	}

	return nil, ErrSendingRequest
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

	observers, err := pp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	return observers, nil
}
