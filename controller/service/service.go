package service

import (
	"bytescheme/common/auth"
	"bytescheme/common/log"
	cservice "bytescheme/common/service"
	"bytescheme/common/util"
	gmodels "bytescheme/controller/generated/models"
	"bytescheme/controller/generated/restapi"
	"bytescheme/controller/generated/restapi/operations"
	"bytescheme/controller/generated/restapi/operations/controller"
	"bytescheme/controller/generated/restapi/operations/store"
	"bytescheme/controller/model"
	"bytescheme/controller/operation"
	"bytescheme/controller/shared"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

// Service is the service
type Service struct {
	ctx             context.Context
	cancel          context.CancelFunc
	registry        *operation.Registry
	server          *restapi.Server
	api             *operations.ControllerAPI
	spec            *loads.Document
	pathPermissions map[string]*auth.PathPermission
}

func getPathPermissionKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

func (service *Service) enrichContext(req *http.Request, principal interface{}) context.Context {
	ctx := req.Context()
	ctx = context.WithValue(ctx, util.XPrincipalKey, principal)
	return ctx
}

func (service *Service) authenticate(credential string) (interface{}, error) {
	for _, authenticator := range shared.Authenticators {
		principal, err := authenticator.Authenticate(credential)
		if err != nil {
			continue
		}
		if principal != nil {
			return principal, nil
		}
	}
	return nil, nil
}

func (service *Service) authorize(r *http.Request, iface interface{}) error {
	apiCtx := service.api.Context()
	// use context to lookup routes
	route, ok := apiCtx.LookupRoute(r)
	if !ok {
		log.Errorf("Route %s not found", r.RequestURI)
		return nil
	}
	pathPermKey := getPathPermissionKey(r.Method, route.PathPattern)
	pathPerm, ok := service.pathPermissions[pathPermKey]
	if !ok {
		log.Errorf("Permission not found for %s %s", r.Method, route.PathPattern)
		return nil
	}
	principal, ok := iface.(*auth.Principal)
	if !ok {
		log.Errorf("Principal type error for %+v", iface)
		err := fmt.Errorf("Unauthorized")
		return model.NewServiceError(401, err)
	}
	if principal.Roles != nil {
		for _, role := range principal.Roles {
			if pathPerm.Permission == role {
				return nil
			}
		}
	}
	log.Errorf("Expected permission %s, found %s for path %s", pathPerm.Permission, "Read", route.PathPattern)
	// Get from the role
	err := fmt.Errorf("Forbidden")
	return model.NewServiceError(403, err)
}

func (service *Service) setHandlers() error {
	// BearerTokenAuth registers a function that takes a token and returns a principal
	// it performs authentication based on an api key Authorization provided in the header
	// BearerTokenAuth func(string) (interface{}, error)
	api := service.api
	api.APIKeyAuth = service.authenticate
	api.APIAuthorizer = cservice.Authorizer(service.authorize)

	api.ControllerGetControllerHandler = controller.GetControllerHandlerFunc(func(params controller.GetControllerParams, principal interface{}) middleware.Responder {
		ctx := service.enrichContext(params.HTTPRequest, principal)
		cntlr, err := service.registry.GetController(ctx, params.ControllerID)
		if err != nil {
			serviceError, ok := err.(*model.ServiceError)
			if !ok {
				serviceError = &model.ServiceError{
					&gmodels.APIError{
						Status:  int32(500),
						Message: err.Error(),
					},
				}
			}
			return controller.NewGetControllerDefault(int(serviceError.Status)).WithPayload(serviceError.APIError)
		}
		return controller.NewGetControllerOK().WithPayload(cntlr)
	})
	api.ControllerListControllersHandler = controller.ListControllersHandlerFunc(func(params controller.ListControllersParams, principal interface{}) middleware.Responder {
		ctx := service.enrichContext(params.HTTPRequest, principal)
		cntlrs, err := service.registry.ListControllers(ctx)
		if err != nil {
			serviceError, ok := err.(*model.ServiceError)
			if !ok {
				serviceError = &model.ServiceError{
					&gmodels.APIError{
						Status:  int32(500),
						Message: err.Error(),
					},
				}
			}
			return controller.NewListControllersDefault(int(serviceError.Status)).WithPayload(serviceError.APIError)
		}
		return controller.NewListControllersOK().WithPayload(cntlrs)
	})
	api.ControllerUpdateControllerHandler = controller.UpdateControllerHandlerFunc(func(params controller.UpdateControllerParams, principal interface{}) middleware.Responder {
		ctx := service.enrichContext(params.HTTPRequest, principal)
		cntlr, err := service.registry.UpdateController(ctx, params.Payload)
		if err != nil {
			serviceError, ok := err.(*model.ServiceError)
			if !ok {
				serviceError = &model.ServiceError{
					&gmodels.APIError{
						Status:  int32(500),
						Message: err.Error(),
					},
				}
			}
			return controller.NewUpdateControllerDefault(int(serviceError.Status)).WithPayload(serviceError.APIError)
		}
		return controller.NewUpdateControllerOK().WithPayload(cntlr)
	})

	api.StoreListStoreKeysHandler = store.ListStoreKeysHandlerFunc(func(params store.ListStoreKeysParams, principal interface{}) middleware.Responder {
		ctx := service.enrichContext(params.HTTPRequest, principal)
		isPrefix := (params.Prefix != nil && *params.Prefix)
		keyValues, err := shared.ListStoreKeys(ctx, shared.Store, params.Key, isPrefix)
		if err != nil {
			serviceError, ok := err.(*model.ServiceError)
			if !ok {
				serviceError = &model.ServiceError{
					&gmodels.APIError{
						Status:  int32(500),
						Message: err.Error(),
					},
				}
			}
			return store.NewListStoreKeysDefault(int(serviceError.Status)).WithPayload(serviceError.APIError)
		}
		return store.NewListStoreKeysOK().WithPayload(keyValues)
	})

	api.StoreUpdateStoreKeysHandler = store.UpdateStoreKeysHandlerFunc(func(params store.UpdateStoreKeysParams, principal interface{}) middleware.Responder {
		ctx := service.enrichContext(params.HTTPRequest, principal)
		keyValues, err := shared.UpdateStoreKeys(ctx, shared.Store, params.Payload)
		if err != nil {
			serviceError, ok := err.(*model.ServiceError)
			if !ok {
				serviceError = &model.ServiceError{
					&gmodels.APIError{
						Status:  int32(500),
						Message: err.Error(),
					},
				}
			}
			return store.NewUpdateStoreKeysDefault(int(serviceError.Status)).WithPayload(serviceError.APIError)
		}
		return store.NewUpdateStoreKeysOK().WithPayload(keyValues)
	})

	api.StoreDeleteStoreKeysHandler = store.DeleteStoreKeysHandlerFunc(func(params store.DeleteStoreKeysParams, principal interface{}) middleware.Responder {
		ctx := service.enrichContext(params.HTTPRequest, principal)
		isPrefix := (params.Prefix != nil && *params.Prefix)
		keys, err := shared.DeleteStoreKeys(ctx, shared.Store, params.Key, isPrefix)
		if err != nil {
			serviceError, ok := err.(*model.ServiceError)
			if !ok {
				serviceError = &model.ServiceError{
					&gmodels.APIError{
						Status:  int32(500),
						Message: err.Error(),
					},
				}
			}
			return store.NewDeleteStoreKeysDefault(int(serviceError.Status)).WithPayload(serviceError.APIError)
		}
		return store.NewDeleteStoreKeysOK().WithPayload(keys)
	})

	pathPerms, err := cservice.GetPathPermissions(service.spec)
	if err != nil {
		return err
	}
	for idx := range pathPerms {
		pathPerm := pathPerms[idx]
		pathPermKey := getPathPermissionKey(pathPerm.Method, pathPerm.Path)
		service.pathPermissions[pathPermKey] = pathPerm
	}
	return nil
}

func (service *Service) startRedirecter() {
	host := service.server.Host
	port := service.server.Port
	go http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostRequested := host
		if strings.HasPrefix(r.Host, "bytescheme.") {
			hostRequested = "controller.bytescheme.com"
		}
		log.Infof("Requested host: %s, resolved host: %s", host, hostRequested)
		http.Redirect(w, r, fmt.Sprintf("https://%s%s", hostRequested, r.RequestURI), http.StatusMovedPermanently)
	}))
}

// NewService creates the http service
func NewService(host string, port int, registry *operation.Registry) (*Service, error) {
	certDir := os.Getenv("HTTPS_CERT_DIR")
	if certDir == "" {
		certDir, _ = os.Getwd()
		certDir = certDir + "/../controller/cert"
	}
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	api := operations.NewControllerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.Host = host
	server.TLSHost = host
	server.Port = port
	server.TLSPort = port + 3
	server.TLSCertificate = flags.Filename(certDir + "/controller.bytescheme.com.cer")
	server.TLSCertificateKey = flags.Filename(certDir + "/controller.bytescheme.com.key")
	svc := &Service{
		registry:        registry,
		server:          server,
		api:             api,
		spec:            swaggerSpec,
		pathPermissions: map[string]*auth.PathPermission{},
	}
	err = svc.setHandlers()
	server.ConfigureAPI()
	if err != nil {
		return nil, err
	}
	return svc, nil
}

// Serve is a blocking call to start the server
func (service *Service) Serve() error {
	service.ctx, service.cancel = context.WithCancel(context.Background())
	util.ShutdownHandler.RegisterCloseable(service)
	service.startRedirecter()
	return service.server.Serve()
}

// Close closes the service
func (service *Service) Close() error {
	service.cancel()
	return nil
}

// IsClosed returns true if the service is already closed
func (service *Service) IsClosed() bool {
	select {
	case <-service.ctx.Done():
		return true
	default:
		return false
	}
}
