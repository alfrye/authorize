package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/server"
	"github.com/spf13/cobra"

 "github.com/alfrye/authorize/internal/authorization/provider"

	mg  "github.com/alfrye/authorize/internal/persistence/mongo"
	mysql "github.com/alfrye/authorize/internal/persistence/mysql"
	postgres "github.com/alfrye/authorize/internal/persistence/postgres"
	 api "github.com/alfrye/authorize/internal/handlers/api"
)

type config struct {
	 db_vendor string
	 port int
	 db_host string
	 db_user string
	 db_name string
	 db_password string
db_port string 
	 db_timeout int
 }
 
var (
	startCmd = &cobra.Command{
		Use: "start",
		Short: "Start the service",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// The runE field is what is executed for the commnd
			fmt.Println("Start command was called ")
			c := config{}	
			c.db_vendor, err = cmd.Flags().GetString("dbvendor")			
			if err != nil {
				return err
			}
			c.port, err = cmd.Flags().GetInt("port")
			if err != nil {
				 return err
		 }

      c.db_name, err = cmd.Flags().GetString("dbName")
			if err != nil {
				 return err
			 }
			 
      c.db_port, err = cmd.Flags().GetString("dbPort")
			if err != nil {
				 return err
			 }
      c.db_user, err = cmd.Flags().GetString("dbUser")

			if err != nil {
				 return err
			 }
			 
      c.db_password, err = cmd.Flags().GetString("dbPassword")

			if err != nil {
				 return err
			 }



		c.db_timeout, err = cmd.Flags().GetInt("dbTimeout")

			if err != nil {
				 return err
			 }

     		  setup(c)
			return nil
		},
}
)

func init() {
  rootCmd.AddCommand(startCmd)
	startCmd.Flags().Int("port", 9010, "Server Port")
	startCmd.Flags().String("dbvendor", "postgresql", "Database vendor")
	startCmd.Flags().String("dbHost", "localhost", "Database Host")
	startCmd.Flags().String("dbPort", "5432", "Database Port")
	startCmd.Flags().String("dbName", "", "Database Name")
	startCmd.Flags().String("dbUser", "", "Database User")
	startCmd.Flags().String("dbPassword", "" , "Database Password")
	startCmd.Flags().Int("dbTimeout", 3000, "Database Timeout")
}

func setup(c config) error {
  // Connect to database and start the server
	fmt.Println("Connecting to database and start server")
	s := server.New(c.port)

	repo := choseRepository(c)
	authProvider := choseAuthProvider()
	authService := authorize.NewAuthService(repo, authProvider)
	nhandler := api.NewAuthHandler(authService)
	s.PopulateRoutes(s.AuthorizeServiceRoutes(nhandler))
	s.Listen()
 return nil
}

func choseRepository(c config) authorize.AuthorizeRepository {

	switch c.db_vendor {
	case "mongo":
		//Set up information for mongo db
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
		if err != nil {

		}
		repo, err := mg.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal("Could not create mongo repo")
		}
		return repo
	case "mysql":
		//Setup information for mysql database
		//	mySQLURL := os.Getenv("MYSQL_URL")
		mySQLHost := os.Getenv("MYSQL_HOST")
		mySQLPort := os.Getenv("MYSQL_PORT")
		mySQLDatabase := os.Getenv("MYSQL_DB")
		mySQLUser := os.Getenv("MYSQL_USER")
		mySQLPassword := os.Getenv("MYSQL_PASS")

		mySQLTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
		if err != nil {

		}

		repo, err := mysql.NewMySQLRepository("mysql", mySQLHost, mySQLPort, mySQLDatabase, mySQLUser, mySQLPassword, mySQLTimeout)
		if err != nil {
			log.Fatal("Could not create mongo repo")
		}
		return repo
	case "postgres":
		//Setup information for postgres database
		log.Printf("Using Database %s", c.db_vendor)
	   repo, err := postgres.NewPostgresRepository(c.db_host, c.db_port,c.db_name,c.db_user,c.db_password, c.db_timeout)
	  if err != nil {
			 log.Fatal("Could not create postgres repo")
		}
	return repo
  default:
	   log.Println("Database is not supported")
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
	case "google":
		authProvider, err := provider.NewGoogleAuthProvider()
		if err != nil {
			log.Fatal("Can not create authentication provider")
		}
		return authProvider
	}

	return nil
}
