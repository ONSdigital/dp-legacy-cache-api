package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/api/mock"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var testCacheID = "testCacheID"
var baseURL = "http://localhost:29100/v1/cache-times/"
var staticTime = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestGetCacheTimeEndpointReturns401WhenTokenIsMisssing(t *testing.T) {
	Convey("Given a GetCacheTime handler", t, func() {
		ctx := context.Background()
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

		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

		Convey("When an existing cache time is requested with its ID", func() {
			request := httptest.NewRequest(http.MethodGet, baseURL+testCacheID, http.NoBody)
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)

			Convey("The status code should be 401", func() {

				So(responseRecorder.Code, ShouldEqual, http.StatusUnauthorized)
			})
		})
	})
}

func TestGetCacheTimeEndpoint(t *testing.T) {
	Convey("Given a GetCacheTime handler", t, func() {
		ctx := context.Background()
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

		dataStoreAPI := setupAPINoAuthWithStore(ctx, dataStoreMock)

		Convey("When an existing cache time is requested with its ID", func() {
			request := httptest.NewRequest(http.MethodGet, baseURL+testCacheID, http.NoBody)
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
				// payload, _ := io.ReadAll(responseRecorder.Body)
				// err := json.Unmarshal(payload, &cacheTime)
				// So(err, ShouldBeNil)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
				So(cacheTime, ShouldEqual, expectedCacheTime)
			})
		})
	})
}

func TestUpdateExistingCacheTime(t *testing.T) {
	Convey("Given an existing cache time", t, func() {
		ctx := context.Background()
		db := make(map[string]models.CacheTime)
		dataStoreMock := &mock.DataStoreMock{
			UpsertCacheTimeFunc: func(ctx context.Context, cacheTime *models.CacheTime) error {
				db[cacheTime.ID] = *cacheTime
				return nil
			},
		}
		dataStoreAPI := setupAPINoAuthWithStore(ctx, dataStoreMock)

		existingCacheTime := models.CacheTime{
			ID:           testCacheID,
			Path:         "existingpath",
			ETag:         "existingetag",
			CollectionID: 123,
			ReleaseTime:  staticTime,
		}
		db[testCacheID] = existingCacheTime

		Convey("When updating the cache time", func() {
			updatedCacheTime := models.CacheTime{
				ID:           testCacheID,
				Path:         "updatedpath",
				ETag:         "updatedetag",
				CollectionID: 123,
				ReleaseTime:  staticTime,
			}
			payload, err := json.Marshal(updatedCacheTime)
			So(err, ShouldBeNil)
			reader := bytes.NewReader(payload)
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, reader)
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)

			Convey("Then the cache time should be updated with status code 204", func() {
				So(responseRecorder.Code, ShouldEqual, http.StatusNoContent)
				updatedRecord, exists := db[testCacheID]
				So(exists, ShouldBeTrue)
				So(updatedRecord, ShouldEqual, updatedCacheTime)
				So(responseRecorder.Body.Len(), ShouldEqual, 0)
			})
		})
	})
}

func TestCreateNewCacheTime(t *testing.T) {
	Convey("Given no existing cache time", t, func() {
		ctx := context.Background()
		db := make(map[string]models.CacheTime)
		dataStoreMock := &mock.DataStoreMock{
			UpsertCacheTimeFunc: func(ctx context.Context, cacheTime *models.CacheTime) error {
				db[cacheTime.ID] = *cacheTime
				return nil
			},
		}
		dataStoreAPI := setupAPINoAuthWithStore(ctx, dataStoreMock)

		Convey("When creating a new cache time", func() {
			newCacheTime := models.CacheTime{
				ID:           testCacheID,
				Path:         "newpath",
				ETag:         "newetag",
				CollectionID: 123,
				ReleaseTime:  staticTime,
			}
			payload, err := json.Marshal(newCacheTime)
			So(err, ShouldBeNil)
			reader := bytes.NewReader(payload)
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, reader)
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)

			Convey("Then a new cache time should be created with status code 204 with an empty response body", func() {
				So(responseRecorder.Code, ShouldEqual, http.StatusNoContent)
				createdRecord, exists := db[testCacheID]
				So(exists, ShouldBeTrue)
				So(createdRecord, ShouldEqual, newCacheTime)
				So(responseRecorder.Body.Len(), ShouldEqual, 0)
			})
		})
	})
}
