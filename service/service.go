package service

import (
	"context"

	dphandlers "github.com/ONSdigital/dp-net/handlers"

	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
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

	router.StrictSlash(true).Path("/health").HandlerFunc(svc.HealthCheck.Handler)
	svc.HealthCheck.Start(ctx)

	aliceChain := alice.New(dphandlers.Identity(svc.Config.ZebedeeURL)).Then(router)
	svc.Server = svc.ServiceList.GetHTTPServer(svc.Config.BindAddr, aliceChain)

	// Setup the API
	svc.API = api.Setup(ctx, router, svc.mongoDB)

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
