package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/gorilla/mux"
)

// AuthHandler provides authorisation checks on requests
type AuthHandler interface {
	Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}

// API provides a struct to wrap the api around
type API struct {
	Router             *mux.Router
	dataStore          DataStore
	datasetPermissions AuthHandler
	permissions        AuthHandler
}

var (
	createPermission = auth.Permissions{Create: true}
	readPermission   = auth.Permissions{Read: true}
	updatePermission = auth.Permissions{Update: true}
	deletePermission = auth.Permissions{Delete: true}
)

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, router *mux.Router, dataStore DataStore, datasetPermissions AuthHandler, permissions AuthHandler) *API {
	api := &API{
		Router:             router,
		dataStore:          dataStore,
		datasetPermissions: datasetPermissions,
		permissions:        permissions,
	}

	// router.HandleFunc("/mongocheck", api.AddDataSets(ctx)).Methods("POST")
	// router.HandleFunc("/mongocheck", api.GetDataSets(ctx)).Methods("GET")

	// // TODO: implement write endpoint here (DIS-328)
	// router.HandleFunc("/v1/cache-times/{id}", func(w http.ResponseWriter, req *http.Request) {
	// 	api.GetCacheTime(ctx, w, req)
	// }).Methods(http.MethodGet)

	router.Path("/mongocheck").Methods("POST").HandlerFunc(dphandlers.CheckIdentity(api.AddDataSets(ctx)))
	router.Path("/mongocheck").Methods("GET").HandlerFunc(dphandlers.CheckIdentity(api.GetDataSets(ctx)))
	router.Path("/v1/cache-times/{id}").Methods("GET").HandlerFunc(api.isAuthenticated(api.isAuthorised(readPermission, func(w http.ResponseWriter, req *http.Request) { api.GetCacheTime(ctx, w, req) })))

	return api
}

// api.put(
// 	"/instances/{instance_id}",
// 	api.isAuthenticated(
// 		api.isAuthorised(updatePermission,
// 			api.isInstancePublished(instanceAPI.Update))),
// )

// isAuthenticated wraps a http handler func in another http handler func that checks the caller is authenticated to
// perform the requested action. handler is the http.HandlerFunc to wrap in an
// authentication check. The wrapped handler is only called if the caller is authenticated
func (api *API) isAuthenticated(handler http.HandlerFunc) http.HandlerFunc {
	return dphandlers.CheckIdentity(handler)
}

// isAuthorised wraps a http.HandlerFunc another http.HandlerFunc that checks the caller is authorised to perform the
// requested action. required is the permissions required to perform the action, handler is the http.HandlerFunc to
// apply the check to. The wrapped handler is only called if the caller has the required permissions.
func (api *API) isAuthorised(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
	return api.permissions.Require(required, handler)
}
