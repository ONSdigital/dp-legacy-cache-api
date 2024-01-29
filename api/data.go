package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

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
		log.Error(ctx, "error retrieving cache time from datastore", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(cacheTime); err != nil {
		log.Error(ctx, "error encoding results to JSON", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
