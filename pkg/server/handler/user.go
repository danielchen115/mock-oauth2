package handler

import (
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"net/http"
)

func UserInfo(service oauth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

	}
}