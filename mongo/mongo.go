package mongo

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
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
			mongoHealth.Collection(m.ActualCollectionName(config.DatasetsCollection)),
		},
	}

	m.healthClient = mongoHealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil
}

// GetDataSets reads all records in the connected database collection
func (m *Mongo) GetDataSets(ctx context.Context) (values []models.DataMessage, err error) {
	filter := bson.M{}

	var results []models.DataMessage

	_, err = m.Connection.Collection(config.DatasetsCollection).Find(ctx, filter, &results)
	if err != nil {
		log.Error(ctx, "error finding collection", err)
		return nil, err
	}
	return results, nil
}

// AddDataSet stores one dataset in the connected database
func (m *Mongo) AddDataSet(ctx context.Context, dataset models.DataMessage) error {
	_, err := m.Connection.Collection(config.DatasetsCollection).InsertOne(ctx, dataset)
	if err != nil {
		log.Error(ctx, "error inserting document into collection:", err)
		return err
	}
	return nil
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