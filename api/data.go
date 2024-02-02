package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	errs "github.com/ONSdigital/dp-legacy-cache-api/apierrors"
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

	// Check request body not empty
	if req.ContentLength <= 0 {
		log.Info(ctx, "createOrUpdateCacheTime endpoint: empty request body")
		http.Error(w, "bad request: empty request body", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields() // disallow unknown fields in the request body

	err := decoder.Decode(&docToInsertOrUpdate)
	if err != nil {
		// Handle error for unknown fields, incorrect field type and decode
		log.Info(ctx, "createOrUpdateCacheTime endpoint: error decoding request body")
		http.Error(w, fmt.Errorf("bad request: %v", err).Error(), http.StatusBadRequest)
		return
	}

	// Validate request body
	err = isValidCacheTime(docToInsertOrUpdate)
	if err != nil {
		log.Info(ctx, "createOrUpdateCacheTime endpoint: cache time failed validation checks")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Upsert document into mongoDB.
	err = api.dataStore.UpsertCacheTime(ctx, docToInsertOrUpdate)
	if err != nil {
		log.Error(ctx, "createOrUpdateCacheTime endpoint: error upserting document", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCacheTime retrieves a cache time for a given ID and writes it to the HTTP response.
func (api *API) GetCacheTime(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Info(ctx, "calling get cache time handler")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	id := vars["id"]

	err := isValidID(id)
	if err != nil {
		log.Info(ctx, "getCacheTime endpoint: id failed validation checks")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	if err := json.NewEncoder(w).Encode(cacheTime); err != nil {
		log.Error(ctx, "error encoding results to JSON", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func isValidCacheTime(cacheTime *models.CacheTime) error {
	e := findIDErrors(cacheTime.ID)

	if cacheTime.ETag == "" {
		e = append(e, errors.New("etag field missing"))
	}
	if cacheTime.Path == "" {
		e = append(e, errors.New("path field missing"))
	}
	if len(e) > 0 {
		return fmt.Errorf("validation errors: %v", e)
	}
	return nil
}

func isValidID(id string) error {
	e := findIDErrors(id)
	if len(e) > 0 {
		return fmt.Errorf("validation errors: %v", e)
	}
	return nil
}

func findIDErrors(id string) []error {
	var e []error

	if len(id) != 32 {
		e = append(e, errors.New("id should be 32 characters in length"))
	}
	if !isLower(id) {
		e = append(e, errors.New("id is not lowercase"))
	}
	if !isHexadecimal(id) {
		e = append(e, errors.New("id is not a valid hexadecimal"))
	}
	return e
}

func isHexadecimal(s string) bool {
	hexRegex := regexp.MustCompile("^[0-9a-fA-F]+$")
	return hexRegex.MatchString(s)
}

func isLower(s string) bool {
	return strings.ToLower(s) == s
}
