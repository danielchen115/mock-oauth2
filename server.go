package main

import (
    "fmt"
    "net/http"
    "context"
    "time"
    "log"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

func main() {
    // mongodb setup
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:secret@mongodb:27017"))
    fmt.Println("hit")
    if err != nil {
        log.Fatal(err)
    }

    // example usage
    collection := client.Database("mock-oauth2").Collection("users")
    //ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
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
