package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

type validatorInput struct {
	Name      string
	Validator validator.Func
}

// CreateServer creates a HTTP server
func CreateServer(versionManager data.VersionManagerHandler, port int) (*http.Server, error) {
	ws := gin.Default()
	ws.Use(cors.Default())

	err := registerValidators()
	if err != nil {
		return nil, err
	}

	err = registerRoutes(ws, versionManager)
	if err != nil {
		return nil, err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ws,
	}

	return httpServer, nil
}

func registerRoutes(ws *gin.Engine, versionManager data.VersionManagerHandler) error {
	versionsMap, err := versionManager.GetAllVersions()
	if err != nil {
		return err
	}

	for version, versionData := range versionsMap {
		versionGroup := ws.Group(version)

		for path, group := range versionData.ApiHandler.GetAllGroups() {
			subGroup := versionGroup.Group(path)
			subGroup.Use(WithElrondProxyFacade(versionData.Facade, version))
			group.Routes(subGroup)
		}
	}

	return nil
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

// WithElrondProxyFacade middleware will set up an ElrondFacade object in the gin context
func WithElrondProxyFacade(elrondProxyFacade ElrondProxyHandler, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(version, elrondProxyFacade)
		c.Next()
	}
}
