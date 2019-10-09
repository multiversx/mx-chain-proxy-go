package process

import (
	"encoding/hex"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AddressPath defines the address path at which the nodes answer
const AddressPath = "/address/"

// AccountProcessor is able to process account requests
type AccountProcessor struct {
	proc   Processor
	keyGen crypto.KeyGenerator
}

// NewAccountProcessor creates a new instance of AccountProcessor
func NewAccountProcessor(proc Processor, keyGen crypto.KeyGenerator) (*AccountProcessor, error) {
	if proc == nil {
		return nil, ErrNilCoreProcessor
	}
	if keyGen == nil {
		return nil, ErrNilKeyGen
	}

	return &AccountProcessor{
		proc:   proc,
		keyGen: keyGen,
	}, nil
}

// GetAccount resolves the request by sending the request to the right observer and replies back the answer
func (ap *AccountProcessor) GetAccount(address string) (*data.Account, error) {
	addressBytes, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	shardId, err := ap.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	observers, err := ap.proc.GetObservers(shardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		responseAccount := &data.ResponseAccount{}

		err = ap.proc.CallGetRestEndPoint(observer.Address, AddressPath+address, responseAccount)
		if err == nil {
			log.Info(fmt.Sprintf("Got account request from observer %v from shard %v", observer.Address, shardId))
			return &responseAccount.AccountData, nil
		}

		log.LogIfError(err)
	}

	return nil, ErrSendingRequest
}

// PublicKeyFromPrivateKey will return the public key corresponding to the private key
func (ap *AccountProcessor) PublicKeyFromPrivateKey(privateKeyHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	privKey, err := ap.keyGen.PrivateKeyFromByteArray(privKeyBytes)
	if err != nil {
		return "", err
	}

	publicKey := privKey.GeneratePublic()
	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(publicKeyBytes), nil
}
