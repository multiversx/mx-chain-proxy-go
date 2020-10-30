package main

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	erdConfig "github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state/factory"
	hasherFactory "github.com/ElrondNetwork/elrond-go/hashing/factory"
	marshalFactory "github.com/ElrondNetwork/elrond-go/marshal/factory"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet"
	"github.com/ElrondNetwork/elrond-proxy-go/observer"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/cache"
	"github.com/ElrondNetwork/elrond-proxy-go/process/database"
	processFactory "github.com/ElrondNetwork/elrond-proxy-go/process/factory"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta"
	"github.com/ElrondNetwork/elrond-proxy-go/testing"
	versionsFactory "github.com/ElrondNetwork/elrond-proxy-go/versions/factory"
	"github.com/pkg/profile"
	"github.com/urfave/cli"
)

var (
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

	log = logger.GetOrCreate("proxy")

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
	// economicsFile defines a flag for the path to the economics toml configuration file
	economicsFile = cli.StringFlag{
		Name:  "economics-config",
		Usage: "The economics configuration file to load",
		Value: "./config/economics.toml",
	}
	// walletKeyPemFile represents the path of the wallet (address) pem file
	walletKeyPemFile = cli.StringFlag{
		Name:  "pem-file",
		Usage: "This represents the path of the walletKey.pem file",
		Value: "./config/walletKey.pem",
	}
	// externalConfigFile defines a flag for the path to the external toml configuration file
	externalConfigFile = cli.StringFlag{
		Name: "config-external",
		Usage: "The path for the external configuration file. This TOML file contains" +
			" external configurations such as ElasticSearch's URL and login information",
		Value: "./config/external.toml",
	}

	// testHttpServerEn used to enable a test (mock) http server that will handle all requests
	testHttpServerEn = cli.BoolFlag{
		Name:  "test-mode",
		Usage: "Enables a test http server that will handle all requests",
	}

	startAsRosetta = cli.BoolFlag{
		Name:  "rosetta",
		Usage: "Starts the proxy as a rosetta server",
	}

	testServer *testing.TestHttpServer
)

func main() {
	log.SetLevel(logger.LogInfo)
	removeLogColors()

	app := cli.NewApp()
	cli.AppHelpTemplate = proxyHelpTemplate
	app.Name = "Elrond Node Proxy CLI App"
	app.Version = "v1.0.0"
	app.Usage = "This is the entry point for starting a new Elrond node proxy"
	app.Flags = []cli.Flag{
		configurationFile,
		economicsFile,
		externalConfigFile,
		profileMode,
		walletKeyPemFile,
		testHttpServerEn,
		startAsRosetta,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Action = startProxy

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
	generalConfig, err := loadMainConfig(configurationFileName)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Initialized with main config from: %s", configurationFile))

	economicsFileName := ctx.GlobalString(economicsFile.Name)
	economicsConfig, err := loadEconomicsConfig(economicsFileName)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Initialized with economics config from: %s", economicsFileName))

	externalConfigurationFileName := ctx.GlobalString(externalConfigFile.Name)
	externalConfig, err := loadExternalConfig(externalConfigurationFileName)
	if err != nil {
		return err
	}

	versionManager, err := createVersionManagerTestOrProduction(ctx, generalConfig, economicsConfig, externalConfig)
	if err != nil {
		return err
	}

	httpServer, err := startWebServer(versionManager, ctx, generalConfig)
	if err != nil {
		return err
	}

	waitForServerShutdown(httpServer)
	return nil
}

func loadMainConfig(filepath string) (*config.Config, error) {
	cfg := &config.Config{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadEconomicsConfig(filepath string) (*erdConfig.EconomicsConfig, error) {
	cfg := &erdConfig.EconomicsConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadExternalConfig(filepath string) (*erdConfig.ExternalConfig, error) {
	cfg := &erdConfig.ExternalConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func createVersionManagerTestOrProduction(
	ctx *cli.Context,
	cfg *config.Config,
	ecCfg *erdConfig.EconomicsConfig,
	exCfg *erdConfig.ExternalConfig,
) (data.VersionManagerHandler, error) {

	var testHTTPServerEnabled bool
	if ctx.IsSet(testHttpServerEn.Name) {
		testHTTPServerEnabled = ctx.GlobalBool(testHttpServerEn.Name)
	}

	if testHTTPServerEnabled {
		log.Info("Starting test HTTP server handling the requests...")
		testServer = testing.NewTestHttpServer()
		log.Info("Test HTTP server running at " + testServer.URL())

		testCfg := &config.Config{
			GeneralSettings: config.GeneralSettingsConfig{
				RequestTimeoutSec:                 10,
				HeartbeatCacheValidityDurationSec: 60,
				ValStatsCacheValidityDurationSec:  60,
				FaucetValue:                       "10000000000",
			},
			Observers: []*data.NodeData{
				{
					ShardId: 0,
					Address: testServer.URL(),
				},
				{
					ShardId: 1,
					Address: testServer.URL(),
				},
				{
					ShardId: core.MetachainShardId,
					Address: testServer.URL(),
				},
			},
			FullHistoryNodes: []*data.NodeData{
				{
					ShardId: 0,
					Address: testServer.URL(),
				},
				{
					ShardId: 1,
					Address: testServer.URL(),
				},
				{
					ShardId: core.MetachainShardId,
					Address: testServer.URL(),
				},
			},
			AddressPubkeyConverter: cfg.AddressPubkeyConverter,
		}

		return createVersionManager(testCfg, ecCfg, exCfg, ctx.GlobalString(walletKeyPemFile.Name), false)
	}

	isRosettaOn := ctx.GlobalBool(startAsRosetta.Name)
	return createVersionManager(cfg, ecCfg, exCfg, ctx.GlobalString(walletKeyPemFile.Name), isRosettaOn)
}

func createVersionManager(
	cfg *config.Config,
	ecConf *erdConfig.EconomicsConfig,
	exCfg *erdConfig.ExternalConfig,
	pemFileLocation string,
	isRosettaOn bool,
) (data.VersionManagerHandler, error) {
	pubKeyConverter, err := factory.NewPubkeyConverter(cfg.AddressPubkeyConverter)
	if err != nil {
		return nil, err
	}

	marshalizer, err := marshalFactory.NewMarshalizer(cfg.Marshalizer.Type)
	if err != nil {
		return nil, err
	}
	hasher, err := hasherFactory.NewHasher(cfg.Hasher.Type)
	if err != nil {
		return nil, err
	}

	shardCoord, err := getShardCoordinator(cfg)
	if err != nil {
		return nil, err
	}

	nodesProviderFactory, err := observer.NewNodesProviderFactory(*cfg)
	if err != nil {
		return nil, err
	}

	observersProvider, err := nodesProviderFactory.CreateObservers()
	if err != nil {
		return nil, err
	}

	fullHistoryNodesProvider, err := nodesProviderFactory.CreateFullHistoryNodes()
	if err != nil {
		if err != observer.ErrEmptyObserversList {
			return nil, err
		}
	}

	bp, err := process.NewBaseProcessor(
		cfg.GeneralSettings.RequestTimeoutSec,
		shardCoord,
		observersProvider,
		fullHistoryNodesProvider,
		pubKeyConverter,
	)
	if err != nil {
		return nil, err
	}

	connector, err := createElasticSearchConnector(exCfg)
	if err != nil {
		return nil, err
	}

	accntProc, err := process.NewAccountProcessor(bp, pubKeyConverter, connector)
	if err != nil {
		return nil, err
	}

	privKeysLoader, err := faucet.NewPrivateKeysLoader(shardCoord, pemFileLocation, pubKeyConverter)
	if err != nil {
		return nil, err
	}

	faucetValue := big.NewInt(0)
	faucetValue.SetString(cfg.GeneralSettings.FaucetValue, 10)
	faucetProc, err := processFactory.CreateFaucetProcessor(ecConf, bp, privKeysLoader, faucetValue, pubKeyConverter)
	if err != nil {
		return nil, err
	}

	txProc, err := process.NewTransactionProcessor(bp, pubKeyConverter, hasher, marshalizer)
	if err != nil {
		return nil, err
	}

	scQueryProc, err := process.NewSCQueryProcessor(bp, pubKeyConverter)
	if err != nil {
		return nil, err
	}

	htbCacher := cache.NewHeartbeatMemoryCacher()
	cacheValidity := time.Duration(cfg.GeneralSettings.HeartbeatCacheValidityDurationSec) * time.Second

	htbProc, err := process.NewHeartbeatProcessor(bp, htbCacher, cacheValidity)
	if err != nil {
		return nil, err
	}
	if !isRosettaOn {
		htbProc.StartCacheUpdate()
	}

	valStatsCacher := cache.NewValidatorsStatsMemoryCacher()
	cacheValidity = time.Duration(cfg.GeneralSettings.ValStatsCacheValidityDurationSec) * time.Second

	valStatsProc, err := process.NewValidatorStatisticsProcessor(bp, valStatsCacher, cacheValidity)
	if err != nil {
		return nil, err
	}
	if !isRosettaOn {
		valStatsProc.StartCacheUpdate()
	}

	nodeStatusProc, err := process.NewNodeStatusProcessor(bp)
	if err != nil {
		return nil, err
	}

	blockProc, err := process.NewBlockProcessor(connector, bp)
	if err != nil {
		return nil, err
	}

	commonApiHandler := api.NewCommonApiHandler()

	facadeArgs := versionsFactory.FacadeArgs{
		AccountProcessor:             accntProc,
		FaucetProcessor:              faucetProc,
		BlockProcessor:               blockProc,
		HeartbeatProcessor:           htbProc,
		NodeStatusProcessor:          nodeStatusProc,
		ScQueryProcessor:             scQueryProc,
		TransactionProcessor:         txProc,
		ValidatorStatisticsProcessor: valStatsProc,
	}

	return versionsFactory.CreateVersionManager(facadeArgs, commonApiHandler)
}

func createElasticSearchConnector(exCfg *erdConfig.ExternalConfig) (process.ExternalStorageConnector, error) {
	if !exCfg.ElasticSearchConnector.Enabled {
		return database.NewDisabledElasticSearchConnector(), nil
	}

	return database.NewElasticSearchConnector(
		exCfg.ElasticSearchConnector.URL,
		exCfg.ElasticSearchConnector.Username,
		exCfg.ElasticSearchConnector.Password,
	)
}

func getShardCoordinator(cfg *config.Config) (sharding.Coordinator, error) {
	maxShardID := uint32(0)
	for _, obs := range cfg.Observers {
		shardID := obs.ShardId
		isMetaChain := shardID == core.MetachainShardId
		if maxShardID < shardID && !isMetaChain {
			maxShardID = shardID
		}
	}

	shardCoordinator, err := sharding.NewMultiShardCoordinator(maxShardID+1, 0)
	if err != nil {
		return nil, err
	}

	return shardCoordinator, nil
}

func startWebServer(versionManager data.VersionManagerHandler, cliContext *cli.Context, generalConfig *config.Config) (*http.Server, error) {
	var err error
	var httpServer *http.Server

	port := generalConfig.GeneralSettings.ServerPort
	asRosetta := cliContext.GlobalBool(startAsRosetta.Name)
	if asRosetta {
		httpServer, err = rosetta.CreateServer(proxyHandler, generalConfig, port)
	} else {
		httpServer, err = api.CreateServer(proxyHandler, port)
	}
	if err != nil {
		return nil, err
	}
	go func() {
		err = httpServer.ListenAndServe()
		if err != nil {
			log.Error("cannot ListenAndServe()", "err", err)
			os.Exit(1)
		}
	}()

	return httpServer, nil
}

func waitForServerShutdown(httpServer *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	shutdownContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = httpServer.Shutdown(shutdownContext)
	_ = httpServer.Close()
}

func removeLogColors() {
	err := logger.RemoveLogObserver(os.Stdout)
	if err != nil {
		panic("error removing log observer: " + err.Error())
	}

	err = logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
	if err != nil {
		panic("error setting log observer: " + err.Error())
	}
}
