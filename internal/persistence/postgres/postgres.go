package persistence

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
	"github.com/jackc/pgx/v5"
)

type (
	postgresRepository struct {
		client   *pgx.Conn
		database string
		timeout  time.Duration
	}
)

func newPostgresClient(postgresURL string ) (*pgx.Conn, error) {

	dbClient, err := pgx.Connect(context.Background(), "databaseURL")
	if err != nil {
		fmt.Println("Can not open database")
		return nil, err

	}

	return dbClient, nil
}

func NewPostgresRepository(dbHost, dbPort, postgresDB, postgresUser, postgresPassword string, timeout int) (authorize.AuthorizeRepository, error) {
	dbAddr := net.JoinHostPort(dbHost, dbPort)
	dbName := "authorize"
	postgresURL := "postgresql://" + postgresUser + ":" + postgresPassword + "@" + dbAddr + "/" + dbName
	// postgresConfig := postgres.Config{
	// 	User:   postgresUser,
	// 	Passwd: postgresPassword,
	// 	Net:    "tcp",
	// 	DBName: postgresDB,
	// 	Addr:   dbAddr,
	// }

	dbClient, err := newPostgresClient(postgresURL)
	if err != nil {
		fmt.Printf("Error Creating client: %s", err)
		return nil, err
	}
	repo := &postgresRepository{
		client:   dbClient,
		database: postgresDB,
		timeout:  time.Duration(timeout) * time.Second,
	}

	return repo, nil
}

func (r *postgresRepository) CreateUser(user models.Users) error {

	// insUser, err := r.client.Query("INSERT into users(name, password,email) Values(?,?,?)")
	// if err != nil {
	// 	fmt.Printf("Error creating sql: %v", err)
	// 	return err
	// }
	//
	// res, err := insUser.Exec(user.Name, user.Password, user.Email)
	// if err != nil {
	// 	fmt.Printf("Error creating sql: %v", err)
	// 	return err
	// }
	//
	fmt.Println("Calling Creat User")

	return nil
}

func (r *postgresRepository) GetUser(username string) (models.Users, error) {
	var user models.Users
	// selUser, err := r.client.Query("SELECT name,password,email from users Where name=?", username)
	// if err != nil {
	// 	fmt.Printf("Error creating sql: %v", err)
	// 	return user, err
	// }
	// for selUser.Next() {
	// 	var userName, userPassword, userEmail string
	//
	// 	err := selUser.Scan(&userName, &userPassword, &userEmail)
	// 	if err != nil {
	// 		fmt.Printf("Error scanning results:%v", err)
	// 		return models.Users{}, err
	// 	}
	// 	user.Name = username
	// 	user.Password = userPassword
	// 	user.Email = userEmail
	//
	// }
  user = models.Users{
		Name: "Test",
		Password: "test",
		Email: "test@acme,com",
	}
	return user, nil
}

func (r *postgresRepository) GetAllUsers() ([]models.Users, error) {
	var users []models.Users

	allUsers, err := r.client.Query(context.Background(),"SELECT name,password,email from users")
	if err != nil {
		fmt.Printf("Error executing sql: %v", err)
		return nil, err
	}

	for allUsers.Next() {
		var userName, userPassword, userEmail string
		err := allUsers.Scan(&userName, &userPassword, &userEmail)

		if err != nil {
			fmt.Printf("Error scanning results:%v", err)
			return nil, err
		}

		users = append(users, models.Users{Name: userName, Password: userPassword, Email: userEmail})

	}

	return users, nil

}
