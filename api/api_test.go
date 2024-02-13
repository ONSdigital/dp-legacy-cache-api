package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/api/mock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		mockMongoDB := &mock.DataStoreMock{}

		Convey("When created in publishing subnet", func() {
			cacheAPI := setupPublishingAPI(mockMongoDB)

			Convey("Then all the routes should be available", func() {
				So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "GET"), ShouldBeTrue)
				So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "PUT"), ShouldBeTrue)
			})
		})

		Convey("When created in web subnet", func() {
			cacheAPI := setupWebAPI(mockMongoDB)

			Convey("Then the PUT endpoint should not have been added", func() {
				So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "GET"), ShouldBeTrue)
				So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "PUT"), ShouldBeFalse)
			})
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func setupAPI(isPublishing bool, dataStore api.DataStore) *api.API {
	mockIdentityHandler := func(h http.Handler) http.Handler {
		return h
	}

	return api.Setup(context.Background(), isPublishing, mux.NewRouter(), dataStore, mockIdentityHandler)
}

func setupPublishingAPI(dataStore api.DataStore) *api.API {
	return setupAPI(true, dataStore)
}

func setupWebAPI(dataStore api.DataStore) *api.API {
	return setupAPI(false, dataStore)
}
