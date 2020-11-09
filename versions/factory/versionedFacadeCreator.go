package factory

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	apiv1_2 "github.com/ElrondNetwork/elrond-proxy-go/api/groups/v1_2"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	facadeVersions "github.com/ElrondNetwork/elrond-proxy-go/facade/versions"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/v1_0"
	"github.com/ElrondNetwork/elrond-proxy-go/process/v1_2"
	"github.com/ElrondNetwork/elrond-proxy-go/versions"
)

// FacadeArgs holds the arguments needed for creating a base facade
type FacadeArgs struct {
	AccountProcessor             facade.AccountProcessor
	FaucetProcessor              facade.FaucetProcessor
	BlockProcessor               facade.BlockProcessor
	HeartbeatProcessor           facade.HeartbeatProcessor
	NodeStatusProcessor          facade.NodeStatusProcessor
	ScQueryProcessor             facade.SCQueryService
	TransactionProcessor         facade.TransactionProcessor
	ValidatorStatisticsProcessor facade.ValidatorStatisticsProcessor
	PubKeyConverter              core.PubkeyConverter
}

// CreateVersionManager creates the version manager instances and populates it with the versions and their handlers
func CreateVersionManager(facadeArgs FacadeArgs) (data.VersionManagerHandler, error) {
	versionManager := versions.NewVersionManager()

	err := addVersionV1_0(facadeArgs, versionManager)
	if err != nil {
		return nil, err
	}

	err = addVersionV1_1(facadeArgs, versionManager)
	if err != nil {
		return nil, err
	}

	err = addVersionV1_2(facadeArgs, versionManager)
	if err != nil {
		return nil, err
	}

	return versionManager, nil
}

func addVersionV1_0(facadeArgs FacadeArgs, versionManager data.VersionManagerHandler) error {
	v1_0Facade, err := createVersionV1_0Facade(facadeArgs)
	if err != nil {
		return err
	}

	apiHandler, err := api.NewApiHandler(v1_0Facade)
	if err != nil {
		return err
	}

	return versionManager.AddVersion("v1.0",
		&data.VersionData{
			Facade:     v1_0Facade,
			ApiHandler: apiHandler,
		},
	)
}

func createVersionV1_0Facade(facadeArgs FacadeArgs) (*facadeVersions.ElrondProxyFacadeV1_0, error) {
	v1_0HandlerArgs := FacadeArgs{
		AccountProcessor:    facadeArgs.AccountProcessor,
		FaucetProcessor:     facadeArgs.FaucetProcessor,
		BlockProcessor:      facadeArgs.BlockProcessor,
		HeartbeatProcessor:  facadeArgs.HeartbeatProcessor,
		NodeStatusProcessor: facadeArgs.NodeStatusProcessor,
		ScQueryProcessor:    facadeArgs.ScQueryProcessor,
		TransactionProcessor: &v1_0.TransactionProcessorV1_0{
			TransactionProcessor: facadeArgs.TransactionProcessor.(*process.TransactionProcessor),
		},
		ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor,
		PubKeyConverter:              facadeArgs.PubKeyConverter,
	}

	commonFacade, err := createVersionedFacade(v1_0HandlerArgs)
	if err != nil {
		return nil, err
	}

	return &facadeVersions.ElrondProxyFacadeV1_0{ElrondProxyFacade: commonFacade.(*facade.ElrondProxyFacade)}, nil
}

func addVersionV1_1(facadeArgs FacadeArgs, versionManager data.VersionManagerHandler) error {
	v1_1Facade, err := createVersionV1_1Facade(facadeArgs)
	if err != nil {
		return err
	}

	apiHandler, err := api.NewApiHandler(v1_1Facade)
	if err != nil {
		return err
	}

	return versionManager.AddVersion("v1.1",
		&data.VersionData{
			Facade:     v1_1Facade,
			ApiHandler: apiHandler,
		})
}

func createVersionV1_1Facade(facadeArgs FacadeArgs) (data.FacadeHandler, error) {
	v1_1HandlerArgs := FacadeArgs{
		AccountProcessor:             facadeArgs.AccountProcessor,
		FaucetProcessor:              facadeArgs.FaucetProcessor,
		BlockProcessor:               facadeArgs.BlockProcessor,
		HeartbeatProcessor:           facadeArgs.HeartbeatProcessor,
		NodeStatusProcessor:          facadeArgs.NodeStatusProcessor,
		ScQueryProcessor:             facadeArgs.ScQueryProcessor,
		TransactionProcessor:         facadeArgs.TransactionProcessor,
		ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor,
		PubKeyConverter:              facadeArgs.PubKeyConverter,
	}

	commonFacade, err := createVersionedFacade(v1_1HandlerArgs)
	if err != nil {
		return nil, err
	}

	return facadeVersions.ElrondProxyFacadeV1_1{ElrondProxyFacade: commonFacade.(*facade.ElrondProxyFacade)}, nil
}

func addVersionV1_2(facadeArgs FacadeArgs, versionManager data.VersionManagerHandler) error {
	v1_2Handler, err := createVersionV1_2Facade(facadeArgs)
	if err != nil {
		return err
	}

	apiHandler, err := api.NewApiHandler(v1_2Handler)
	if err != nil {
		return err
	}

	accountsGroup, err := apiHandler.GetGroup("/address")
	if err != nil {
		return err
	}

	accountsGroupV1_2, err := apiv1_2.NewAccountsGroupV1_2(accountsGroup, v1_2Handler)
	if err != nil {
		return err
	}

	err = apiHandler.UpdateGroup("/address", accountsGroupV1_2.Group())
	if err != nil {
		return err
	}

	return versionManager.AddVersion("v1.2",
		&data.VersionData{
			Facade:     v1_2Handler,
			ApiHandler: apiHandler,
		},
	)
}

func createVersionV1_2Facade(facadeArgs FacadeArgs) (data.FacadeHandler, error) {
	v1_2HandlerArgs := FacadeArgs{
		AccountProcessor:             facadeArgs.AccountProcessor,
		FaucetProcessor:              facadeArgs.FaucetProcessor,
		BlockProcessor:               facadeArgs.BlockProcessor,
		HeartbeatProcessor:           facadeArgs.HeartbeatProcessor,
		NodeStatusProcessor:          facadeArgs.NodeStatusProcessor,
		ScQueryProcessor:             facadeArgs.ScQueryProcessor,
		TransactionProcessor:         facadeArgs.TransactionProcessor,
		ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor,
		PubKeyConverter:              facadeArgs.PubKeyConverter,
	}

	commonFacade, err := createVersionedFacade(v1_2HandlerArgs)
	if err != nil {
		return nil, err
	}

	newAccountsProcessor := v1_2.AccountProcessorV1_2{
		AccountProcessor: facadeArgs.AccountProcessor.(*process.AccountProcessor),
	}
	customFacade := &facadeVersions.ElrondProxyFacadeV1_2{
		ElrondProxyFacade: commonFacade.(*facade.ElrondProxyFacade),
		AccountsProcessor: newAccountsProcessor,
	}

	return customFacade, nil
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
		args.PubKeyConverter,
	)
}
