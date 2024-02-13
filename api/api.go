package api

import (
	"context"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router          *mux.Router
	dataStore       DataStore
	identityHandler func(http.Handler) http.Handler
}

// Setup function sets up the api and returns an API
func Setup(ctx context.Context, isPublishing bool, r *mux.Router, dataStore DataStore, identityHandler func(http.Handler) http.Handler) *API {
	api := &API{
		Router:          r,
		dataStore:       dataStore,
		identityHandler: identityHandler,
	}

	api.get(
		"/v1/cache-times/{id}",
		func(w http.ResponseWriter, req *http.Request) { api.GetCacheTime(ctx, w, req) },
	)

	if isPublishing {
		api.put(
			"/v1/cache-times/{id}",
			api.isAuthenticated(func(w http.ResponseWriter, req *http.Request) { api.CreateOrUpdateCacheTime(ctx, w, req) }),
		)
	}

	return api
}

func (api *API) isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		checkIdentityHandler := api.identityHandler(dphandlers.CheckIdentity(handler))
		checkIdentityHandler.ServeHTTP(w, req)
	}
}

func (api *API) get(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods(http.MethodGet)
}

func (api *API) put(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods(http.MethodPut)
}
