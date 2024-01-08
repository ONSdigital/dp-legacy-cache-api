package service

import (
	"context"

	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"github.com/ONSdigital/dp-legacy-cache-api/mongo"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config        *config.Config
	Server        HTTPServer
	Router        *mux.Router
	API           *api.API
	ServiceList   *ExternalServiceList
	HealthCheck   HealthChecker
	MongoDBClient mongo.MongoDBClient
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
    log.Info(ctx, "running service")

    log.Info(ctx, "using service configuration", log.Data{"config": cfg})

    // Get HTTP Server and ... // TODO: Add any middleware that your service requires
    r := mux.NewRouter()

    s := serviceList.GetHTTPServer(cfg.BindAddr, r)
	mongoClient := &mongo.Mongo{
		MongoConfig: cfg.MongoConfig,
	}
	err := mongoClient.Init(ctx)
 
    if err != nil {
        log.Fatal(ctx, "failed to initialize MongoDB", err)
        return nil, err
    }

    // Setup the API
    a := api.Setup(ctx, r, mongoClient)

    hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
    if err != nil {
        log.Fatal(ctx, "could not instantiate healthcheck", err)
        return nil, err
    }

    if err := registerCheckers(ctx, hc); err != nil {
        return nil, errors.Wrap(err, "unable to register checkers")
    }

    r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
    hc.Start(ctx)

    // Run the HTTP server in a new go-routine
    go func() {
        if err := s.ListenAndServe(); err != nil {
            svcErrors <- errors.Wrap(err, "failure in HTTP listen and serve")
        }
    }()

    return &Service{
        Config:        cfg,
        Router:        r,
        API:           a,
        HealthCheck:   hc,
        ServiceList:   serviceList,
        Server:        s,
        MongoDBClient: mongoClient, // Assign the passed-in interface to the service
    }, nil
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

		// TODO: Close other dependencies, in the expected order
		if svc.MongoDBClient != nil {
			if err := svc.MongoDBClient.Close(ctx); err != nil {
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

func registerCheckers(_ context.Context, _ HealthChecker) (err error) {
	// TODO: add other health checks here, as per dp-upload-service

	return nil
}
