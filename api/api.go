package api

import (
	"context"
	"net/http"

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

	// TODO: implement write endpoint here (DIS-328)
	r.HandleFunc("/v1/cache-times/{id}", func(w http.ResponseWriter, req *http.Request) {
		api.GetCacheTime(ctx, w, req)
	}).Methods(http.MethodGet)
	return api
}
