package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/log.go/v2/log"
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
	log.Info(ctx, "api contains example endpoint, remove hello.go as soon as possible")

	return func(w http.ResponseWriter, req *http.Request) {
		// Deconstruct json into our models.DataMessage struct
		var input models.DataMessage
		err := json.NewDecoder(req.Body).Decode(&input)

		if err != nil {
			log.Error(ctx, "error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		log.Info(ctx, "received data", log.Data{"input:": input})

		err = api.dataStore.AddDataSet(ctx, input)
		if err != nil {
			log.Error(ctx, "failed to insert document", err)
			return
		}

		log.Info(ctx, "successfully inserted document")

		// Respond with the inserted document and StatusCreated
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(input)
		if err != nil {
			log.Error(ctx, "error encoding results to JSON", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
