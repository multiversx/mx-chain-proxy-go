package rosetta

import (
	"fmt"
	"net/http"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/services"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

var log = logger.GetOrCreate("rosetta")

// CreateServer creates a HTTP server
func CreateServer(elrondFacade api.ElrondProxyHandler, port int) (*http.Server, error) {
	elrondClient := client.NewElrondClient(elrondFacade)

	networkConfig, err := elrondClient.GetNetworkConfig()
	if err != nil {
		return nil, err
	}

	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserterServer, err := asserter.NewServer(
		services.SupportedOperationTypes,
		false,
		[]*types.NetworkIdentifier{
			{
				Blockchain: services.ElrondBlockchainName,
				Network:    networkConfig.ChainID,
			},
		},
	)
	if err != nil {
		log.Error("cannot create asserter", "err", err)
	}

	// Create network service
	networkAPIService := services.NewNetworkAPIService(elrondClient)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserterServer,
	)

	// Create account service
	accountAPIService := services.NewAccountAPIService(elrondClient)
	accountAPIController := server.NewAccountAPIController(
		accountAPIService,
		asserterServer,
	)

	// Create block service
	blockAPIService := services.NewBlockAPIService(elrondClient)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		asserterServer,
	)

	// Create construction service
	constructionAPIService := services.NewConstructionAPIService(elrondClient)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserterServer,
	)

	router := server.NewRouter(networkAPIController, accountAPIController, blockAPIController, constructionAPIController)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsRouter,
	}

	return httpServer, nil
}
