package main

import (
	"fmt"

	"github.com/alfrye/authorize/cmd"
)

func main() {
	fmt.Println("Starting Authorize API with go-kit")
	cmd.Execute()
}

// func choseRepository() authorize.AuthorizeRepository {
//
// 	switch os.Getenv("DB_URL") {
// 	case "mongo":
// 		//Set up information for mongo db
// 		mongoURL := os.Getenv("MONGO_URL")
// 		mongoDB := os.Getenv("MONGO_DB")
// 		mongoTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
// 		if err != nil {
//
// 		}
// 		repo, err := mg.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
// 		if err != nil {
// 			log.Fatal("Could not create mongo repo")
// 		}
// 		return repo
// 	case "mysql":
// 		//Setup information for mysql database
// 		//	mySQLURL := os.Getenv("MYSQL_URL")
// 		mySQLHost := os.Getenv("MYSQL_HOST")
// 		mySQLPort := os.Getenv("MYSQL_PORT")
// 		mySQLDatabase := os.Getenv("MYSQL_DB")
// 		mySQLUser := os.Getenv("MYSQL_USER")
// 		mySQLPassword := os.Getenv("MYSQL_PASS")
//
// 		mySQLTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
// 		if err != nil {
//
// 		}
//
// 		repo, err := mysql.NewMySQLRepository("mysql", mySQLHost, mySQLPort, mySQLDatabase, mySQLUser, mySQLPassword, mySQLTimeout)
// 		if err != nil {
// 			log.Fatal("Could not create mongo repo")
// 		}
// 		return repo
//
// 	}
//
// 	return nil
//
// }


