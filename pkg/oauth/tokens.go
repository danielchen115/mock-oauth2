package oauth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AccessToken = jwt.Token
type RefreshToken = jwt.Token

func NewAccessToken(config TokenConfig, user *User) *AccessToken {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID.Hex()
	claims["cid"] = config.ClientID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Second * time.Duration(config.AccessTokenDuration)).Unix()
	claims["scope"] = user.Scope
	return token
}

func NewRefreshToken(config TokenConfig, user *User) *RefreshToken {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID.Hex()
	claims["cid"] = config.ClientID
	claims["cse"] = config.ClientSecret
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Second * time.Duration(config.RefreshTokenDuration)).Unix()
	claims["scope"] = user.Scope
	return token
}