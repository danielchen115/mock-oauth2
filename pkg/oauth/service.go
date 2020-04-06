package oauth

import (
	"context"
	"fmt"
)

type Service interface {
	ImportUsers(ctx context.Context, fields []Fields) error
	Authorize(ctx context.Context, redirectURI string) (uri string, err error)
	SetCurrentUser(ctx context.Context, id string) (user *User, err error)
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
	user, err := s.userCollection.Find(ctx, s.currentUserID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?code=%s", redirectURI, user.ID), nil
}

func (s service) SetCurrentUser(ctx context.Context, id string) (user *User, err error) {
	s.currentUserID = id
	return s.userCollection.Find(ctx, s.currentUserID)
}