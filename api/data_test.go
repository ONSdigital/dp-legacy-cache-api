package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/api/mock"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var testCacheID = "testCacheID"

func TestGetCacheTimeEndpoint(t *testing.T) {
	Convey("Given a GetCacheTime handler", t, func() {
		ctx := context.Background()
		staticTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
		dataStoreMock := &mock.DataStoreMock{
			GetCacheTimeFunc: func(ctx context.Context, id string) (*models.CacheTime, error) {
				switch id {
				case testCacheID:
					return &models.CacheTime{
						ID:           testCacheID,
						Path:         "testpath",
						ETag:         "testetag",
						CollectionID: 123,
						ReleaseTime:  staticTime,
					}, nil
				default:
					return nil, errors.New("Something went wrong")
				}
			},
		}
		datasetPermissions := getAuthorisationHandlerMock()
		permissions := getAuthorisationHandlerMock()
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock, datasetPermissions, permissions)

		Convey("When an existing cache time is requested with its ID", func() {
			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:29100/v1/cache-times/%s", testCacheID), http.NoBody)
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)

			Convey("The matched cache time is returned with status code 200", func() {
				expectedCacheTime := models.CacheTime{
					ID:           testCacheID,
					Path:         "testpath",
					ETag:         "testetag",
					CollectionID: 123,
					ReleaseTime:  staticTime,
				}
				cacheTime := models.CacheTime{}
				payload, _ := io.ReadAll(responseRecorder.Body)
				err := json.Unmarshal(payload, &cacheTime)
				So(err, ShouldBeNil)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
				So(cacheTime, ShouldEqual, expectedCacheTime)
			})
		})
	})
}
