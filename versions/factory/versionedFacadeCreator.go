package factory

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/api"
	apiv_next "github.com/multiversx/mx-chain-proxy-go/api/groups/v_next"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/facade"
	facadeVersions "github.com/multiversx/mx-chain-proxy-go/facade/versions"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/v_next"
	"github.com/multiversx/mx-chain-proxy-go/versions"
)

// FacadeArgs holds the arguments needed for creating a base facade
type FacadeArgs struct {
	ActionsProcessor             facade.ActionsProcessor
	AccountProcessor             facade.AccountProcessor
	FaucetProcessor              facade.FaucetProcessor
	BlockProcessor               facade.BlockProcessor
	BlocksProcessor              facade.BlocksProcessor
	NodeGroupProcessor           facade.NodeGroupProcessor
	NodeStatusProcessor          facade.NodeStatusProcessor
	ScQueryProcessor             facade.SCQueryService
	TransactionProcessor         facade.TransactionProcessor
	ValidatorStatisticsProcessor facade.ValidatorStatisticsProcessor
	ProofProcessor               facade.ProofProcessor
	PubKeyConverter              core.PubkeyConverter
	ESDTSuppliesProcessor        facade.ESDTSupplyProcessor
	StatusProcessor              facade.StatusProcessor
	AboutInfoProcessor           facade.AboutInfoProcessor
}

// CreateVersionsRegistry creates the version registry instances and populates it with the versions and their handlers
func CreateVersionsRegistry(facadeArgs FacadeArgs, apiConfigParser ApiConfigParser) (data.VersionsRegistryHandler, error) {
	versionsRegistry := versions.NewVersionsRegistry()

	err := addVersionV1_0(facadeArgs, versionsRegistry, apiConfigParser)
	if err != nil {
		return nil, err
	}

	err = addVersionV1_0AsDefault(versionsRegistry, apiConfigParser)
	if err != nil {
		return nil, err
	}

	// un-comment these lines if you want to start proxy also with the v_next

	// err = addVersionV_next(facadeArgs, versionsRegistry)
	// if err != nil {
	//	return nil, err
	// }

	return versionsRegistry, nil
}

func addVersionV1_0AsDefault(versionRegistry data.VersionsRegistryHandler, apiConfigParser ApiConfigParser) error {
	versionsMap, err := versionRegistry.GetAllVersions()
	if err != nil {
		return err
	}

	v1_0handler, ok := versionsMap["v1.0"]
	if !ok {
		return versions.ErrVersionNotFound
	}

	return versionRegistry.AddVersion("", v1_0handler)
}

func addVersionV1_0(facadeArgs FacadeArgs, versionRegistry data.VersionsRegistryHandler, apiConfigParser ApiConfigParser) error {
	v1_0Facade, err := createVersionV1_0Facade(facadeArgs)
	if err != nil {
		return err
	}

	apiHandler, err := api.NewApiHandler(v1_0Facade)
	if err != nil {
		return err
	}

	apiConfig, err := apiConfigParser.GetConfigForVersion("v1_0")
	if err != nil {
		return err
	}

	return versionRegistry.AddVersion("v1.0",
		&data.VersionData{
			Facade:     v1_0Facade,
			ApiHandler: apiHandler,
			ApiConfig:  *apiConfig,
		},
	)
}

func createVersionV1_0Facade(facadeArgs FacadeArgs) (*facadeVersions.ProxyFacadeV1_0, error) {
	v1_0HandlerArgs := FacadeArgs{
		ActionsProcessor:             facadeArgs.ActionsProcessor,
		AccountProcessor:             facadeArgs.AccountProcessor,
		FaucetProcessor:              facadeArgs.FaucetProcessor,
		BlockProcessor:               facadeArgs.BlockProcessor,
		BlocksProcessor:              facadeArgs.BlocksProcessor,
		NodeGroupProcessor:           facadeArgs.NodeGroupProcessor,
		NodeStatusProcessor:          facadeArgs.NodeStatusProcessor,
		ScQueryProcessor:             facadeArgs.ScQueryProcessor,
		TransactionProcessor:         facadeArgs.TransactionProcessor,
		ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor,
		ProofProcessor:               facadeArgs.ProofProcessor,
		PubKeyConverter:              facadeArgs.PubKeyConverter,
		ESDTSuppliesProcessor:        facadeArgs.ESDTSuppliesProcessor,
		StatusProcessor:              facadeArgs.StatusProcessor,
		AboutInfoProcessor:           facadeArgs.AboutInfoProcessor,
	}

	commonFacade, err := createVersionedFacade(v1_0HandlerArgs)
	if err != nil {
		return nil, err
	}

	return &facadeVersions.ProxyFacadeV1_0{ProxyFacade: commonFacade.(*facade.ProxyFacade)}, nil
}

func addVersionV_next(facadeArgs FacadeArgs, versionsRegistry data.VersionsRegistryHandler) error {
	v_nextHandler, err := createVersionV_nextFacade(facadeArgs)
	if err != nil {
		return err
	}

	apiHandler, err := api.NewApiHandler(v_nextHandler)
	if err != nil {
		return err
	}

	accountsGroup, err := apiHandler.GetGroup("/address")
	if err != nil {
		return err
	}

	accountsGroupV_next, err := apiv_next.NewAccountsGroupV_next(accountsGroup, v_nextHandler)
	if err != nil {
		return err
	}

	err = apiHandler.UpdateGroup("/address", accountsGroupV_next.Group())
	if err != nil {
		return err
	}

	return versionsRegistry.AddVersion("v_next",
		&data.VersionData{
			Facade:     v_nextHandler,
			ApiHandler: apiHandler,
		},
	)
}

func createVersionV_nextFacade(facadeArgs FacadeArgs) (data.FacadeHandler, error) {
	v_nextHandlerArgs := FacadeArgs{
		AccountProcessor:             facadeArgs.AccountProcessor,
		FaucetProcessor:              facadeArgs.FaucetProcessor,
		BlockProcessor:               facadeArgs.BlockProcessor,
		BlocksProcessor:              facadeArgs.BlocksProcessor,
		NodeGroupProcessor:           facadeArgs.NodeGroupProcessor,
		NodeStatusProcessor:          facadeArgs.NodeStatusProcessor,
		ScQueryProcessor:             facadeArgs.ScQueryProcessor,
		TransactionProcessor:         facadeArgs.TransactionProcessor,
		ValidatorStatisticsProcessor: facadeArgs.ValidatorStatisticsProcessor,
		PubKeyConverter:              facadeArgs.PubKeyConverter,
		ESDTSuppliesProcessor:        facadeArgs.ESDTSuppliesProcessor,
		StatusProcessor:              facadeArgs.StatusProcessor,
	}

	commonFacade, err := createVersionedFacade(v_nextHandlerArgs)
	if err != nil {
		return nil, err
	}

	newAccountsProcessor := v_next.AccountProcessorV_next{
		AccountProcessor: facadeArgs.AccountProcessor.(*process.AccountProcessor),
	}
	customFacade := &facadeVersions.ProxyFacadeV_next{
		ProxyFacade:       commonFacade.(*facade.ProxyFacade),
		AccountsProcessor: newAccountsProcessor,
	}

	return customFacade, nil
}

func createVersionedFacade(args FacadeArgs) (data.FacadeHandler, error) {
	// no need to check the arguments because they are initiated before arriving here and we assume that the constructor
	// always returns a good instance of the object (or an error otherwise)
	// Also, there are nil checks on the facade's constructors

	return facade.NewProxyFacade(
		args.ActionsProcessor,
		args.AccountProcessor,
		args.TransactionProcessor,
		args.ScQueryProcessor,
		args.NodeGroupProcessor,
		args.ValidatorStatisticsProcessor,
		args.FaucetProcessor,
		args.NodeStatusProcessor,
		args.BlockProcessor,
		args.BlocksProcessor,
		args.ProofProcessor,
		args.PubKeyConverter,
		args.ESDTSuppliesProcessor,
		args.StatusProcessor,
		args.AboutInfoProcessor,
	)
}
