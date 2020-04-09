package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func Token(service oauth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := func() error {
			if req.Method != http.MethodPost {
				return errors.New("method not allowed")
			}
			tokenParams := oauth.TokenParam{
				ClientID:    req.URL.Query().Get("client_id"),
				RedirectURI: req.URL.Query().Get("redirect_uri"),
				GrantType:   req.URL.Query().Get("grant_type"),
				Code:        req.URL.Query().Get("code"),
			}
			token, err := service.Token(context.TODO(), tokenParams)
			if err != nil {
				return err
			}
			type Response struct {
				AccessToken  string `json:"access_token,omitempty"`
				TokenType    string `json:"token_type,omitempty"`
				ExpiresIn    int    `json:"expires_in,omitempty"`
				RefreshToken string `json:"refresh_token,omitempty"`
				Scope        string `json:"scope,omitempty"`
			}
			enc := json.NewEncoder(res)
			return enc.Encode(Response{
				AccessToken:  token,
				TokenType:    "bearer",
				ExpiresIn:    service.GetAccessTokenDuration(),

			})
		}(); err != nil {
			status := http.StatusInternalServerError
			switch err {
			case mongo.ErrNoDocuments:
				status = http.StatusNotFound
			}
			http.Error(res, err.Error(), status)
		}
	}
}
