package mongo

import (
	"context"

	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
// 	mongodb "github.com/ONSdigital/dp-mongodb/v3/mongodb"
)

type Mongo struct {
// 	config.MongoConfig
	mongodriver.MongoDriverConfig

	Connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
}



// type MongoDBClient interface {
//     NewMongoStore(ctx context.Context) error
//     Close(ctx context.Context) error
//     IsConnected(ctx context.Context) bool
// 	GetDataSets(ctx context.Context) ([]models.DataMessage, error)
// 	AddDataSet(ctx context.Context, dataset models.DataMessage) (*mongodriver.CollectionInsertResult, error)
// }

// func NewMongoStore(_ context.Context, cfg config.MongoDB) (m *Mongo, err error) {
// 	m = &Mongo{MongoDriverConfig: cfg}
//
// 	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)
// 	if err != nil {

func NewMongoStore(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	// instantiate mongo
	m = &Mongo{MongoDriverConfig: cfg}

	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)

	if err != nil {
		return nil, err
	}
	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
		mongohealth.Database(m.Database): {
			mongohealth.Collection(m.ActualCollectionName(config.DatasetsCollection)),
		},
	}

	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil


// 	m = &Mongo{MongoDriverConfig: cfg}
//
//     	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)
//     	if err != nil {
//     		return nil, err
//     	}
//
//     	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
//     		mongohealth.Database(m.Database): {
//     			mongohealth.Collection(m.ActualCollectionName(config.RolesCollection)),
//     			mongohealth.Collection(m.ActualCollectionName(config.PoliciesCollection)),
//     		},
//     	}
//     	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)
//
//     	return m, nil

}

func (m *Mongo) GetDataSets(ctx context.Context) (values []models.DataMessage, err error) {
	filter := bson.M{}

	var results []models.DataMessage

	// Finding the documents in the collection
	_, err = m.Connection.Collection("datasets").Find(ctx, filter, &results)
	if err != nil {
		log.Error(ctx, "Error finding collection: %v", err)
		return nil, err
	}
	return results, nil
}

func (m *Mongo) AddDataSet(ctx context.Context, dataset models.DataMessage) (*mongodriver.CollectionInsertResult, error) {
    result, err := m.Connection.Collection("datasets").InsertOne(ctx, dataset)
    if err != nil {
        log.Error(ctx, "Error inserting document into collection:", err)
        return nil, err
    }
    return result, nil
}

// Close represents mongo session closing within the context deadline
func (m *Mongo) Close(ctx context.Context) error {
	return m.Connection.Close(ctx)
}

// test function to check connection
func (m *Mongo) IsConnected(ctx context.Context) bool {
	if m.Connection == nil {
		return false
	}

	// Pinging the MongoDB server
	err := m.Connection.Ping(ctx, 5)
	return err == nil
}
