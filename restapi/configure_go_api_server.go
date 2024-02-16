// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	awesome "github.com/ifaisalalam/Go-awesome-service"
	"github.com/ifaisalalam/Go-awesome-service/handlers"
	"github.com/ifaisalalam/Go-awesome-service/middleware/auth"
	"github.com/ifaisalalam/Go-awesome-service/restapi/operations"
)

//go:generate swagger generate server --target ../../APIServer --name GoAPIServer --spec ../swagger/swagger.yml --principal interface{}

const (
	// EnvNameBeta ...
	EnvNameBeta string = "BETA"
	// EnvNameGamma ...
	EnvNameGamma = "GAMMA"
	// EnvNameProd ...
	EnvNameProd = "PROD"
)

// SupportedEnv contains list of supported Env name values.
var SupportedEnv = []string{EnvNameBeta, EnvNameGamma, EnvNameProd}

// EnvName contains the current Env name.
var EnvName = func() string {
	if env, ok := os.LookupEnv("ENV"); ok {
		for _, name := range SupportedEnv {
			if name == env {
				return name
			}
		}
		panic("unsupported app env")
	}
	return EnvNameBeta
}()

func configureFlags(api *operations.GoAPIServerAPI) {
	//api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{}
}

func configureAPI(api *operations.GoAPIServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	if EnvName == EnvNameBeta {
		api.UseSwaggerUI()
	}
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	registerHandlers(api)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func registerHandlers(api *operations.GoAPIServerAPI) {
	svc := awesome.NewService(context.Background(), EnvName)

	api.LivenessHandler = handlers.NewLivenessHandler()
	api.ReadinessHandler = handlers.NewReadinessHandler()

	shortURLHandler := handlers.NewShortURLHandler(svc.Shortener)
	api.CreateShortURLHandler = shortURLHandler.CreateShortURLHandler()
	api.GetTargetURLHandler = shortURLHandler.GetTargetURLHandler()
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return auth.NewNoOpAuthMiddleware(handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
