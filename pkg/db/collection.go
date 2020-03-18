package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCollection struct {
	client mongo.Client
	collection mongo.Collection
}

type User struct {
	mandatoryFields map[string]interface{}
	accessToken string
	refreshToken string
	metadata map[string]interface{}
}

type Users struct {
	user []User
}