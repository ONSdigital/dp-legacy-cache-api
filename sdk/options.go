package sdk

import "time"

type Options struct {
	ServiceAuthToken string
	UserAccessToken  string
	Limit            int
	Offset           int
	ReleaseTime      time.Time
}
