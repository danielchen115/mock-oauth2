package oauth_test

import (
	"context"
	"fmt"
)

func ExampleInsertUser() {
	userCollection := NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user := User{AccessToken: "1", Fields: Fields{"name": "tom"}}
	userCollection.Insert(context.TODO(), &user)
	fmt.Println(user.AccessToken)
	// Output: 1
}

func ExampleInsertManyUsers() {
	userCollection := NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user1 := User{AccessToken: "1", Fields: Fields{"name": "tom"}}
	user2 := User{AccessToken: "2", Fields: Fields{"name": "dave"}}
	userCollection.InsertMany(context.TODO(), []*User{&user1, &user2})
	fmt.Println(user1.AccessToken)
	fmt.Println(user2.AccessToken)
	// Output:
	// 1
	// 2
}