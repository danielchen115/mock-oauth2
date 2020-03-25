package handler

import (
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"net/http"
)

func Authorize(service oauth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		redirectURI := req.URL.Query().Get("redirect_uri")
		uri, err := service.Authorize(req.Context(), redirectURI)
		if err != nil {
			panic(err)
		}
		http.Redirect(res, req, uri, http.StatusTemporaryRedirect)

	}
}