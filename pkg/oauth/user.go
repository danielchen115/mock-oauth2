package oauth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type UserCollection interface {
	Insert(ctx context.Context, user *User) error
	InsertMany(ctx context.Context, users []*User) error
	Find(ctx context.Context, id string) (*User, error)
}

type Fields map[string]interface{}

type userCollection struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Fields       Fields             `json:"fields"`
	Scope        []string           `json:"scope" bson:"scope,omitempty"`
	AccessToken  string             `json:"accessToken" bson:"accessToken,omitempty"`
	RefreshToken string             `json:"refreshToken" bson:"refreshToken,omitempty"`
}

func NewUserCollection(URI string, database string) UserCollection {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	collection := client.Database(database).Collection("users")
	if err != nil {
		log.Fatal(err)
	}
	return userCollection{client, collection}
}

func (uc userCollection) Insert(ctx context.Context, user *User) error {
	res, err := uc.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(string)
	return nil
}

func (uc userCollection) InsertMany(ctx context.Context, users []*User) error {
	objects := make([]interface{}, len(users))
	for i, user := range users {
		objects[i] = user
	}
	res, err := uc.collection.InsertMany(ctx, objects)
	if err != nil {
		return err
	}
	for i, id := range res.InsertedIDs {
		users[i].ID = id.(string)
	}
	return nil
}

func (uc userCollection) Find(ctx context.Context, id string) (*User, error) {
	var user User
	err := uc.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
