package handler

import (
	"errors"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Authorize(service oauth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := func() error {
			if req.Method != http.MethodGet {
				return errors.New("method not allowed")
			}
				redirectURI := req.URL.Query().Get("redirect_uri")
			if redirectURI == "" {
				return errors.New("redirect URI not found")
			}
			uri, err := service.Authorize(req.Context(), redirectURI)
			if err != nil {
				return err
			}
			http.Redirect(res, req, uri, http.StatusTemporaryRedirect)
			return nil
		} (); err != nil {
			status := http.StatusInternalServerError
			switch err {
			case mongo.ErrNoDocuments:
					status = http.StatusNotFound
			}
			http.Error(res, err.Error(), status)
		}
	}
}