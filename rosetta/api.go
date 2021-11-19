package rosetta

import (
	"fmt"
	"net/http"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/services"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/services/offline"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

var log = logger.GetOrCreate("rosetta")

// CreateServer creates a HTTP server
func CreateServer(
	elrondFacade api.ElrondProxyHandler,
	generalConfig *config.Config,
	port int,
	isOffline bool,
	offlineConfigPath string,
) (*http.Server, error) {
	if !isOffline {
		return createOnlineServer(elrondFacade, generalConfig, port)
	}

	cfg, err := configuration.LoadOfflineConfig(generalConfig, offlineConfigPath)
	if err != nil {
		return nil, err
	}

	asserterServer, err := createAsserter(cfg.Network)
	if err != nil {
		return nil, err
	}

	offlineService := offline.NewOfflineService()

	accountAPIController := server.NewAccountAPIController(offlineService, asserterServer)
	blockAPIController := server.NewBlockAPIController(offlineService, asserterServer)
	mempoolAPIController := server.NewMempoolAPIController(offlineService, asserterServer)

	elrondProvider, err := provider.NewOfflineElrondProvider(elrondFacade, cfg.ElrondNetworkConfig)
	if err != nil {
		return nil, err
	}

	constructionAPIService := services.NewConstructionAPIService(elrondProvider, cfg, cfg.ElrondNetworkConfig, isOffline)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserterServer,
	)

	networkAPIService := services.NewNetworkAPIService(elrondProvider, cfg, true)
	networkAPIController := server.NewNetworkAPIController(networkAPIService, asserterServer)

	log.Info("elrond rosetta server is in offline mode")

	return createHttpServer(port, networkAPIController, accountAPIController, blockAPIController, constructionAPIController, mempoolAPIController)
}

func createOnlineServer(
	elrondFacade api.ElrondProxyHandler,
	generalConfig *config.Config,
	port int,
) (*http.Server, error) {
	elrondProvider, err := provider.NewElrondProvider(elrondFacade)
	if err != nil {
		log.Error("cannot create elrond provider", "err", err)
		return nil, err
	}

	networkConfig, err := elrondProvider.GetNetworkConfig()
	if err != nil {
		log.Error("cannot get network config", "err", err)
		return nil, err
	}

	cfg := configuration.LoadConfiguration(networkConfig, generalConfig)
	asserterServer, err := createAsserter(cfg.Network)
	if err != nil {
		return nil, err
	}

	// Create network service
	networkAPIService := services.NewNetworkAPIService(elrondProvider, cfg, false)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserterServer,
	)

	// Create account service
	accountAPIService := services.NewAccountAPIService(elrondProvider, cfg)
	accountAPIController := server.NewAccountAPIController(
		accountAPIService,
		asserterServer,
	)

	// Create block service
	blockAPIService := services.NewBlockAPIService(elrondProvider, cfg, networkConfig)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		asserterServer,
	)

	// Create construction service
	constructionAPIService := services.NewConstructionAPIService(elrondProvider, cfg, networkConfig, false)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserterServer,
	)

	// Create mempool service
	mempoolAPIService := services.NewMempoolApiService(elrondProvider, cfg, networkConfig)
	mempoolAPIController := server.NewMempoolAPIController(
		mempoolAPIService,
		asserterServer,
	)

	log.Info("rosetta server started in online mode")

	return createHttpServer(port, networkAPIController, accountAPIController, blockAPIController, constructionAPIController, mempoolAPIController)
}

func createHttpServer(port int, routers ...server.Router,
) (*http.Server, error) {
	router := server.NewRouter(
		routers...,
	)

	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsRouter,
	}

	return httpServer, nil
}

func createAsserter(network *types.NetworkIdentifier) (*asserter.Asserter, error) {
	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserterServer, err := asserter.NewServer(
		services.SupportedOperationTypes,
		false,
		[]*types.NetworkIdentifier{
			network,
		},
		nil,
		false,
		"",
	)
	if err != nil {
		return nil, err
	}

	return asserterServer, nil
}
