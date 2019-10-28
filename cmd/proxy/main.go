package main

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/logger"
	"github.com/ElrondNetwork/elrond-go/data/state/addressConverters"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/cache"
	"github.com/ElrondNetwork/elrond-proxy-go/testing"
	"github.com/pkg/profile"
	"github.com/urfave/cli"
)

var (
	log *logger.Logger

	proxyHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`

	// profileMode defines a flag for profiling the binary
	profileMode = cli.StringFlag{
		Name:  "profile-mode",
		Usage: "Profiling mode. Available options: cpu, mem, mutex, block",
		Value: "",
	}
	// configurationFile defines a flag for the path to the main toml configuration file
	configurationFile = cli.StringFlag{
		Name:  "config",
		Usage: "The main configuration file to load",
		Value: "./config/config.toml",
	}
	// configurationFile defines a flag for the path to the main toml configuration file
	economicsFile = cli.StringFlag{
		Name:  "economicsConfig",
		Usage: "The economics configuration file to load",
		Value: "./config/economics.toml",
	}
	// initialBalancesSkFile represents the path of the initialBalancesSk.pem file
	initialBalancesSkFile = cli.StringFlag{
		Name:  "pem-file",
		Usage: "This represents the path of the initialBalancesSk.pem file",
		Value: "./config/initialBalancesSk.pem",
	}
	// testHttpServerEn used to enable a test (mock) http server that will handle all requests
	testHttpServerEn = cli.BoolFlag{
		Name:  "test-http-server-enable",
		Usage: "Enables a test http server that will handle all requests",
	}

	testServer *testing.TestHttpServer
)

func main() {
	log = logger.DefaultLogger()
	log.SetLevel(logger.LogInfo)

	app := cli.NewApp()
	cli.AppHelpTemplate = proxyHelpTemplate
	app.Name = "Elrond Node Proxy CLI App"
	app.Version = "v1.0.0"
	app.Usage = "This is the entry point for starting a new Elrond node proxy"
	app.Flags = []cli.Flag{
		configurationFile,
		economicsFile,
		profileMode,
		initialBalancesSkFile,
		testHttpServerEn,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Action = func(c *cli.Context) error {
		return startProxy(c)
	}

	defer func() {
		if testServer != nil {
			testServer.Close()
		}
	}()

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startProxy(ctx *cli.Context) error {
	profileMode := ctx.GlobalString(profileMode.Name)
	switch profileMode {
	case "cpu":
		p := profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook)
		defer p.Stop()
	case "mem":
		p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
		defer p.Stop()
	case "mutex":
		p := profile.Start(profile.MutexProfile, profile.ProfilePath("."), profile.NoShutdownHook)
		defer p.Stop()
	case "block":
		p := profile.Start(profile.BlockProfile, profile.ProfilePath("."), profile.NoShutdownHook)
		defer p.Stop()
	}

	log.Info("Starting proxy...")

	configurationFileName := ctx.GlobalString(configurationFile.Name)
	generalConfig, err := loadMainConfig(configurationFileName, log)
	if err != nil {
		return err
	}

	economicsFileName := ctx.GlobalString(economicsFile.Name)
	economicsConfig, err := loadEconomicsConfig(economicsFileName, log)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Initialized with config from: %s", economicsFileName))

	stop := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	epf, err := createElrondProxyFacade(ctx, generalConfig, economicsConfig)
	if err != nil {
		return err
	}

	startWebServer(epf, generalConfig.GeneralSettings.ServerPort)

	go func() {
		<-sigs
		log.Info("terminating at user's signal...")
		stop <- true
	}()

	log.Info("Application is now running...")
	<-stop

	return nil
}

func loadMainConfig(filepath string, log *logger.Logger) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath, log)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadEconomicsConfig(filepath string, log *logger.Logger) (*config.EconomicsConfig, error) {
	cfg := &config.EconomicsConfig{}
	err := core.LoadTomlFile(cfg, filepath, log)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func createElrondProxyFacade(
	ctx *cli.Context,
	cfg *config.Config,
	ecCfg *config.EconomicsConfig,
) (*facade.ElrondProxyFacade, error) {

	var testHttpServerEnabled bool
	if ctx.IsSet(testHttpServerEn.Name) {
		testHttpServerEnabled = ctx.GlobalBool(testHttpServerEn.Name)
	}

	if testHttpServerEnabled {
		log.Info("Starting test HTTP server handling the requests...")
		testServer = testing.NewTestHttpServer()
		log.Info("Test HTTP server running at " + testServer.URL())

		testCfg := &config.Config{
			GeneralSettings: config.GeneralSettingsConfig{
				RequestTimeoutSec:                 10,
				HeartbeatCacheValidityDurationSec: 6000,
			},
			Observers: []*data.Observer{
				{
					ShardId: 0,
					Address: testServer.URL(),
				},
			},
		}
		testEcCfg := &config.EconomicsConfig{
			FeeSettings: config.FeeSettings{
				MinGasLimit: "1",
				MinGasPrice: "5",
			},
		}

		return createFacade(testCfg, testEcCfg, ctx.GlobalString(initialBalancesSkFile.Name))
	}

	return createFacade(cfg, ecCfg, ctx.GlobalString(initialBalancesSkFile.Name))
}

func createFacade(
	cfg *config.Config,
	ecConf *config.EconomicsConfig,
	pemFileLocation string,
) (*facade.ElrondProxyFacade, error) {
	addrConv, err := addressConverters.NewPlainAddressConverter(32, "")
	if err != nil {
		return nil, err
	}

	shardCoord, err := getShardCoordinator(cfg)
	if err != nil {
		return nil, err
	}

	bp, err := process.NewBaseProcessor(addrConv, cfg.GeneralSettings.RequestTimeoutSec, shardCoord)
	if err != nil {
		return nil, err
	}

	err = bp.ApplyConfig(cfg)
	if err != nil {
		return nil, err
	}

	accntProc, err := process.NewAccountProcessor(bp)
	if err != nil {
		return nil, err
	}

	privKeysLoader, err := faucet.NewPrivateKeysLoader(addrConv, shardCoord, pemFileLocation)
	if err != nil {
		return nil, err
	}

	faucetValue := big.NewInt(0)
	faucetValue.SetString(cfg.GeneralSettings.FaucetValue, 10)
	faucetProc, err := process.NewFaucetProcessor(ecConf, bp, privKeysLoader, faucetValue)
	if err != nil {
		return nil, err
	}

	txProc, err := process.NewTransactionProcessor(bp)
	if err != nil {
		return nil, err
	}

	gvpProc, err := process.NewVmValuesProcessor(bp)
	if err != nil {
		return nil, err
	}

	htbCacher := cache.NewHeartbeatMemoryCacher()
	cacheValidity := time.Duration(cfg.GeneralSettings.HeartbeatCacheValidityDurationSec) * time.Second

	htbProc, err := process.NewHeartbeatProcessor(bp, htbCacher, cacheValidity)
	if err != nil {
		return nil, err
	}
	htbProc.StartCacheUpdate()

	return facade.NewElrondProxyFacade(accntProc, txProc, gvpProc, htbProc, faucetProc)
}

func getShardCoordinator(cfg *config.Config) (sharding.Coordinator, error) {
	maxShardId := uint32(0)
	for _, observer := range cfg.Observers {
		shardId := observer.ShardId
		if maxShardId < shardId {
			maxShardId = shardId
		}
	}

	shardCoordinator, err := sharding.NewMultiShardCoordinator(maxShardId+1, 0)
	if err != nil {
		return nil, err
	}

	return shardCoordinator, nil
}

func startWebServer(proxyHandler api.ElrondProxyHandler, port int) {
	go func() {
		err := api.Start(proxyHandler, port)
		log.LogIfError(err)
	}()
}
