package handler

import (
	"context"
	"encoding/json"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func CurrentUser(service oauth.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := func() error {
			switch req.Method {
			case http.MethodPost:
				var payload struct {
					ID string `json:"id"`
				}
				if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
					return err
				}
				objectID, err := primitive.ObjectIDFromHex(payload.ID)
				if err != nil {
					return err
				}
				user, err := service.SetCurrentUser(context.TODO(), objectID)
				if err != nil {
					return err
				}
				return json.NewEncoder(res).Encode(user)
			}
			return nil
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
