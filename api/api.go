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
func Setup(ctx context.Context, isPublishing bool, r *mux.Router, dataStore DataStore) *API {
	api := &API{
		Router:    r,
		dataStore: dataStore,
	}

	r.HandleFunc("/v1/cache-times/{id}", func(w http.ResponseWriter, req *http.Request) {
		api.GetCacheTime(ctx, w, req)
	}).Methods(http.MethodGet)

	if isPublishing {
		r.HandleFunc("/v1/cache-times/{id}", func(w http.ResponseWriter, req *http.Request) {
			api.CreateOrUpdateCacheTime(ctx, w, req)
		}).Methods(http.MethodPut)
	}
	return api
}
