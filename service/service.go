package service

import (
	"context"
	"net/http"

	clientsidentity "github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/config"

	// dphandlers "github.com/ONSdigital/dp-net/handlers"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config         *config.Config
	Server         HTTPServer
	Router         *mux.Router
	API            *api.API
	IdentityClient *clientsidentity.Client
	ServiceList    *ExternalServiceList
	HealthCheck    HealthChecker
	mongoDB        DataStore
}

// New creates a new service instance
func New(cfg *config.Config, serviceList *ExternalServiceList) *Service {
	svc := &Service{
		Config:      cfg,
		ServiceList: serviceList,
	}
	return svc
}

// Run the service
func (svc *Service) Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) error {
	var err error
	log.Info(ctx, "running service")
	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	// Get HTTP Server and ... // TODO: Add any middleware that your service requires
	router := mux.NewRouter()

	svc.Server = serviceList.GetHTTPServer(cfg.BindAddr, router)

	svc.mongoDB, err = serviceList.GetMongoDB(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise mongo DB", err)
		return err
	}

	svc.HealthCheck, err = serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)

	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return err
	}

	if err = registerCheckers(ctx, svc.HealthCheck, svc.mongoDB); err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", errors.Wrap(err, "unable to register checkers"))
		return err
	}

	// Get Identity Client (only if private endpoints are enabled)
	//if svc.Config.EnablePrivateEndpoints {
	svc.IdentityClient = clientsidentity.New(svc.Config.ZebedeeURL)
	//}

	m := svc.createMiddleware(svc.Config)
	svc.Server = svc.ServiceList.GetHTTPServer(svc.Config.BindAddr, m.Then(router))

	// Create Dataset API
	// urlBuilder := url.NewBuilder(svc.Config.WebsiteURL)
	datasetPermissions, permissions := getAuthorisationHandlers(ctx, svc.Config)
	// svc.API = api.Setup(ctx, svc.config, r, ds, urlBuilder, downloadGenerators, datasetPermissions, permissions)

	// Setup the API
	svc.API = api.Setup(ctx, router, svc.mongoDB, datasetPermissions, permissions)

	router.StrictSlash(true).Path("/health").HandlerFunc(svc.HealthCheck.Handler)
	svc.HealthCheck.Start(ctx)

	// Run the HTTP server in a new go-routine
	go func() {
		if err = svc.Server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in HTTP listen and serve")
		}
	}()

	return nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		if svc.mongoDB != nil {
			if err := svc.mongoDB.Close(ctx); err != nil {
				log.Error(ctx, "failed to close MongoDB connection", err)
				hasShutdownError = true
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	healthChecker HealthChecker,
	dataStore DataStore,
) (err error) {
	hasErrors := false

	if err = healthChecker.AddCheck("Mongo DB", dataStore.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for mongo db", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}

// CreateMiddleware creates an Alice middleware chain of handlers
// to forward collectionID from cookie from header
func (svc *Service) createMiddleware(cfg *config.Config) alice.Chain {
	// healthcheck
	healthcheckHandler := newMiddleware(svc.HealthCheck.Handler, "/health")
	middleware := alice.New(healthcheckHandler)

	// Only add the identity middleware when running in publishing.
	if cfg.EnablePrivateEndpoints {
		middleware = middleware.Append(dphandlers.IdentityWithHTTPClient(svc.IdentityClient))
	}

	// collection ID
	middleware = middleware.Append(dphandlers.CheckHeader(dphandlers.CollectionID))

	return middleware
}

// newMiddleware creates a new http.Handler to intercept /health requests.
func newMiddleware(healthcheckHandler func(http.ResponseWriter, *http.Request), path string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" && req.URL.Path == path {
				healthcheckHandler(w, req)
				return
			} else if req.Method == "GET" && req.URL.Path == "/healthcheck" {
				http.NotFound(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}

func getAuthorisationHandlers(ctx context.Context, cfg *config.Config) (datasetPermissions, permissions api.AuthHandler) {
	// if !cfg.EnablePermissionsAuth {
	// 	log.Info(ctx, "feature flag not enabled defaulting to nop auth impl", log.Data{"feature": "ENABLE_PERMISSIONS_AUTH"})
	// 	return &auth.NopHandler{}, &auth.NopHandler{}
	// }

	log.Info(ctx, "feature flag enabled", log.Data{"feature": "ENABLE_PERMISSIONS_AUTH"})

	authClient := auth.NewPermissionsClient(dphttp.NewClient())
	authVerifier := auth.DefaultPermissionsVerifier()

	// for checking caller permissions when we have a datasetID, collection ID and user/service token
	datasetPermissions = auth.NewHandler(
		auth.NewDatasetPermissionsRequestBuilder(cfg.ZebedeeURL, "dataset_id", mux.Vars),
		authClient,
		authVerifier,
	)

	// for checking caller permissions when we only have a user/service token
	permissions = auth.NewHandler(
		auth.NewPermissionsRequestBuilder(cfg.ZebedeeURL),
		authClient,
		authVerifier,
	)

	return datasetPermissions, permissions
}
