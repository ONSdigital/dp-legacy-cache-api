package mongo

import (
	"context"

	"github.com/ONSdigital/dp-legacy-cache-api/config"
	mongolock "github.com/ONSdigital/dp-mongodb/v3/dplock"
	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
)

type Mongo struct {
	config.MongoConfig

	Connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
	lockClient   *mongolock.Lock
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

// test function to check connection
func (m *Mongo) IsConnected(ctx context.Context) bool {
	if m.Connection == nil {
		return false
	}

	// Pinging the MongoDB server
	err := m.Connection.Ping(ctx, 5)
	return err == nil
}
