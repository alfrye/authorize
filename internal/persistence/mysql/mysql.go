package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
	"github.com/go-sql-driver/mysql"
)

type (
	mysqlRepository struct {
		client   *sql.DB
		database string
		timeout  time.Duration
	}
)

func newMySQLClient(mySQLURL string, mySQLConfig mysql.Config) (*sql.DB, error) {

	dbClient, err := sql.Open(mySQLURL, mySQLConfig.FormatDSN())
	if err != nil {
		fmt.Println("Can not open database")
		return nil, err

	}

	return dbClient, nil
}

func NewMySQLRepository(mySQLURL, dbHost, dbPort, mysqlDB, mysqlUser, mysqlPassword string, timeout int) (authorize.AuthorizeRepository, error) {
	dbAddr := net.JoinHostPort(dbHost, dbPort)
	mysqlConfig := mysql.Config{
		User:   mysqlUser,
		Passwd: mysqlPassword,
		Net:    "tcp",
		DBName: mysqlDB,
		Addr:   dbAddr,
	}

	dbClient, err := newMySQLClient(mySQLURL, mysqlConfig)
	if err != nil {
		fmt.Printf("Error Creating client: %s", err)
		return nil, err
	}
	repo := &mysqlRepository{
		client:   dbClient,
		database: mysqlDB,
		timeout:  time.Duration(timeout) * time.Second,
	}

	return repo, nil
}

func (r *mysqlRepository) CreateUser(user models.Users) error {

	insUser, err := r.client.Prepare("INSERT into users(name, password,email) Values(?,?,?)")
	if err != nil {
		fmt.Printf("Error creating sql: %v", err)
		return err
	}

	res, err := insUser.Exec(user.Name, user.Password, user.Email)
	if err != nil {
		fmt.Printf("Error creating sql: %v", err)
		return err
	}

	fmt.Println(res)

	return nil
}

func (r *mysqlRepository) GetUser(username string) (models.Users, error) {
	var user models.Users
	selUser, err := r.client.Query("SELECT name,password,email from users Where name=?", username)
	if err != nil {
		fmt.Printf("Error creating sql: %v", err)
		return user, err
	}
	for selUser.Next() {
		var userName, userPassword, userEmail string

		err := selUser.Scan(&userName, &userPassword, &userEmail)
		if err != nil {
			fmt.Printf("Error scanning results:%v", err)
			return models.Users{}, err
		}
		user.Name = username
		user.Password = userPassword
		user.Email = userEmail

	}

	return user, nil
}

func (r *mysqlRepository) GetAllUsers() ([]models.Users, error) {
	var users []models.Users

	allUsers, err := r.client.Query("SELECT name,password,email from users")
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

func (r *mysqlRepository) GetClient(clientID string) (models.Client, error) {
	return models.Client{}, errors.New("not implemented")
}

func (r *mysqlRepository) CreateClient(client models.Client) error {
	return errors.New("not implemented")
}

func (r *mysqlRepository) SaveAuthCode(code models.AuthCode) error {
	return errors.New("not implemented")
}

func (r *mysqlRepository) GetAuthCode(code string) (models.AuthCode, error) {
	return models.AuthCode{}, errors.New("not implemented")
}

func (r *mysqlRepository) DeleteAuthCode(code string) error {
	return errors.New("not implemented")
}

func (r *mysqlRepository) CreateSession(s models.Session) error {
	return errors.New("not implemented")
}

func (r *mysqlRepository) GetSession(sessionID string) (models.Session, error) {
	return models.Session{}, errors.New("not implemented")
}

func (r *mysqlRepository) DeleteSession(sessionID string) error {
	return errors.New("not implemented")
}

func (r *mysqlRepository) SaveRefreshToken(rt models.RefreshToken) error {
	return errors.New("not implemented")
}

func (r *mysqlRepository) GetRefreshToken(token string) (models.RefreshToken, error) {
	return models.RefreshToken{}, errors.New("not implemented")
}

func (r *mysqlRepository) DeleteRefreshToken(token string) error {
	return errors.New("not implemented")
}
