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
		// TODO: put hardcoded strings in env file
		userCollection := oauth.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
		config, err := oauth.LoadConfig("config_example.yml", ".")
		if err != nil {
			return err
		}
		service := oauth.NewService(config, userCollection)

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
		return service.ImportUsers(context.TODO(), users)
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}