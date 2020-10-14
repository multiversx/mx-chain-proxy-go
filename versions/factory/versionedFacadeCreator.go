package factory

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	versions2 "github.com/ElrondNetwork/elrond-proxy-go/facade/versions"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/v1_0"
	"github.com/ElrondNetwork/elrond-proxy-go/process/v1_1"
	"github.com/ElrondNetwork/elrond-proxy-go/versions"
)

type FacadeArgs struct {
	AccountProcessor             facade.AccountProcessor
	FaucetProcessor              facade.FaucetProcessor
	BlockProcessor               facade.BlockProcessor
	HeartbeatProcessor           facade.HeartbeatProcessor
	NodeStatusProcessor          facade.NodeStatusProcessor
	ScQueryProcessor             facade.SCQueryService
	TransactionProcessor         facade.TransactionProcessor
	ValidatorStatisticsProcessor facade.ValidatorStatisticsProcessor
}

func CreateVersionManager(facadeArgs FacadeArgs, commonApiHandler data.ApiHandler) (data.VersionManagerHandler, error) {
	versionManager := versions.NewVersionManager(commonApiHandler)

	v1_0Handler, err := createVersionV1_0(facadeArgs)
	if err != nil {
		return nil, err
	}
	err = versionManager.AddVersion("v1.0", &data.VersionData{Facade: v1_0Handler})
	if err != nil {
		return nil, err
	}

	v1_1Handler, err := createVersionV1_1(facadeArgs)
	if err != nil {
		return nil, err
	}
	err = versionManager.AddVersion("v1.1", &data.VersionData{Facade: v1_1Handler})
	if err != nil {
		return nil, err
	}

	return versionManager, nil
}

func createVersionV1_0(facadeArgs FacadeArgs) (data.FacadeHandler, error) {
	v1_0HandlerArgs := FacadeArgs{
		AccountProcessor: &v1_0.AccountProcessorV1_0{
			AccountProcessor: facadeArgs.AccountProcessor.(*process.AccountProcessor),
		},
		FaucetProcessor: &v1_0.FaucetProcessorV1_0{
			FaucetProcessor: facadeArgs.FaucetProcessor.(*process.FaucetProcessor),
		},
		BlockProcessor: &v1_0.BlockProcessorV1_0{
			BlockProcessor: facadeArgs.BlockProcessor.(*process.BlockProcessor),
		},
		HeartbeatProcessor: &v1_0.HeartbeatProcessorV1_0{
			HeartbeatProcessor: facadeArgs.HeartbeatProcessor.(*process.HeartbeatProcessor),
		},
		NodeStatusProcessor: &v1_0.NodeStatusProcessorV1_0{
			NodeStatusProcessor: facadeArgs.NodeStatusProcessor.(*process.NodeStatusProcessor),
		},
		ScQueryProcessor: &v1_0.ScQueryProcessorV1_0{
			SCQueryProcessor: facadeArgs.ScQueryProcessor.(*process.SCQueryProcessor),
		},
		TransactionProcessor: &v1_0.TransactionProcessorV1_0{
			TransactionProcessor: facadeArgs.TransactionProcessor.(*process.TransactionProcessor),
		},
		ValidatorStatisticsProcessor: &v1_0.ValidatorStatisticsProcessorV1_0{
			ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor.(*process.ValidatorStatisticsProcessor),
		},
	}

	commonFacade, err := createVersionedFacade(v1_0HandlerArgs)
	if err != nil {
		return nil, err
	}

	return versions2.ElrondProxyFacadeV1_0{ElrondProxyFacade: commonFacade.(*facade.ElrondProxyFacade)}, nil
}

func createVersionV1_1(facadeArgs FacadeArgs) (data.FacadeHandler, error) {
	v1_1HandlerArgs := FacadeArgs{
		AccountProcessor: &v1_1.AccountProcessorV1_1{
			AccountProcessor: facadeArgs.AccountProcessor.(*process.AccountProcessor),
		},
		FaucetProcessor: &v1_1.FaucetProcessorV1_1{
			FaucetProcessor: facadeArgs.FaucetProcessor.(*process.FaucetProcessor),
		},
		BlockProcessor: &v1_1.BlockProcessorV1_1{
			BlockProcessor: facadeArgs.BlockProcessor.(*process.BlockProcessor),
		},
		HeartbeatProcessor: &v1_1.HeartbeatProcessorV1_1{
			HeartbeatProcessor: facadeArgs.HeartbeatProcessor.(*process.HeartbeatProcessor),
		},
		NodeStatusProcessor: &v1_1.NodeStatusProcessorV1_1{
			NodeStatusProcessor: facadeArgs.NodeStatusProcessor.(*process.NodeStatusProcessor),
		},
		ScQueryProcessor: &v1_1.ScQueryProcessorV1_1{
			SCQueryProcessor: facadeArgs.ScQueryProcessor.(*process.SCQueryProcessor),
		},
		TransactionProcessor: &v1_1.TransactionProcessorV1_1{
			TransactionProcessor: facadeArgs.TransactionProcessor.(*process.TransactionProcessor),
		},
		ValidatorStatisticsProcessor: &v1_1.ValidatorStatisticsProcessorV1_1{
			ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor.(*process.ValidatorStatisticsProcessor),
		},
	}

	commonFacade, err := createVersionedFacade(v1_1HandlerArgs)
	if err != nil {
		return nil, err
	}

	return versions2.ElrondProxyFacadeV1_1{ElrondProxyFacade: commonFacade.(*facade.ElrondProxyFacade)}, nil
}

func createVersionedFacade(args FacadeArgs) (data.FacadeHandler, error) {
	// no need to check the arguments because they are initiated before arriving here and we assume that the constructor
	// always returns a good instance of the object (or an error otherwise)
	// Also, there are nil checks on the facade's constructors

	return facade.NewElrondProxyFacade(
		args.AccountProcessor,
		args.TransactionProcessor,
		args.ScQueryProcessor,
		args.HeartbeatProcessor,
		args.ValidatorStatisticsProcessor,
		args.FaucetProcessor,
		args.NodeStatusProcessor,
		args.BlockProcessor,
	)
}
