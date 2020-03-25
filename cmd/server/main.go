package main

import (
	"fmt"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"github.com/danielchen115/mock-oauth2/pkg/server"
	"os"
)

const addr = "localhost:8090"

func main() {
	err := func() error {
		// TODO: put hardcoded strings in env file
		userCollection := oauth.NewUserCollection("mongodb://root:secret@localhost:27017", "mock-oauth2")
		config, err := oauth.LoadConfig("config_example.yml", ".")
		if err != nil {
			return err
		}
		svc := oauth.NewService(config, userCollection)
		s := server.New(svc, addr)
		fmt.Printf("listening on %s...\n", addr)
		return s.ListenAndServe()
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}