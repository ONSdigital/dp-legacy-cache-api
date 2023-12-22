package api

import (
	"context"

	"github.com/ONSdigital/dp-legacy-cache-api/mongo"
	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router      *mux.Router
	MongoClient *mongo.Mongo
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, mongoDB *mongo.Mongo) *API {
	api := &API{
		Router:      r,
		MongoClient: mongoDB,
	}

	// TODO: remove hello world example handler route
	r.HandleFunc("/hello", HelloHandler(ctx)).Methods("GET")
	r.HandleFunc("/mongocheck", api.AddDataSets(ctx)).Methods("POST")
	r.HandleFunc("/mongocheck", api.GetDataSets(ctx)).Methods("GET")
	return api
}
