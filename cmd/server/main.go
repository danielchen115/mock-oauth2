package main

import (
	"fmt"
	"github.com/danielchen115/mock-oauth2/pkg/oauth"
	"github.com/danielchen115/mock-oauth2/pkg/server"
	"os"
)

func main() {
	err := func() error {
		config, err := oauth.LoadConfig("config.yml", ".")
		if err != nil {
			return err
		}
		dbConfig := config.Database
		dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%d", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port)
		userCollection := oauth.NewUserCollection(dbURI, dbConfig.Database)
		serverURI := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
		svc := oauth.NewService(config, userCollection)
		s := server.New(svc, serverURI)
		fmt.Printf("listening on %s...\n", serverURI)
		return s.ListenAndServe()
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}