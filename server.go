package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/danielchen115/mock-oauth2/pkg/db"
)

func main() {
	// user connection setup
    userCollection := db.NewUserCollection("mongodb://root:secret@mongodb:27017", "mock-oauth2")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:secret@mongodb:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// example usage
	collection := client.Database("mock-oauth2").Collection("users")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"name": "test", "token": "test_hash"})
	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Println("Inserted document:", id)

	//web server
	http.HandleFunc("/", ServerHandler)
	http.ListenAndServe(":7080", nil)
}

func ServerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mock OAuth 2 Server")
}
