package models

import "time"

type CacheTime struct {
	ID           string    `json:"_id" bson:"_id" validate:"required"`   // MD5 of the path, required
	Path         string    `json:"path" bson:"path" validate:"required"` // Path for which caching is set, required
	ETag         string    `json:"etag" bson:"etag" validate:"required"` // ETag for cache validation, required
	CollectionID int       `json:"collection_id" bson:"collection_id"`   // Collection ID - used for grouping and filtering of cache-time objects.
	ReleaseTime  time.Time `json:"release_time" bson:"release_time"`     // Release time in ISO-8601 format
}
