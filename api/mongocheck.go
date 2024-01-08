package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

func (api *API) GetDataSets(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Info(ctx, "Calling get datasets handler")

		results, err := api.MongoClient.GetDataSets(ctx)
		if err != nil {
			log.Error(ctx, "Error retrieving datasets from db:", err)
			return
		}

		// Setting the header and encoding the results to JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Error(ctx, "Error encoding results to JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// HelloHandler returns function containing a simple hello world example of an api handler
func (api *API) AddDataSets(ctx context.Context) http.HandlerFunc {
	log.Info(ctx, "api contains example endpoint, remove hello.go as soon as possible")

	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		// Deconstruct json into our models.DataMessage struct
		var input models.DataMessage
		err := json.NewDecoder(req.Body).Decode(&input)

		if err != nil {
			log.Error(ctx, "error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// log the received data
		fmt.Println("received data:", input)


		result, err := api.MongoClient.AddDataSet(ctx, input)
		if err != nil {
			log.Error(ctx, "failed to insert document: %w", err)
			return
		}

		insertedID := result.InsertedId
		log.Info(ctx, "Document inserted successfully", log.Data{"insertedID": insertedID})

		// Respond with the inserted document and StatusCreated
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(input)
		if err != nil {
			log.Error(ctx, "Error encoding results to JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
