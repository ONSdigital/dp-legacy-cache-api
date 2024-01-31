package steps

import (
	"context"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
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
	MongoClient    *mongo.Mongo
}

func NewComponent() (*Component, error) {
	c := &Component{
		HTTPServer:     &http.Server{ReadHeaderTimeout: 3 * time.Second},
		errorChan:      make(chan error),
		ServiceRunning: false,
		Config:         nil,
	}

	var err error
	c.Config, err = config.Get()

	if err != nil {
		return nil, err
	}

	initMock := &mock.InitialiserMock{
		DoGetHealthCheckFunc: c.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  c.DoGetHTTPServer,
		DoGetMongoDBFunc:     c.DoGetMongoDB,
	}

	c.svcList = service.NewServiceList(initMock)
	c.svc = service.New(c.Config, c.svcList)
	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c, nil
}

func (c *Component) Reset() *Component {
	c.apiFeature.Reset()
	return c
}

func (c *Component) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *Component) InitialiseService() (http.Handler, error) {
	err := c.svc.Run(context.Background(), c.Config, c.svcList, "1", "", "", c.errorChan)
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
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *Component) DoGetMongoDB(_ context.Context, _ *config.Config) (service.DataStore, error) {
	return &mock.DataStoreMock{
		CloseFunc:       func(ctx context.Context) error { return nil },
		GetDataSetsFunc: func(ctx context.Context) ([]models.DataMessage, error) { return []models.DataMessage{}, nil },
	}, nil
}
