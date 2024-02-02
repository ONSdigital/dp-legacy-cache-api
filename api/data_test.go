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
	dprequest "github.com/ONSdigital/dp-net/request"
	. "github.com/smartystreets/goconvey/convey"
)

var testCacheID = "a1b2c3d4e5f67890123456789abcdef0"
var baseURL = "http://localhost:29100/v1/cache-times/"
var staticTime = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

func TestGetCacheTimeEndpointReturns200WhenAuthIsOnAndTokenIsMissing(t *testing.T) {
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
					return nil, errors.New("something went wrong")
				}
			},
		}

		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

		Convey("When an existing cache time is requested with its ID and request is NOT authenticated", func() {
			request := httptest.NewRequest(http.MethodGet, baseURL+testCacheID, http.NoBody)
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)

			Convey("The status code should be 200", func() {
				So(responseRecorder.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestUpdateCacheTimeReturns401WhenAuthIsOnAndAuthTokenIsMissing(t *testing.T) {
	Convey("Given a GetCacheTime handler", t, func() {
		ctx := context.Background()
		dataStoreMock := &mock.DataStoreMock{}
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

		Convey("When trying an unauthenticated request to the put cache-times endpoint", func() {
			payload, err := json.Marshal("")
			So(err, ShouldBeNil)
			reader := bytes.NewReader(payload)
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, reader)
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
					return nil, errors.New("something went wrong")
				}
			},
		}
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

		Convey("When an existing cache time is requested with its ID with no authenticationen set", func() {
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

		Convey("When updating the cache time with an authenticated request", func() {
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
			request := createRequestWithAuth(http.MethodPut, baseURL+testCacheID, reader)
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

		Convey("When creating a new cache time request with authentication", func() {
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
			request := createRequestWithAuth(http.MethodPut, baseURL+testCacheID, reader)
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
func TestCreateOrUpdateCacheTimeReturnsErr(t *testing.T) {
	const validBody = `{"path": "testpath", "etag": "testetag"}`

	Convey("Given an API", t, func() {
		ctx := context.Background()
		dataStoreMock := &mock.DataStoreMock{}
		dataStoreAPI := setupAPIWithStore(ctx, dataStoreMock)

		Convey("When no request body is provided and the CreateOrUpdateCacheTime endpoint is called", func() {
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, http.NoBody)
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with 'empty request body' in the response", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, "empty request body")
			})
		})

		Convey("When an etag and/or path is not provided and the CreateOrUpdateCacheTime endpoint is called", func() {
			staticTimeString := staticTime.Format(time.RFC3339)
			body := `{"collection_id": 123, "release_time":"` + staticTimeString + `"}`
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, bytes.NewBufferString(body))
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with the missing fields in the response", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, "[etag field missing path field missing]")
			})
		})

		Convey("When an extra field is provided and the CreateOrUpdateCacheTime endpoint is called", func() {
			body := `{"path": "testpath", "etag": "testetag", "extra_field": "hello" }`
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, bytes.NewBufferString(body))
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with an error about the unknown field", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, `json: unknown field "extra_field"`)
			})
		})

		Convey("When the field type provided is not the expected and the CreateOrUpdateCacheTime endpoint is called", func() {
			body := `{"path": 1234, "etag": "testetag"}`
			request := httptest.NewRequest(http.MethodPut, baseURL+testCacheID, bytes.NewBufferString(body))
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with a type mismatch error in the response", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, "json: cannot unmarshal number into Go struct field CacheTime.path of type string")
			})
		})

		Convey("When the id provided is not 32 characters in length and the CreateOrUpdateCacheTime endpoint is called", func() {
			body := validBody
			request := httptest.NewRequest(http.MethodPut, baseURL+"abc", bytes.NewBufferString(body))
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with an ID length error in the response", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, "[id should be 32 characters in length]")
			})
		})

		Convey("When the id provided is not lowercase and the CreateOrUpdateCacheTime endpoint is called", func() {
			idWithUpperCase := "1A2B3C4D5E6F7890A1B2C3D4E5F67890"
			body := validBody
			request := httptest.NewRequest(http.MethodPut, baseURL+idWithUpperCase, bytes.NewBufferString(body))
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with an ID format error in the response", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, "[id is not lowercase]")
			})
		})
		Convey("When the id provided is not hexadecimal CreateOrUpdateCacheTime endpoint is called", func() {
			idNotHexadecimal := "1a2b3c4d5g6h7890g1h2i3j4k5l67890"
			body := validBody
			request := httptest.NewRequest(http.MethodPut, baseURL+idNotHexadecimal, bytes.NewBufferString(body))
			responseRecorder := httptest.NewRecorder()
			dataStoreAPI.Router.ServeHTTP(responseRecorder, request)
			Convey("Then a 400 is returned with an ID format error in the response", func() {
				So(responseRecorder.Code, ShouldEqual, 400)
				So(responseRecorder.Body.String(), ShouldContainSubstring, "[id is not a valid hexadecimal]")
			})
		})
	})
}

func createRequestWithAuth(method, target string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, target, body)
	ctx := r.Context()
	ctx = dprequest.SetCaller(ctx, "someone@ons.gov.uk")
	r = r.WithContext(ctx)
	return r
}
