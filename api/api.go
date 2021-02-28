package api

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/hashing/factory"
	"github.com/ElrondNetwork/elrond-go/hashing/sha256"
	"github.com/ElrondNetwork/elrond-proxy-go/api/middleware"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
) (*http.Server, error) {
	ws := gin.Default()
	ws.Use(cors.Default())

	err := registerValidators()
	if err != nil {
		return nil, err
	}

	err = registerRoutes(ws, versionsRegistry, apiLoggingConfig, credentialsConfig)
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
) error {
	versionsMap, err := versionsRegistry.GetAllVersions()
	if err != nil {
		return err
	}

	if apiLoggingConfig.LoggingEnabled {
		responseLoggerMiddleware := middleware.NewResponseLoggerMiddleware(time.Duration(apiLoggingConfig.ThresholdInMicroSeconds) * time.Microsecond)
		ws.Use(responseLoggerMiddleware.MiddlewareHandlerFunc())
	}

	for version, versionData := range versionsMap {
		versionGroup := ws.Group(version)
		for path, group := range versionData.ApiHandler.GetAllGroups() {
			subGroup := versionGroup.Group(path)
			group.RegisterRoutes(subGroup, versionData.ApiConfig, getAuthenticationFunc(credentialsConfig))
		}
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
		hasher = sha256.Sha256{} // fallback in case the hasher creation failed
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
