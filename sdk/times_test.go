package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-legacy-cache-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testCacheTime1 = models.CacheTime{
		CollectionID: "collection1",
		ID:           "test-cache-time-1",
		Path:         "/test/path/1",
		ReleaseTime:  nil,
	}
	testCacheTime2 = models.CacheTime{
		CollectionID: "collection2",
		ID:           "test-cache-time-2",
		Path:         "/test/path/2",
		ReleaseTime:  nil,
	}

	testCacheTimeItems = []*models.CacheTime{&testCacheTime1, &testCacheTime2}
	testCacheTimesList = models.CacheTimesList{
		Count:  2,
		Items:  testCacheTimeItems,
		Offset: 0,
		Limit:  20,
	}

	testPaginationCacheTimeItems = []*models.CacheTime{&testCacheTime1}
	testPaginationCacheTimesList = models.CacheTimesList{
		Count:  1,
		Items:  testPaginationCacheTimeItems,
		Offset: 0,
		Limit:  20,
	}
)

func TestGetCacheTimes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given legacy cache API returns successfully", t, func() {
		body, err := json.Marshal(testCacheTimesList)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		legacyCacheAPIClient := newLegacyCacheAPIClient(t, httpClient)

		Convey("When GetCacheTimes is called", func() {
			cacheTimesList, err := legacyCacheAPIClient.GetCacheTimes(ctx, Auth{}, Options{})

			Convey("Then the expected cache times are returned", func() {
				So(*cacheTimesList, ShouldResemble, testCacheTimesList)
				So(cacheTimesList.Count, ShouldEqual, 2)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/cache-times")
					})
				})
			})
		})

		Convey("When GetCacheTimes is called with pagination params", func() {
			body, err := json.Marshal(testPaginationCacheTimesList)
			if err != nil {
				t.Errorf("failed to setup test data, error: %v", err)
			}

			httpClient := newMockHTTPClient(
				&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				},
				nil)

			legacyCacheAPIClient := newLegacyCacheAPIClient(t, httpClient)
			cacheTimesList, err := legacyCacheAPIClient.GetCacheTimes(ctx, Auth{}, Options{Limit: 1, Offset: 0})

			Convey("Then the expected cache times are returned", func() {
				So(*cacheTimesList, ShouldResemble, testPaginationCacheTimesList)
				So(cacheTimesList.Count, ShouldEqual, 1)
				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/cache-times")
					})
				})
			})
		})
	})

	Convey("When GetCacheTimes is called with no response body returned", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			nil)

		legacyCacheAPIClient := newLegacyCacheAPIClient(t, httpClient)
		cacheTimesList, err := legacyCacheAPIClient.GetCacheTimes(ctx, Auth{}, Options{})

		Convey("Then no cache times are returned", func() {
			So(cacheTimesList, ShouldBeNil)

			Convey("And no error is returned", func() {
				So(err, ShouldNotBeNil)

				Convey("And client.Do should be called once with the expected parameters", func() {
					doCalls := httpClient.DoCalls()
					So(doCalls, ShouldHaveLength, 1)
					So(doCalls[0].Req.URL.Path, ShouldEqual, "/cache-times")
				})
			})
		})
	})

	Convey("Given a request is made with invalid query parameters", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusBadRequest}, nil)
		legacyCacheAPIClient := newLegacyCacheAPIClient(t, httpClient)

		Convey("When GetCacheTimes is called", func() {
			cacheTimesList, err := legacyCacheAPIClient.GetCacheTimes(ctx, Auth{}, Options{Limit: -1})

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusBadRequest)

				Convey("And the expected cache list should be empty", func() {
					So(cacheTimesList, ShouldEqual, &models.CacheTimesList{})

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
					})
				})
			})
		})
	})

	Convey("Given a 500 response from legacy cache API", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		legacyCacheAPIClient := newLegacyCacheAPIClient(t, httpClient)

		Convey("When GetCacheTimes is called", func() {
			cacheTimesList, err := legacyCacheAPIClient.GetCacheTimes(ctx, Auth{}, Options{})

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected cache times list should be empty", func() {
					So(cacheTimesList.Items, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/cache-times")
					})
				})
			})
		})
	})
}
