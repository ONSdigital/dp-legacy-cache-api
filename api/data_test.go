package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/api/mock"
	errs "github.com/ONSdigital/dp-legacy-cache-api/apierrors"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var testCacheID = "testCacheID"
var baseURL = "http://localhost:29100/v1/cache-times/"
var staticTime = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

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
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

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
				payload, _ := io.ReadAll(responseRecorder.Body)
				err := json.Unmarshal(payload, &cacheTime)
				So(err, ShouldBeNil)
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
				So(cacheTime, ShouldEqual, expectedCacheTime)
			})
		})
	})
}

func TestGetCacheTimeReturnsError404(t *testing.T) {
	Convey("Given a GetCacheTime handler", t, func() {
		ctx := context.Background()
		mockedDataStore := &mock.DataStoreMock{
			GetCacheTimeFunc: func(ctx context.Context, id string) (*models.CacheTime, error) {
				return nil, errs.ErrCacheTimeNotFound
			},
		}
		dataStoreAPI := setupAPIWithStore(ctx, mockedDataStore)
		Convey("When a non-existent cache time is requested with its ID", func() {
			var nonExistentCacheID = "abcdef0a1b2c3d4e5f67890123456789"
			r := httptest.NewRequest(http.MethodGet, baseURL+nonExistentCacheID, http.NoBody)
			w := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(w, r)
			Convey("Then unmatched cache time is not found with status code 404", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetCacheTimeReturnsError500(t *testing.T) {
	Convey("Given a GetCacheTime handler", t, func() {
		ctx := context.Background()
		mockedDataStore := &mock.DataStoreMock{
			GetCacheTimeFunc: func(ctx context.Context, id string) (*models.CacheTime, error) {
				return nil, errs.ErrDataStore
			},
		}
		dataStoreAPI := setupAPIWithStore(ctx, mockedDataStore)
		Convey("When there is an error with the datastore ", func() {
			r := httptest.NewRequest(http.MethodGet, baseURL+testCacheID, http.NoBody)
			w := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(w, r)
			Convey("Then return an internal server error with status code 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
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
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

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
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

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
