package api

import (
	"context"

// 	"github.com/ONSdigital/dp-healthcheck/healthcheck"
 	"github.com/ONSdigital/dp-legacy-cache-api/models"
//  	"github.com/ONSdigital/dp-legacy-cache-api/config"
)

// moq -out mock/dataStore.go -pkg mock . DataStore
// moq -out ../service/mock/store.go -pkg mock . DataStore
// moq -out mock/bundler.go -pkg mock . DataBundler

//DataStore defines the behaviour of a PermissionsStore

type DataStore interface {
    Close(ctx context.Context) error
    IsConnected(ctx context.Context) bool
    GetDataSets(ctx context.Context) ([]models.DataMessage, error)
//     AddDataSet(ctx context.Context, dataset models.DataMessage) (*mongodriver.CollectionInsertResult, error)

}