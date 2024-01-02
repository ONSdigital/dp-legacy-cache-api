package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type DataMessage struct {
	Message string `json:"message,omitempty"`
}

func (api *API) GetDataSets(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Info(ctx, "Calling get datasets handler")

		// Empty filter to fetch all documents
		filter := bson.M{}

		// Slice to hold the results
		var results []DataMessage

		// Getting the collection
		collection := api.MongoClient.Connection.Collection("datasets")

		// Finding documents in the collection
		_, err := collection.Find(ctx, filter, &results)
		if err != nil {
			log.Error(ctx, "Error finding collection: %v", err)
			return
		}

		// Setting the header and encoding the results to JSON
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Error(ctx, "Error encoding results to JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// HelloHandler returns function containing a simple hello world example of an api handler
func (api *API) AddDataSets(ctx context.Context) http.HandlerFunc {
	log.Info(ctx, "api contains example endpoint, remove hello.go as soon as possible")

	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		//Deconstruct json into our DataMessage struct
		var input DataMessage
		err := json.NewDecoder(req.Body).Decode(&input)
		
		if err != nil {
			log.Error(ctx, "error decoding request body", err)
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
		json.NewEncoder(w).Encode(input) 
	}
}
