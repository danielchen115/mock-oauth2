package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"os"
)

func main() {
	err := func() error {
		if len(os.Args) < 2 {
			return errors.New("expects 1 argument, got none")
		}
		config, err := oauth.LoadConfig("config_example.yml", ".")
		if err != nil {
			return err
		}
		dbConfig := config.Database
		dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%d", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port)
		userCollection := oauth.NewUserCollection(dbURI, dbConfig.Database)
		svc := oauth.NewService(config, userCollection)
		f, err := os.Open(os.Args[1])
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(f)
		var users []oauth.Fields
		err = decoder.Decode(&users)
		if err != nil {
			return err
		}
		return svc.ImportUsers(context.TODO(), users)
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}