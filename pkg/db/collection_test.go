package db_test

import (
	"context"
	"fmt"
	"github.com/danielchen115/mock-oauth2/pkg/db"
)

func ExampleInsertUser() {
	userCollection := db.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user := db.User{AccessToken: "1", Fields: db.Fields{"name":"tom"}}
	userCollection.Insert(context.TODO(), &user)
	fmt.Println(user.AccessToken)
	// Output: 1
}

func ExampleInsertManyUsers() {
	userCollection := db.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user1 := db.User{AccessToken: "1", Fields: db.Fields{"name":"tom"}}
	user2 := db.User{AccessToken: "2", Fields: db.Fields{"name":"dave"}}
	userCollection.InsertMany(context.TODO(), []*db.User{&user1, &user2})
	fmt.Println(user1.AccessToken)
	fmt.Println(user2.AccessToken)
	// Output:
	// 1
	// 2
}