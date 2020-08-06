package rosetta

import (
	"fmt"
	"net/http"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

var log = logger.GetOrCreate("rosetta")

// Start will boot up the api and appropriate routes, handlers and validators
func StartRosetta(elrondFacade interface{}, port int) error {
	network := &types.NetworkIdentifier{
		Blockchain: "Elrond",
		Network:    "Testnet",
	}

	// The asserter automatically rejects incorrectly formatted
	// requests.
	asserter, err := asserter.NewServer(
		[]string{"Transfer", "Reward"},
		false,
		[]*types.NetworkIdentifier{network},
	)
	if err != nil {
		log.Error("cannot create asserter", "err", err)
	}

	networkAPIService := NewNetworkAPIService(network)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserter,
	)

	blockAPIService := NewBlockAPIService(network)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		asserter,
	)

	router := server.NewRouter(networkAPIController, blockAPIController)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	log.Info("Listening on port", "port", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), corsRouter)
}
