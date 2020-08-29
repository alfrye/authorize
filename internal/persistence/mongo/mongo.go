package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	// Database struct {
	// 	Conn   *Connection
	// 	Client *mongo.Client
	// }

	mongoRepository struct {
		client   *mongo.Client
		database string
		timeout  time.Duration
	}

	// Connection struct {
	// 	Name string
	// 	Host string
	// 	Port string
	// }
)

func newMongoCient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil

}

// NewMongoRepository instatiates a a mongo repository
func NewMongoRepository(mongoURL, mongoDB string, timeout int) (authorize.AuthorizeRepository, error) {
	repo := &mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(timeout) * time.Second,
	}

	client, err := newMongoCient(mongoURL, timeout)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	repo.client = client

	return repo, nil
}

// func (db *Database) Connect() *Database {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	db.Client = client
// 	return db

// }

// CreateUser creates a new user in the database
func (r *mongoRepository) CreateUser(user models.Users) error {
	collection := r.client.Database(r.database).Collection("users")

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Inserted record", res)
	return nil
}

// GetUser retrieves a user from the database
func (r *mongoRepository) GetUser(userName string) (models.Users, error) {
	var user models.Users
	collection := r.client.Database(r.database).Collection("users")
	filter := bson.M{"name": userName}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return models.Users{}, err
	}
	return user, nil
}
