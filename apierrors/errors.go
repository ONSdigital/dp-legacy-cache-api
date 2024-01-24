package apierrors

import (
	"errors"
)

// A list of error messages for Dataset API
var (
	ErrCacheTimeNotFound = errors.New("cachetime not found")
	ErrDataStore         = errors.New("DataStore error")
)
