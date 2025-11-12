package persistence

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	postgresRepository struct {
		client   *pgxpool.Pool
		database string
		timeout  time.Duration
	}
)

func newPostgresClient(postgresURL string ) (*pgxpool.Pool, error) {
	ctx := context.Background()
  pool, err := pgxpool.New(ctx, postgresURL)
//	dbClient, err := pgx.Connect(context.Background(), "databaseURL")
	if err != nil {
		log.Println("Can not open database")
		return nil, err
	}
// verify the connection
  if err := pool.Ping(ctx); err != nil {
		log.Printf("unable to ping database: %v", err)
		return nil, err
	}

	return pool, nil
}

func NewPostgresRepository(dbHost, dbPort, postgresDB, postgresUser, postgresPassword string, timeout int) (authorize.AuthorizeRepository, error) {
	dbAddr := net.JoinHostPort(dbHost, dbPort)
	dbName := "postgres"
	postgresURL := "postgresql://" + postgresUser + ":" + postgresPassword + "@" + dbAddr + "/" + dbName
	// postgresConfig := postgres.Config{
	// 	User:   postgresUser,
	// 	Passwd: postgresPassword,
	// 	Net:    "tcp",
	// 	DBName: postgresDB,
	// 	Addr:   dbAddr,
	// }
  log.Printf("Connecting to database:%s",postgresURL)
 	dbClient, err := newPostgresClient(postgresURL)
	if err != nil {
		log.Printf("Error Creating client: %s", err)
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
//var user models.Users
sql := "Select * from users;"
 err := r.client.QueryRow(context.Background(), sql)
	if err != nil {
		 log.Printf("Error querying the users table:%v", err)
	 }

	//  defer rows.Close()
	//
	//  for rows.Next() {
	// 	 err := rows.Scan(
	// 		 &user.FirstName,
	// 		 &user.LastName,
	// 		 &user.Email)
	//
	// if err != nil {
	// 	 log.Printf("Error Scanning rows: %w", err)
	//  }
	//
	
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
	selUser, err := r.client.Query(context.Background(),"SELECT firstname,lastname,email from users Where firstnam=$1", username)
	if err != nil {
		fmt.Printf("Error creating sql: %v", err)
		return user, err
	}
	for selUser.Next() {
		var userFirstName,userLastName, userEmail string

		err := selUser.Scan(&userFirstName, &userLastName, &userEmail)
		if err != nil {
			fmt.Printf("Error scanning results:%v", err)
			return models.Users{}, err
		}
		user.FirstName = userFirstName
		user.LastName = userLastName
		user.Email = userEmail

	}
	//  user = models.Users{
	// 	Name: "Test",
	// 	Password: "test",
	// 	Email: "test@acme,com",
	// }
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
