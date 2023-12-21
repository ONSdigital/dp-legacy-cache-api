package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-api/mongo"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()

		// Create a mock MongoDB client
		mockMongoDB := createMockMongoDBClient()

		api := Setup(ctx, r, mockMongoDB)

		// TODO: remove hello world example handler route test case
		Convey("When created the following routes should have been added", func() {
			// Replace the check below with any newly added api endpoints
			So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/mongocheck", "POST"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

// Mock MongoDB client creation function
func createMockMongoDBClient() *mongo.Mongo {
	// Return a mock or dummy MongoDB client
	// The implementation of this depends on your specific MongoDB client structure and needs
	return &mongo.Mongo{} // Replace with actual mock implementation
}
