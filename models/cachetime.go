package models

import "time"

type CacheTime struct {
	ID           string     `bson:"_id" json:"_id"`                                         // MD5 of the path
	Path         string     `bson:"path" json:"path"`                                       // Path for which caching is set
	CollectionID string     `bson:"collection_id,omitempty" json:"collection_id,omitempty"` // Collection ID - used for grouping and filtering of cache-time objects.
	ReleaseTime  *time.Time `bson:"release_time,omitempty" json:"release_time,omitempty"`   // Release time in ISO-8601 format
}

type CacheTimesList struct {
	Items      *[]CacheTime `json:"items"`
	Count      int          `json:"count"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
	TotalCount int          `json:"total_count"`
}
