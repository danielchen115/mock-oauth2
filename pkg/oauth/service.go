package oauth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Service interface {
	ImportUsers(ctx context.Context, fields []Fields) error
	Authorize(ctx context.Context, redirectURI string) (uri string, err error)
	SetCurrentUser(ctx context.Context, id primitive.ObjectID) (user *User, err error)
	Tokens(ctx context.Context, params TokenParam) (at string, rt string, err error)
	GetAccessTokenDuration() int
	ExtractClaims(tokenString string) (*RefreshClaims, error)
	RefreshTokenValid(token string, params TokenParam) error
}

type TokenParam struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	GrantType    string
	Code         string
	RefreshToken string
}

type RefreshClaims struct {
	jwt.StandardClaims
	ClientID     string `json:"cid,omitempty"`
	ClientSecret string `json:"cse,omitempty"`
}

type service struct {
	config         *Config
	userCollection UserCollection
	currentUserID  primitive.ObjectID
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
	if params.GrantType == "authorization_code" {
		return s.parseAuthorize(ctx, params)
	}
	if params.GrantType == "refresh_token" {
		return s.parseRefresh(ctx, params)
	}
	return "", "", errors.New("unsupported grant type")
}

func (s *service) parseAuthorize(ctx context.Context, params TokenParam) (at string, rt string, err error) {
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

func (s *service) parseRefresh(ctx context.Context, params TokenParam) (at string, rt string, err error) {
	if params.ClientID != s.config.Token.ClientID {
		return "", "", errors.New("unknown client ID")
	}
	if params.ClientSecret != s.config.Token.ClientSecret {
		return "", "", errors.New("invalid client secret")
	}
	err = s.RefreshTokenValid(params.RefreshToken, params)
	if err != nil {
		return "", "", err
	}
	user, err := s.userCollection.Find(ctx, s.currentUserID)
	if err != nil {
		return "", "", err
	}
	err = user.AddRefreshToken(s.config.Token)
	return user.AccessToken, user.RefreshToken, err
}

func (s *service) GetAccessTokenDuration() int {
	return s.config.Token.AccessTokenDuration
}

func (s *service) ExtractClaims(tokenString string) (*RefreshClaims, error) {
	var claims RefreshClaims
	signingSecret := []byte(s.config.Token.SigningSecret)
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return signingSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid JWT token")
}

func (s *service) RefreshTokenValid(token string, params TokenParam) error {
	claims, err := s.ExtractClaims(token)
	if err != nil {
		return err
	}
	var expTime time.Time
	expClaim := claims.ExpiresAt
	expTime = time.Unix(expClaim, 0)
	if time.Now().After(expTime) {
		return errors.New("token expired")
	}
	if claims.ClientID != params.ClientID {
		return errors.New("invalid client id")
	}
	if claims.ClientSecret != params.ClientSecret {
		return errors.New("invalid client secret")
	}
	return nil
}
