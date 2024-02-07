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

// Setup function sets up the api with authentication and returns an api
func Setup(ctx context.Context, router *mux.Router, dataStore DataStore, identityHandler func(http.Handler) http.Handler) *API {
	api := &API{
		Router:          router,
		dataStore:       dataStore,
		identityHandler: identityHandler,
	}

	api.get(
		"/v1/cache-times/{id}",
		func(w http.ResponseWriter, req *http.Request) { api.GetCacheTime(ctx, w, req) },
	)

	api.put(
		"/v1/cache-times/{id}",
		api.isAuthenticated(func(w http.ResponseWriter, req *http.Request) { api.CreateOrUpdateCacheTime(ctx, w, req) }),
	)
	return api
}

// isAuthenticated wraps a http handler func in another http handler func that checks the callers identity and if it is authenticated to
// perform the requested action. handler is the http.HandlerFunc to wrap in an
// authentication check. The wrapped handler is only called if the caller is authenticated
func (api *API) isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		checkIdentityHandler := api.identityHandler(dphandlers.CheckIdentity(handler)) // -> handler chain: call identityhandler first then checkIdentity then when  identity exists the final handler is called
		checkIdentityHandler.ServeHTTP(w, req)
	}
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
