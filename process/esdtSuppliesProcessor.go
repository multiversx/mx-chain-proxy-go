package process

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type esdtSuppliesProcessor struct {
	baseProc    Processor
	scQueryProc SCQueryService
}

func NewESDTSuppliesProcessor(baseProc Processor, scQueryProc SCQueryService) (*esdtSuppliesProcessor, error) {
	if baseProc == nil {
		return nil, ErrNilCoreProcessor
	}
	if scQueryProc == nil {
		return nil, ErrNilSCQueryService
	}

	return &esdtSuppliesProcessor{
		baseProc:    baseProc,
		scQueryProc: scQueryProc,
	}, nil
}

func (esp *esdtSuppliesProcessor) GetESDTSupply(token string) (*data.GenericAPIResponse, error) {
	return nil, nil
}
