package oauth_test

import (
	"context"
	"fmt"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
)

func ExampleInsertUser() {
	userCollection := oauth.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user := oauth.User{AccessToken: "1", Fields: oauth.Fields{"name": "tom"}}
	userCollection.Insert(context.TODO(), &user)
	fmt.Println(user.AccessToken)
	// Output: 1
}

func ExampleInsertManyUsers() {
	userCollection := oauth.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user1 := oauth.User{AccessToken: "1", Fields: oauth.Fields{"name": "tom"}}
	user2 := oauth.User{AccessToken: "2", Fields: oauth.Fields{"name": "dave"}}
	userCollection.InsertMany(context.TODO(), []*oauth.User{&user1, &user2})
	fmt.Println(user1.AccessToken)
	fmt.Println(user2.AccessToken)
	// Output:
	// 1
	// 2
}

func ExampleFindOneUser() {
	userCollection := oauth.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
	user1 := oauth.User{AccessToken: "1", Fields: oauth.Fields{"name": "tom"}}
	userCollection.Insert(context.TODO(), &user1)
	user, _ := userCollection.Find(context.TODO(), user1.ID)
	fmt.Println(user.Fields["name"])
	// Output:
	// tom
}