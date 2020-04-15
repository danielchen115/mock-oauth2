package oauth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ImportUsers(ctx context.Context, fields []Fields) error
	Authorize(ctx context.Context, redirectURI string) (uri string, err error)
	SetCurrentUser(ctx context.Context, id primitive.ObjectID) (user *User, err error)
	Tokens(ctx context.Context, params TokenParam) (at string, rt string, err error)
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
	currentUserID primitive.ObjectID
}

func NewService(config *Config, userCollection UserCollection) Service {
	return &service{config: config, userCollection: userCollection}
}

func (s *service) ImportUsers(ctx context.Context, fieldsArr []Fields) error {
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

func (s *service) Authorize(ctx context.Context, redirectURI string) (uri string, err error) {
	id := s.currentUserID
	_, err = s.userCollection.Find(ctx, id)
	if err != nil {
		return "", errors.New("user not found")
	}
	code := base64.URLEncoding.EncodeToString([]byte(id.Hex()))
	return fmt.Sprintf("%s?code=%s", redirectURI, code), nil
}

func (s *service) SetCurrentUser(ctx context.Context, id primitive.ObjectID) (user *User, err error) {
	s.currentUserID = id
	return s.userCollection.Find(ctx, id)
}

func (s *service) Tokens(ctx context.Context, params TokenParam) (at string, rt string, err error) {
	if params.GrantType != "authorization_code" {
		return "", "", errors.New("unsupported grant type")
	}
	if params.RedirectURI == "" {
		return "", "", errors.New("redirect URI not found")
	}
	if params.ClientID != s.config.Token.ClientID {
		return "", "", errors.New("unknown client ID")
	}
	code, _ := base64.URLEncoding.DecodeString(params.Code)
	if string(code) != s.currentUserID.Hex() {
		return "", "", errors.New("authorization code invalid")
	}
	user, err := s.userCollection.Find(ctx, s.currentUserID)
	if err != nil {
		return "", "", err
	}
	err = user.AddAccessToken(s.config.Token)
	if err != nil {
		return "", "", err
	}
	err = user.AddRefreshToken(s.config.Token)
	return user.AccessToken, user.RefreshToken, err
}

func (s *service) GetAccessTokenDuration() int {
	return s.config.Token.AccessTokenDuration
}