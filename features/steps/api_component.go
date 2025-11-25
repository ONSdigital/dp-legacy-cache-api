package steps

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"github.com/ONSdigital/dp-legacy-cache-api/mongo"
	"github.com/ONSdigital/dp-legacy-cache-api/service"
	"github.com/ONSdigital/dp-legacy-cache-api/service/mock"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type Component struct {
	componenttest.ErrorFeature
	svcList        *service.ExternalServiceList
	svc            *service.Service
	errorChan      chan error
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
	apiFeature     *componenttest.APIFeature
	authFeature    *componenttest.AuthorizationFeature
	MongoClient    *mongo.Mongo
}

func NewComponent(mongoURI, mongoDatabaseName string) (*Component, error) {
	c := &Component{
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	parsedURI, err := url.Parse(mongoURI)
	if err != nil {
		return nil, err
	}
	hostPort := parsedURI.Host

	c.Config.IsPublishing = true
	c.Config.ClusterEndpoint = hostPort
	c.Config.Database = mongoDatabaseName

	c.MongoClient, err = mongo.NewMongoStore(context.Background(), c.Config.MongoConfig)
	if err != nil {
		return nil, err
	}

	initMock := &mock.InitialiserMock{
		DoGetHealthCheckFunc: c.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  c.DoGetHTTPServer,
		DoGetMongoDBFunc:     c.DoGetMongoDB,
	}

	c.svcList = service.NewServiceList(initMock)

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)
	c.authFeature = componenttest.NewAuthorizationFeature()

	c.Config.ZebedeeURL = c.authFeature.FakeAuthService.ResolveURL("")

	return c, nil
}

func (c *Component) Reset() *Component {
	c.apiFeature.Reset()
	c.authFeature.Reset()

	c.authFeature.FakeAuthService.NewHandler().Get("/identity").Reply(200).BodyString(`{ "identifier": "svc-authenticated"}`)

	return c
}

func (c *Component) Close() error {
	ctx := context.Background()

	if c.svc != nil && c.ServiceRunning {
		c.authFeature.Close()
		if err := c.MongoClient.Connection.DropDatabase(ctx); err != nil {
			return err
		}
		if err := c.svc.Close(ctx); err != nil {
			return err
		}
		c.ServiceRunning = false
	}
	return nil
}

func (c *Component) InitialiseService() (http.Handler, error) {
	var err error
	c.svc, err = service.Run(context.Background(), c.Config, c.svcList, "1", "", "", c.errorChan)
	if err != nil {
		return nil, err
	}

	c.ServiceRunning = true
	return c.HTTPServer.Handler, nil
}

func (c *Component) DoGetHealthcheckOk(_ *config.Config, _, _, _ string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer = &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              bindAddr,
		Handler:           router,
	}
	return c.HTTPServer
}

func (c *Component) DoGetMongoDB(_ context.Context, _ *config.Config) (service.DataStore, error) {
	return c.MongoClient, nil
}
