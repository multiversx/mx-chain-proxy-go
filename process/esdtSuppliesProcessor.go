package process

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type esdtSuppliesProcessor struct {
}

func NewESDTSuppliesProcessor() (*esdtSuppliesProcessor, error) {
	return &esdtSuppliesProcessor{}, nil
}

func (esp *esdtSuppliesProcessor) GetESDTSupply(token string) (*data.GenericAPIResponse, error) {
	return nil, nil
}
