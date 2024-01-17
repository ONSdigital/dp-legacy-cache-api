package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

var (
	read   = auth.Permissions{Read: true}
	update = auth.Permissions{Update: true}
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

	// GetPermissionsRequestBuilder for authorising access to datasets.
	// datasetPermissionsRequestBuilder := auth.NewDatasetPermissionsRequestBuilder("http://localhost:8082", "dataset_id", mux.Vars)

	// datasetsPermissions := auth.NewHandler(
	// 	datasetPermissionsRequestBuilder,
	// 	auth.DefaultPermissionsClient(),
	// 	auth.DefaultPermissionsVerifier(),
	// )

	// GetPermissionsRequestBuilder for authorising general CMD access (cases where we don't have a collection ID & dataset ID).
	permissionsRequestBuilder := auth.NewPermissionsRequestBuilder("http://localhost:8082")

	permissions := auth.NewHandler(
		permissionsRequestBuilder,
		auth.DefaultPermissionsClient(),
		auth.DefaultPermissionsVerifier(),
	)

	//r.HandleFunc("/datasets", permissions.Require(read, getDatasetsHandlerFunc)).Methods("GET")

	r.HandleFunc("/mongocheck", permissions.Require(update, api.AddDataSets(ctx))).Methods("POST")
	r.HandleFunc("/mongocheck", permissions.Require(read, api.GetDataSets(ctx))).Methods("GET")
	return api
}

// an example http.HandlerFunc for getting a dataset
func getDatasetsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Context(), "get datasets stub invoked")
	w.Write([]byte("datasets info here"))
}

// an example http.HandlerFunc for getting a dataset
func getDatasetHandlerFunc(w http.ResponseWriter, r *http.Request) {
	datasetID := mux.Vars(r)["dataset_id"]
	log.Info(r.Context(), "get dataset stub invoked", log.Data{"dataset_id": datasetID})
	w.Write([]byte("dataset info here"))
}
