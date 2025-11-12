package sdk

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	apiError "github.com/ONSdigital/dp-legacy-cache-api/sdk/errors"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter

type Clienter interface {
	Checker(ctx context.Context, check *healthcheck.CheckState) error
	Health() *health.Client
	URL() string
	GetCacheTimes(ctx context.Context, auth Auth, options Options) (*models.CacheTimesList, apiError.Error)
}
