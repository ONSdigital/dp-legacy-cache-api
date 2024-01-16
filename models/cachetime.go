package models

import "time"

type CacheTime struct {
	ID           string    `json:"_id" bson:"_id"`                     // MD5 of the path
	Path         string    `json:"path" bson:"path"`                   // Path for which caching is set
	ETag         string    `json:"etag" bson:"etag"`                   // ETag for cache validation
	CollectionID int       `json:"collection_id" bson:"collection_id"` // Collection ID - used for grouping and filtering of cache-time objects.
	ReleaseTime  time.Time `json:"release_time" bson:"release_time"`   // Release time in ISO-8601 format
}
