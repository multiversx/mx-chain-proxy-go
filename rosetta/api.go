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
) (*http.Server, error) {
	if !isOffline {
		return createServerOnline(elrondFacade, generalConfig, port)
	}

	cfg := configuration.LoadOfflineConfig(generalConfig)
	asserterServer, err := asserter.NewServer(services.SupportedOperationTypes,
		false,
		[]*types.NetworkIdentifier{
			cfg.Network,
		},
		nil,
		false,
		"",
	)
	if err != nil {
		log.Error("cannot create asserter", "err", err)
		return nil, err
	}

	offlineService := offline.NewOfflineService()

	accountAPIController := server.NewAccountAPIController(offlineService, asserterServer)
	blockAPIController := server.NewBlockAPIController(offlineService, asserterServer)
	mempoolAPIController := server.NewMempoolAPIController(offlineService, asserterServer)

	elrondProvider, err := provider.NewOfflineElrondProvider(elrondFacade, cfg.ElrondNetworkConfig)
	if err != nil {
		log.Error("cannot create elrond provider", "err", err)
		return nil, err
	}

	constructionAPIService := services.NewConstructionAPIService(elrondProvider, cfg, cfg.ElrondNetworkConfig, isOffline)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserterServer,
	)

	networkAPIService := services.NewNetworkAPIService(elrondProvider, cfg, true)
	networkAPIController := server.NewNetworkAPIController(networkAPIService, asserterServer)

	router := server.NewRouter(
		networkAPIController,
		accountAPIController,
		blockAPIController,
		mempoolAPIController,
		constructionAPIController,
	)

	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsRouter,
	}

	log.Info("elrond rosetta server is in offline mode")

	return httpServer, nil
}

func createServerOnline(
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

	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserterServer, err := asserter.NewServer(
		services.SupportedOperationTypes,
		false,
		[]*types.NetworkIdentifier{
			cfg.Network,
		},
		nil,
		false,
		"",
	)
	if err != nil {
		log.Error("cannot create asserter", "err", err)
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

	router := server.NewRouter(
		networkAPIController,
		accountAPIController,
		blockAPIController,
		constructionAPIController,
		mempoolAPIController,
	)

	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsRouter,
	}

	return httpServer, nil
}
