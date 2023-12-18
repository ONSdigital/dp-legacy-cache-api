package mongo

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-legacy-cache-api/config"
	"go.mongodb.org/mongo-driver/bson"

	mongolock "github.com/ONSdigital/dp-mongodb/v3/dplock"
	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
)

type Mongo struct {
	config.MongoConfig

	Connection *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
	lockClient *mongolock.Lock
}

func (m *Mongo) Init(ctx context.Context) (err error) {

	//instantiate mongo 
	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)

	if err != nil {
		return err
	}
	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
		mongohealth.Database(m.Database): {
			mongohealth.Collection(m.ActualCollectionName(config.DatasetsCollection)),
		},
	}

	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return nil
}

// Close represents mongo session closing within the context deadline
func (m *Mongo) Close(ctx context.Context) error {
	return m.Connection.Close(ctx)
}

//test function to check connection
func (m *Mongo) IsConnected(ctx context.Context) bool {
    if m.Connection == nil {
        return false
    }

    // Pinging the MongoDB server
    err := m.Connection.Ping(ctx, 5)
    return err == nil
}

func (m *Mongo) CreateCollectionTest(ctx context.Context) error {
    // Connect to the collection
    collection := m.Connection.Collection("test-collection")

    // Example document to insert
    docToInsert := bson.M{"message": "hello"}

    // Insert a document
    _, err := collection.InsertOne(ctx, docToInsert)
    if err != nil {
        return fmt.Errorf("failed to insert document: %w", err)
    }

    // Find the inserted document
    var foundDoc bson.M
    err = collection.FindOne(ctx, bson.M{"message": "hello"}, &foundDoc)
    if err != nil {
        return fmt.Errorf("failed to find document: %w", err)
    }

    fmt.Println("Found document:", foundDoc)

	for key, value := range foundDoc {
		fmt.Println(key, value)
	}

    return nil
}

