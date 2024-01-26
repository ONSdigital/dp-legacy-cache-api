package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-playground/validator/v10"
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

	// Check request body not empty
	if req.ContentLength == 0 {
		log.Error(ctx, "empty request body", http.ErrContentLength) // TO-DO: confirm ErrContentLength is correct implementation
		http.Error(w, "Bad Request: Empty Request Body", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&docToInsertOrUpdate)

	if err != nil {
		log.Error(ctx, "error decoding request body", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fmt.Println(docToInsertOrUpdate)
	// Validate request body
	validErrs := isValidCacheTime(docToInsertOrUpdate)
	if validErrs != nil {
		errMsg := "Validation Errors: "
		for _, vErr := range validErrs {
			// Construct a single error message string from all validation errors
			errMsg += vErr.Error() + "; "
		}

		// Create an error object from the concatenated error message
		combinedError := errors.New(errMsg)

		// Log the error along with the context and the error object
		log.Error(ctx, errMsg, combinedError)

		// Send an appropriate HTTP response back to the client
		// The status code 400 Bad Request is typically used for validation errors
		http.Error(w, errMsg, http.StatusBadRequest)
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

func isValidCacheTime(cacheTime *models.CacheTime) []error {
	errs := []error{}
	var validate *validator.Validate = validator.New()
	err := validate.Struct(cacheTime)
	if err != nil {
		// Use the FieldError interface to get detailed information
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New(err.StructNamespace()+" "+err.Tag()))
		}
	}
	return errs
}

// func CreateDataset(reader io.Reader) (*Dataset, error) {
// 	b, err := io.ReadAll(reader)
// 	if err != nil {
// 		return nil, errs.ErrUnableToReadMessage
// 	}

// 	var dataset Dataset

// 	err = json.Unmarshal(b, &dataset)
// 	if err != nil {
// 		return nil, errs.ErrUnableToParseJSON
// 	}

// 	return &dataset, nil
// }
