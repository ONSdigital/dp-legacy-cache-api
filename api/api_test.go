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
		router := mux.NewRouter()
		ctx := context.Background()

		mockMongoDB := &mock.DataStoreMock{}
		Convey("When created in publishing subnet", func() {
			var isPublishing = true
			cacheAPI := api.Setup(ctx, isPublishing, router, mockMongoDB)

			Convey("Then all the routes should be available", func() {
				So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "GET"), ShouldBeTrue)
				So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "PUT"), ShouldBeTrue)
			})
		})
		Convey("When created in web subnet", func() {
			var isPublishing = false
			cacheAPI := api.Setup(ctx, isPublishing, router, mockMongoDB)

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

func setupAPIWithStore(ctx context.Context, isPublishing bool, dataStore api.DataStore) *api.API {
	return api.Setup(ctx, isPublishing, mux.NewRouter(), dataStore)
}
