package oauth

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ImportUsers(ctx context.Context, fields []Fields) error
	Authorize(ctx context.Context, redirectURI string) (uri string, err error)
	SetCurrentUser(ctx context.Context, id string) (user *User, err error)
	Token(ctx context.Context, params TokenParam) (string, error)
	GetAccessTokenDuration() int
}

type TokenParam struct {
	ClientID string
	RedirectURI string
	GrantType string
	Code string
}

type service struct {
	config *Config
	userCollection UserCollection
	currentUserID string
}

func NewService(config *Config, userCollection UserCollection) Service {
	return service{config: config, userCollection: userCollection}
}

func (s service) ImportUsers(ctx context.Context, fieldsArr []Fields) error {
	var users []*User
	Outer:
	for _, fields := range fieldsArr {
		user := User{Fields: make(Fields)}
		for _, spec := range s.config.Import.Fields {
			value, ok := fields[spec.Name]
			if !ok && spec.Required {
				fmt.Printf("Required field \"%s\" not found, user not imported.\n", spec.Name)
				continue Outer
			}
			user.Fields[spec.Name] = value
		}
		users = append(users, &user)
	}
	return s.userCollection.InsertMany(ctx, users)
}

func (s service) Authorize(ctx context.Context, redirectURI string) (uri string, err error) {
	id, _ := primitive.ObjectIDFromHex(s.currentUserID)
	user, err := s.userCollection.Find(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?code=%s", redirectURI, user.ID.Hex()), nil
}

func (s service) SetCurrentUser(ctx context.Context, id string) (user *User, err error) {
	s.currentUserID = id
	hexID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.userCollection.Find(ctx, hexID)
}

func (s service) Token(ctx context.Context, params TokenParam) (string, error) {
	if params.GrantType != "authorization_code" {
		return "", errors.New("unsupported grant type")
	}
	if params.RedirectURI == "" {
		return "", errors.New("redirect URI not found")
	}
	if params.ClientID != s.config.Token.ClientID {
		return "", errors.New("unknown client ID")
	}
	if params.Code != s.currentUserID {
		return "", errors.New("authorization code invalid")
	}
	userID, err := primitive.ObjectIDFromHex(s.currentUserID)
	if err != nil {
		return "", err
	}
	user, err := s.userCollection.Find(ctx, userID)
	if err != nil {
		return "", err
	}
	err = user.AddAccessToken(s.config.Token)
	return user.AccessToken, err
}

func (s service) GetAccessTokenDuration() int {
	return s.config.Token.AccessTokenDuration
}