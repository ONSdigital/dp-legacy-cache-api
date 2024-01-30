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
		cacheAPI := api.Setup(ctx, router, mockMongoDB)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(cacheAPI.Router, "/mongocheck", "POST"), ShouldBeTrue)
			So(hasRoute(cacheAPI.Router, "/mongocheck", "GET"), ShouldBeTrue)
			So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "PUT"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func setupAPIWithStore(ctx context.Context, dataStore api.DataStore) *api.API {
	return api.Setup(ctx, mux.NewRouter(), dataStore)
}

func setupAPINoAuthWithStore(ctx context.Context, dataStore api.DataStore) *api.API {
	return api.SetupNoAuth(ctx, mux.NewRouter(), dataStore)
}
