package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

type DataMessage struct {
	Message string `json:"message,omitempty"`
}

// HelloHandler returns function containing a simple hello world example of an api handler
func (api *API) MongoCheckHandler(ctx context.Context) http.HandlerFunc {
	log.Info(ctx, "api contains example endpoint, remove hello.go as soon as possible")

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		//Deconstruct json into our DataMessage struct
		var input DataMessage
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			log.Error(ctx, "error decoding request body", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		//log the received data
		fmt.Println("received data:", input)

		//Insert Document to MongoDB
		collection := api.MongoClient.Connection.Collection("datasets")

		result, err := collection.InsertOne(ctx, input)
		if err != nil {
			log.Error(ctx, "failed to insert document: %w", err)
			return
		}

		insertedID := result.InsertedId
		log.Info(ctx, "Document inserted successfully", log.Data{"insertedID": insertedID})

		// Respond with the inserted document and StatusCreated
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(input) // Send back the input as the inserted document

	}
}
