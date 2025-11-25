package sdk

import (
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/v3/request"
)

type Auth struct {
	ServiceAuthToken string `json:"-"`
	UserAccessToken  string `json:"-"`
}

func (a *Auth) Add(req *http.Request) {
	if a.ServiceAuthToken != "" {
		dprequest.AddServiceTokenHeader(req, a.ServiceAuthToken)
	}

	if a.UserAccessToken != "" {
		dprequest.AddFlorenceHeader(req, a.UserAccessToken)
	}
}
