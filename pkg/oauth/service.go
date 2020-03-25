package oauth

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ImportUsers(ctx context.Context, fields []Fields) error
	Authorize(ctx context.Context, redirectURI string) (uri string, err error)
}

type service struct {
	config *Config
	userCollection UserCollection
}

func NewService(config *Config, userCollection UserCollection) Service {
	return service{config: config, userCollection: userCollection}
}

func (s service) ImportUsers(ctx context.Context, fieldsArr []Fields) error {
	var users []*User
	Outer:
	for _, fields := range fieldsArr {
		user := User{Fields: make(Fields)}
		for _, spec := range s.config.Fields {
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
	id, _ := primitive.ObjectIDFromHex("5e76824e6a9946d454b731c5")
	user, err := s.userCollection.Find(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?code=%s", redirectURI, user.ID.Hex()), nil
}