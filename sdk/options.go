package sdk

import "time"

type Options struct {
	Limit       int
	Offset      int
	ReleaseTime time.Time
}
