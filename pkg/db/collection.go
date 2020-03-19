package db

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type userCollection struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type User struct {
	mandatoryFields map[string]interface{}
	accessToken     string
	refreshToken    string
	metadata        json.RawMessage
}

type Users struct {
	user []User
}

func NewUserCollection(URI string, database string) *userCollection {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	collection := client.Database(database).Collection("users")
	if err != nil {
		log.Fatal(err)
	}
	uc := userCollection{client, collection}
	return &uc
}