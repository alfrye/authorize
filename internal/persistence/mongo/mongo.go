package persistence

import (
	"context"
	"errors"
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

func (r *mongoRepository) GetAllUsers() ([]models.Users, error) {
	var users []models.Users
	collection := r.client.Database(r.database).Collection("users")
	filter := bson.D{{}}
	cur, err := collection.Find(context.Background(), filter, options.Find()) //.Decode(&users)
	if err != nil {
		fmt.Println(err)
		return nil, err

	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var user models.Users
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		// results =append(results, elem)
		users = append(users, user)
	}

	return users, nil

}

func (r *mongoRepository) GetClient(clientID string) (models.Client, error) {
	return models.Client{}, errors.New("not implemented")
}

func (r *mongoRepository) CreateClient(client models.Client) error {
	return errors.New("not implemented")
}

func (r *mongoRepository) SaveAuthCode(code models.AuthCode) error {
	return errors.New("not implemented")
}

func (r *mongoRepository) GetAuthCode(code string) (models.AuthCode, error) {
	return models.AuthCode{}, errors.New("not implemented")
}

func (r *mongoRepository) DeleteAuthCode(code string) error {
	return errors.New("not implemented")
}

func (r *mongoRepository) CreateSession(s models.Session) error {
	return errors.New("not implemented")
}

func (r *mongoRepository) GetSession(sessionID string) (models.Session, error) {
	return models.Session{}, errors.New("not implemented")
}

func (r *mongoRepository) DeleteSession(sessionID string) error {
	return errors.New("not implemented")
}

func (r *mongoRepository) SaveRefreshToken(rt models.RefreshToken) error {
	return errors.New("not implemented")
}

func (r *mongoRepository) GetRefreshToken(token string) (models.RefreshToken, error) {
	return models.RefreshToken{}, errors.New("not implemented")
}

func (r *mongoRepository) DeleteRefreshToken(token string) error {
	return errors.New("not implemented")
}
