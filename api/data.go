package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/ONSdigital/dp-legacy-cache-api/apierrors"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetDataSets reads all messages from the datastore
func (api *API) GetDataSets(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Info(ctx, "calling get datasets handler")

		results, err := api.dataStore.GetDataSets(ctx)
		if err != nil {
			log.Error(ctx, "error retrieving datasets from db", err)
			return
		}

		// Setting the header and encoding the results to JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Error(ctx, "error encoding results to JSON", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// AddDataSets stores a message in the datastore
func (api *API) AddDataSets(ctx context.Context) http.HandlerFunc {
	log.Info(ctx, "calling add datsets handler")

	return func(w http.ResponseWriter, req *http.Request) {
		var docToInsert models.DataMessage
		err := json.NewDecoder(req.Body).Decode(&docToInsert)

		if err != nil {
			log.Error(ctx, "error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		log.Info(ctx, "received data", log.Data{"doc to insert:": docToInsert})

		err = api.dataStore.AddDataSet(ctx, docToInsert)
		if err != nil {
			log.Error(ctx, "failed to insert document", err)
			return
		}

		log.Info(ctx, "successfully inserted document")

		// Respond with the inserted document and StatusCreated
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(docToInsert)
		if err != nil {
			log.Error(ctx, "error encoding results to JSON", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// CreateOrUpdateCacheTime handles the creation or update of a cache time
func (api *API) CreateOrUpdateCacheTime(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Info(ctx, "calling create or update cache time handler")

	vars := mux.Vars(req)
	id := vars["id"]

	var docToInsertOrUpdate = &models.CacheTime{
		ID: id,
	}

	err := json.NewDecoder(req.Body).Decode(&docToInsertOrUpdate)
	if err != nil {
		log.Error(ctx, "error decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = api.dataStore.UpsertCacheTime(ctx, docToInsertOrUpdate)
	if err != nil {
		log.Error(ctx, "error upserting document", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCacheTime retrieves a cache time for a given ID and writes it to the HTTP response.
func (api *API) GetCacheTime(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Info(ctx, "calling get cache time handler")
	vars := mux.Vars(req)
	id := vars["id"]

	cacheTime, err := api.dataStore.GetCacheTime(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrCacheTimeNotFound) {
			log.Info(ctx, "getCacheTime endpoint: api.dataStore.GetCacheTime document not found")
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			log.Error(ctx, "getCacheTime endpoint: api.dataStore.GetCacheTime internal server error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(cacheTime); err != nil {
		log.Error(ctx, "error encoding results to JSON", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
