package persistence

import (
	"context"
	"fmt"
	"log"

	"github.com/alfrye/authorize/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Database struct {
		Conn   *Connection
		Client *mongo.Client
	}

	Connection struct {
		Name string
		Host string
		Port string
	}
)

func (db *Database) Connect() *Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	db.Client = client
	return db

}

// CreateUser creates a new user in the database
func (db *Database) CreateUser(user models.Users) {
	collection := db.Client.Database("security").Collection("users")

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Inserted record", res)
}

// GetUser retrieves a user from the database
func (db *Database) GetUser(userName string) models.Users {
	var user models.Users
	collection := db.Client.Database("security").Collection("users")
	filter := bson.D{{"name", userName}}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
	}
	return user
}
