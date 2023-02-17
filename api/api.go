package api

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/hashing/factory"
	"github.com/multiversx/mx-chain-core-go/hashing/sha256"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/api/middleware"
	"github.com/multiversx/mx-chain-proxy-go/config"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"gopkg.in/go-playground/validator.v8"
)

var log = logger.GetOrCreate("api")

type validatorInput struct {
	Name      string
	Validator validator.Func
}

// CreateServer creates a HTTP server
func CreateServer(
	versionsRegistry data.VersionsRegistryHandler,
	port int,
	apiLoggingConfig config.ApiLoggingConfig,
	credentialsConfig config.CredentialsConfig,
	statusMetricsExtractor middleware.StatusMetricsExtractor,
	rateLimitTimeWindowInSeconds int,
	isProfileModeActivated bool,
	shouldStartSwaggerUI bool,
) (*http.Server, error) {
	ws := gin.Default()
	ws.Use(cors.Default())

	err := registerValidators()
	if err != nil {
		return nil, err
	}

	err = registerRoutes(ws, versionsRegistry, apiLoggingConfig, credentialsConfig, statusMetricsExtractor, rateLimitTimeWindowInSeconds, isProfileModeActivated, shouldStartSwaggerUI)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ws,
	}

	return httpServer, nil
}

func registerValidators() error {
	validators := []validatorInput{
		{Name: "skValidator", Validator: skValidator},
	}
	for _, validatorFunc := range validators {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation(validatorFunc.Name, validatorFunc.Validator)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func registerRoutes(
	ws *gin.Engine,
	versionsRegistry data.VersionsRegistryHandler,
	apiLoggingConfig config.ApiLoggingConfig,
	credentialsConfig config.CredentialsConfig,
	statusMetricsExtractor middleware.StatusMetricsExtractor,
	rateLimitTimeWindowInSeconds int,
	isProfileModeActivated bool,
	shouldStartSwaggerUI bool,
) error {
	versionsMap, err := versionsRegistry.GetAllVersions()
	if err != nil {
		return err
	}

	if shouldStartSwaggerUI {
		ws.Use(static.ServeRoot("/", "config/swagger"))
	}

	if apiLoggingConfig.LoggingEnabled {
		responseLoggerMiddleware := middleware.NewResponseLoggerMiddleware(time.Duration(apiLoggingConfig.ThresholdInMicroSeconds) * time.Microsecond)
		ws.Use(responseLoggerMiddleware.MiddlewareHandlerFunc())
	}

	// TODO: maybe add a flag when starting proxy if metrics should be exposed or not
	metricsMiddleware, err := middleware.NewMetricsMiddleware(statusMetricsExtractor)
	if err != nil {
		return err
	}

	for version, versionData := range versionsMap {
		limitsMap := getLimitsMapForVersion(versionData)
		rateLimitTimeWindowDuration := time.Duration(rateLimitTimeWindowInSeconds) * time.Second
		rateLimiter, err := middleware.NewRateLimiter(limitsMap, rateLimitTimeWindowDuration)
		if err != nil {
			return err
		}
		startRateLimiterReset(rateLimitTimeWindowInSeconds, rateLimiter, version)
		versionGroup := ws.Group(version)
		for path, group := range versionData.ApiHandler.GetAllGroups() {
			subGroup := versionGroup.Group(path)
			group.RegisterRoutes(
				subGroup,
				versionData.ApiConfig,
				getAuthenticationFunc(credentialsConfig),
				rateLimiter.MiddlewareHandlerFunc(),
				metricsMiddleware.MiddlewareHandlerFunc(),
			)
		}
	}

	if isProfileModeActivated {
		pprof.Register(ws)
	}

	return nil
}

func getAuthenticationFunc(credentialsConfig config.CredentialsConfig) gin.HandlerFunc {
	if len(credentialsConfig.Credentials) == 0 {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				data.GenericAPIResponse{
					Data:  nil,
					Error: "no credentials found on server",
					Code:  data.ReturnCodeInternalError,
				},
			)
		}
	}

	var hasher hashing.Hasher
	var err error
	hasher, err = factory.NewHasher(credentialsConfig.Hasher.Type)
	if err != nil {
		log.Warn("cannot create hasher from config. Will use Sha256 as default", "error", err)
		hasher = sha256.NewSha256() // fallback in case the hasher creation failed
	}

	accounts := gin.Accounts{}
	for _, pair := range credentialsConfig.Credentials {
		accounts[pair.Username] = pair.Password
	}

	authenticationFunction := func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, data.GenericAPIResponse{
				Data:  nil,
				Error: "this endpoint requires Basic Authentication",
				Code:  data.ReturnCodeRequestError,
			})
			return
		}

		userPassword, ok := accounts[user]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, data.GenericAPIResponse{
				Data:  nil,
				Error: "username does not exist",
				Code:  data.ReturnCodeRequestError,
			})
			return
		}

		if userPassword != hex.EncodeToString(hasher.Compute(pass)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, data.GenericAPIResponse{
				Data:  nil,
				Error: "invalid password",
				Code:  data.ReturnCodeRequestError,
			})
			return
		}
	}

	return authenticationFunction
}

func getLimitsMapForVersion(versionData *data.VersionData) map[string]uint64 {
	limitsMap := make(map[string]uint64)
	for packageName, packageConfig := range versionData.ApiConfig.APIPackages {
		for _, routeConfig := range packageConfig.Routes {
			if routeConfig.RateLimit > 0 {
				mapKey := fmt.Sprintf("/%s%s", packageName, routeConfig.Name)
				limitsMap[mapKey] = routeConfig.RateLimit
			}
		}
	}

	return limitsMap
}

func startRateLimiterReset(rateLimiterDuration int, rl middleware.RateLimiterHandler, version string) {
	go func() {
		for {
			time.Sleep(time.Duration(rateLimiterDuration) * time.Second)
			rl.ResetMap(version)
		}
	}()
}

// skValidator validates a secret key from user input for correctness
func skValidator(
	_ *validator.Validate,
	_ reflect.Value,
	_ reflect.Value,
	_ reflect.Value,
	_ reflect.Type,
	_ reflect.Kind,
	_ string,
) bool {
	return true
}
