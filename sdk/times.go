package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-legacy-cache-api/api"
	"github.com/ONSdigital/dp-legacy-cache-api/models"
	apiError "github.com/ONSdigital/dp-legacy-cache-api/sdk/errors"
)

// GetCacheTimes gets a list of cache times
func (cli *Client) GetCacheTimes(ctx context.Context, auth Auth, opts Options) (*models.CacheTimesList, apiError.Error) {
	path := fmt.Sprintf("%s/cache-times", cli.hcCli.URL)
	var cacheTimesList models.CacheTimesList

	queryParams := url.Values{}
	if !opts.ReleaseTime.IsZero() {
		queryParams.Add(api.QueryParamReleaseTime, opts.ReleaseTime.Format(time.RFC3339))
	}

	if opts.Limit > 0 {
		queryParams.Add(api.QueryParamLimit, strconv.Itoa(opts.Limit))
	}

	if opts.Offset > 0 {
		queryParams.Add(api.QueryParamOffset, strconv.Itoa(opts.Offset))
	}

	respInfo, apiErr := cli.callLegacyCacheAPI(ctx, path, http.MethodGet, auth, queryParams, nil)
	if apiErr != nil {
		return &cacheTimesList, apiErr
	}

	if err := json.Unmarshal(respInfo.Body, &cacheTimesList); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal CacheTimesList response - error is: %v", err),
		}
	}

	return &cacheTimesList, nil
}
