package api

import (
	"context"

	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router    *mux.Router
	dataStore DataStore
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, dataStore DataStore) *API {
	api := &API{
		Router:    r,
		dataStore: dataStore,
	}

	r.HandleFunc("/mongocheck", api.AddDataSets(ctx)).Methods("POST")
	r.HandleFunc("/mongocheck", api.GetDataSets(ctx)).Methods("GET")
	return api
}
