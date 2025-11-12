package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	errs "github.com/ONSdigital/dp-legacy-cache-api/apierrors"
	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	mongoHealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type Mongo struct {
	mongoDriver.MongoDriverConfig

	Connection   *mongoDriver.MongoConnection
	healthClient *mongoHealth.CheckMongoClient
}

// NewMongoStore creates a connection to mongo database
func NewMongoStore(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	m = &Mongo{MongoDriverConfig: cfg}
	m.Connection, err = mongoDriver.Open(&m.MongoDriverConfig)

	if err != nil {
		return nil, err
	}
	databaseCollectionBuilder := map[mongoHealth.Database][]mongoHealth.Collection{
		mongoHealth.Database(m.Database): {
			mongoHealth.Collection(m.ActualCollectionName(config.CacheTimesCollection)),
		},
	}

	m.healthClient = mongoHealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// Close closes the connection to the database
func (m *Mongo) Close(ctx context.Context) error {
	return m.Connection.Close(ctx)
}

// IsConnected return the connection status to the db
func (m *Mongo) IsConnected(ctx context.Context) bool {
	if m.Connection == nil {
		return false
	}

	// Pinging the MongoDB server
	err := m.Connection.Ping(ctx, 5)
	return err == nil
}

// GetCacheTime returns a cache time with its given id
func (m *Mongo) GetCacheTime(ctx context.Context, id string) (*models.CacheTime, error) {
	filter := bson.M{"_id": id}

	var result models.CacheTime
	err := m.Connection.Collection(m.ActualCollectionName(config.CacheTimesCollection)).FindOne(ctx, filter, &result)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocumentFound) {
			log.Info(ctx, "api.dataStore.GetCacheTime document not found")
			return nil, errs.ErrCacheTimeNotFound
		}
		log.Error(ctx, "error targeting api.dataStore.GetCacheTime", err)
		return nil, errs.ErrDataStore
	}
	return &result, nil
}

// GetCacheTimes
func (m *Mongo) GetCacheTimes(ctx context.Context, offset, limit int, releaseTime time.Time) ([]*models.CacheTime, int, error) {
	var filter bson.M

	if !releaseTime.IsZero() {
		filter = bson.M{
			"release_time": buildDateTimeFilter(releaseTime),
		}
	}

	results := []*models.CacheTime{}
	totalCount, err := m.Connection.Collection(m.ActualCollectionName(config.CacheTimesCollection)).
		Find(ctx,
			filter,
			&results,
			mongoDriver.Offset(offset),
			mongoDriver.Limit(limit),
			mongoDriver.Sort(bson.D{{Key: "_id", Value: 1}}),
		)
	if err != nil {
		log.Error(ctx, "error targeting api.dataStore.GetCacheTimes", err)
		return nil, 0, errs.ErrDataStore
	}

	log.Info(ctx, "results from query", log.Data{"results": results, "total_count": totalCount})

	return results, totalCount, nil
}

// UpsertCacheTime adds or overrides an existing cache time
func (m *Mongo) UpsertCacheTime(ctx context.Context, cacheTime *models.CacheTime) (err error) {
	update := bson.M{
		"$set": bson.M{"path": cacheTime.Path, "collection_id": cacheTime.CollectionID, "release_time": cacheTime.ReleaseTime},
	}
	selector := bson.M{"_id": cacheTime.ID}

	_, err = m.Connection.Collection(m.ActualCollectionName(config.CacheTimesCollection)).UpsertOne(ctx, selector, update)

	return err
}

// buildDateTimeFilter builds a bson filter for the datetime with a window of duration around the datetime
func buildDateTimeFilter(date time.Time) bson.M {
	const timeWindow = 2 * time.Second

	startTime := date.Add(timeWindow * -1)
	endTime := date.Add(timeWindow)

	return bson.M{
		"$gte": startTime,
		"$lte": endTime,
	}
}
