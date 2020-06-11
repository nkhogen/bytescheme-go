// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import(
		"fmt"
		"strings"
		_ "bytescheme/controller/generated/statik"
		"github.com/rakyll/statik/fs"
		
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"bytescheme/controller/generated/restapi/operations"
	"bytescheme/controller/generated/restapi/operations/controller"
	"bytescheme/controller/generated/restapi/operations/store"
)

//go:generate swagger generate server --target ../../generated --name Controller --spec ../../controller-swagger.json

func configureFlags(api *operations.ControllerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ControllerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization" header is set
	if api.APIKeyAuth == nil {
		api.APIKeyAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (ApiKey) Authorization from header param [Authorization] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	if api.StoreDeleteStoreKeysHandler == nil {
		api.StoreDeleteStoreKeysHandler = store.DeleteStoreKeysHandlerFunc(func(params store.DeleteStoreKeysParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation store.DeleteStoreKeys has not yet been implemented")
		})
	}
	if api.ControllerGetControllerHandler == nil {
		api.ControllerGetControllerHandler = controller.GetControllerHandlerFunc(func(params controller.GetControllerParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation controller.GetController has not yet been implemented")
		})
	}
	if api.ControllerListControllersHandler == nil {
		api.ControllerListControllersHandler = controller.ListControllersHandlerFunc(func(params controller.ListControllersParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation controller.ListControllers has not yet been implemented")
		})
	}
	if api.StoreListStoreKeysHandler == nil {
		api.StoreListStoreKeysHandler = store.ListStoreKeysHandlerFunc(func(params store.ListStoreKeysParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation store.ListStoreKeys has not yet been implemented")
		})
	}
	if api.ControllerUpdateControllerHandler == nil {
		api.ControllerUpdateControllerHandler = controller.UpdateControllerHandlerFunc(func(params controller.UpdateControllerParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation controller.UpdateController has not yet been implemented")
		})
	}
	if api.StoreUpdateStoreKeysHandler == nil {
		api.StoreUpdateStoreKeysHandler = store.UpdateStoreKeysHandlerFunc(func(params store.UpdateStoreKeysParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation store.UpdateStoreKeys has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request received for %s\n", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/v1") {
			handler.ServeHTTP(w, r)
		} else {
			statikFS, err := fs.New()
			if err != nil {
				fmt.Printf("Cannot create statik FS. Error: %s\n", err.Error())
				statikFS = http.Dir(".")
			}
			file, err := statikFS.Open(r.URL.Path)
			if err == nil {
				file.Close()
			} else {
				r.URL.Path="index.html"
			}
			http.FileServer(statikFS).ServeHTTP(w, r)
		}
	})
}
