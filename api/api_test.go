package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/api/mock"
	"github.com/ONSdigital/dp-legacy-cache-api/mocks"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

type AuthHandler interface {
	Require(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		router := mux.NewRouter()
		ctx := context.Background()

		mockMongoDB := &mock.DataStoreMock{}
		datasetPermissions := getAuthorisationHandlerMock()
		permissions := getAuthorisationHandlerMock()
		cacheAPI := api.Setup(ctx, router, mockMongoDB, datasetPermissions, permissions)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(cacheAPI.Router, "/mongocheck", "PUT"), ShouldBeTrue)
			So(hasRoute(cacheAPI.Router, "/mongocheck", "GET"), ShouldBeTrue)
			So(hasRoute(cacheAPI.Router, "/v1/cache-times/{id}", "GET"), ShouldBeTrue)
		})
	})
}

func getAuthorisationHandlerMock() *mocks.AuthHandlerMock {
	return &mocks.AuthHandlerMock{
		Required: &mocks.PermissionCheckCalls{Calls: 0},
	}
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func setupAPIWithStore(ctx context.Context, dataStore api.DataStore, datasetPermissions AuthHandler, permissions AuthHandler) *api.API {

	return api.Setup(ctx, mux.NewRouter(), dataStore, datasetPermissions, permissions)
}
