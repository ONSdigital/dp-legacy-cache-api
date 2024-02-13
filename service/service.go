package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/config"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	Router      *mux.Router
	API         *api.API
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
	mongoDB     DataStore
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")

	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	router := mux.NewRouter()
	router.Use(ensureJSONHeaderMiddleware)

	httpServer := serviceList.GetHTTPServer(cfg.BindAddr, router)

	mongoDB, err := serviceList.GetMongoDB(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise mongo DB", err)
		return nil, err
	}

	identityHandler := dphandlers.Identity(cfg.ZebedeeURL)

	legacyCacheAPI := api.Setup(ctx, cfg.IsPublishing, router, mongoDB, identityHandler)

	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, mongoDB); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	router.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the HTTP server in a new go-routine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in HTTP listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      router,
		API:         legacyCacheAPI,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      httpServer,
		mongoDB:     mongoDB,
	}, nil
}

func ensureJSONHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
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
