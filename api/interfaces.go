package api

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
)

//go:generate moq -out mock/dataStore.go -pkg mock . DataStore
//go:generate moq -out ../service/mock/store.go -pkg mock . DataStore

// DataStore defines the behaviour of a DataStore
type DataStore interface {
	Checker(ctx context.Context, state *healthcheck.CheckState) error
	Close(ctx context.Context) error
	IsConnected(ctx context.Context) bool
	GetCacheTime(ctx context.Context, id string) (*models.CacheTime, error)
	GetCacheTimes(ctx context.Context, offset int, limit int, releaseTime time.Time) ([]*models.CacheTime, int, error)
	UpsertCacheTime(ctx context.Context, cacheTime *models.CacheTime) error
}
