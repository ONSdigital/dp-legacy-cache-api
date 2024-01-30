package api

import (
	"context"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router    *mux.Router
	dataStore DataStore
}

// Setup function sets up the api with authentication and returns an api
func Setup(ctx context.Context, router *mux.Router, dataStore DataStore) *API {
	api := &API{
		Router:    router,
		dataStore: dataStore,
	}

	api.get(
		"/v1/cache-times/{id}",
		api.isAuthenticated(
			func(w http.ResponseWriter, req *http.Request) { api.GetCacheTime(ctx, w, req) }),
	)

	api.put(
		"/v1/cache-times/{id}",
		api.isAuthenticated(func(w http.ResponseWriter, req *http.Request) { api.CreateOrUpdateCacheTime(ctx, w, req) }),
	)

	api.post(
		"/mongocheck",
		api.isAuthenticated(
			func(w http.ResponseWriter, req *http.Request) { api.AddDataSets(ctx) }),
	)

	api.get(
		"/mongocheck",
		api.isAuthenticated(
			func(w http.ResponseWriter, req *http.Request) { api.GetDataSets(ctx) }),
	)

	return api
}

// Setup function sets up the api and returns an api
func SetupNoAuth(ctx context.Context, router *mux.Router, dataStore DataStore) *API {
	api := &API{
		Router:    router,
		dataStore: dataStore,
	}

	api.get(
		"/v1/cache-times/{id}",
		func(w http.ResponseWriter, req *http.Request) { api.GetCacheTime(ctx, w, req) },
	)

	api.put(
		"/v1/cache-times/{id}",
		func(w http.ResponseWriter, req *http.Request) { api.CreateOrUpdateCacheTime(ctx, w, req) },
	)

	api.post(
		"/mongocheck",
		func(w http.ResponseWriter, req *http.Request) { api.AddDataSets(ctx) },
	)

	api.get(
		"/mongocheck",
		func(w http.ResponseWriter, req *http.Request) { api.GetDataSets(ctx) },
	)

	return api
}

// isAuthenticated wraps a http handler func in another http handler func that checks the caller is authenticated to
// perform the requested action. handler is the http.HandlerFunc to wrap in an
// authentication check. The wrapped handler is only called if the caller is authenticated
func (api *API) isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return dphandlers.CheckIdentity(handler)
}

// get registers a GET http.HandlerFunc.
func (api *API) get(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods(http.MethodGet)
}

// put registers a PUT http.HandlerFunc.
func (api *API) put(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods(http.MethodPut)
}

// post registers a POST http.HandlerFunc.
func (api *API) post(path string, handler http.HandlerFunc) {
	api.Router.HandleFunc(path, handler).Methods(http.MethodPost)
}
