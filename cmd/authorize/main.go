package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alfrye/authorize/internal/authorization/provider"
	"github.com/alfrye/authorize/internal/authorize"
	api "github.com/alfrye/authorize/internal/handlers/api"
	persistence "github.com/alfrye/authorize/internal/persistence/mongo"
	"github.com/alfrye/authorize/internal/server"
)

func main() {

	fmt.Println("starting point for Authorize")
	s := server.New("9010")
	repo := choseRepository()
	authProvider := choseAuthProvider()
	authService := authorize.NewAuthService(repo, authProvider)
	nhandler := api.NewAuthHandler(authService)
	s.PopulateRoutes(s.AuthorizeServiceRoutes(nhandler))
	s.Listen()

}

func choseRepository() authorize.AuthorizeRepository {

	switch os.Getenv("DB_URL") {
	case "mongo":
		//Set up information for mongo db
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
		if err != nil {

		}
		repo, err := persistence.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal("Could not create mongo repo")
		}
		return repo
	}

	return nil

}

func choseAuthProvider() authorize.AuthProvider {
	switch os.Getenv("AUTH_PROVIDER") {
	case "local":
		authProvider, err := provider.NewLocalAuthProvider()
		if err != nil {
			log.Fatal("Can not create authentication provider")
		}
		return authProvider
	}

	return nil
}
