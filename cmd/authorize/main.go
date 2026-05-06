package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/alfrye/authorize/internal/authorization/provider"
	"github.com/alfrye/authorize/internal/authorize"
	api "github.com/alfrye/authorize/internal/handlers/api"
	"github.com/alfrye/authorize/internal/oidc"
	mg "github.com/alfrye/authorize/internal/persistence/mongo"
	mysql "github.com/alfrye/authorize/internal/persistence/mysql"
	sqlite "github.com/alfrye/authorize/internal/persistence/sqlite"
	"github.com/alfrye/authorize/internal/server"
)

func main() {
	fmt.Println("starting point for Authorize")

	port := getEnvOrDefault("PORT", "9010")
	s := server.New(port)
	repo := choseRepository()
	authProvider := choseAuthProvider()
	authService := authorize.NewAuthService(repo, authProvider)
	nhandler := api.NewAuthHandler(authService)
	s.PopulateRoutes(s.AuthorizeServiceRoutes(nhandler))

	oidcProvider, err := createOIDCProvider(repo)
	if err != nil {
		log.Printf("Warning: could not create OIDC provider: %v", err)
	} else {
		s.AttachOIDCRoutes(oidcProvider.Router())
	}

	s.Listen()
}

func createOIDCProvider(repo authorize.AuthorizeRepository) (*oidc.Provider, error) {
	issuer := getEnvOrDefault("ISSUER", "https://localhost:9010")
	dbPath := getEnvOrDefault("DB_PATH", "./data/auth.db")
	accessTokenTTL := parseDuration("JWT_EXPIRY", 5*time.Minute)
	refreshTokenTTL := parseDuration("REFRESH_EXPIRY", 30*24*time.Hour)
	authCodeTTL := parseDuration("CODE_EXPIRY", 5*time.Minute)
	sessionMaxAge := parseDuration("SESSION_MAX_AGE", 24*time.Hour)
	cacheSizeMB := parseInt("CACHE_SIZE_MB", 32)

	return oidc.NewProvider(oidc.Config{
		Issuer:          issuer,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		AuthCodeTTL:     authCodeTTL,
		SessionMaxAge:   sessionMaxAge,
		CacheSizeMB:     int64(cacheSizeMB),
		DBPath:          dbPath,
	}, repo)
}

func choseRepository() authorize.AuthorizeRepository {
	dbVendor := os.Getenv("DB_VENDOR")
	if dbVendor == "" {
		dbVendor = os.Getenv("DB_URL")
	}

	switch dbVendor {
	case "sqlite":
		dbPath := getEnvOrDefault("DB_PATH", "./data/auth.db")
		repo, err := sqlite.NewSQLiteRepository(dbPath)
		if err != nil {
			log.Fatalf("Could not create sqlite repo: %v", err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongoDB := os.Getenv("MONGO_DB")
		mongoTimeout, err := strconv.Atoi(getEnvOrDefault("DB_TIMEOUT", "3000"))
		if err != nil {
			mongoTimeout = 3000
		}
		repo, err := mg.NewMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if err != nil {
			log.Fatal("Could not create mongo repo")
		}
		return repo
	case "mysql":
		mySQLHost := os.Getenv("MYSQL_HOST")
		mySQLPort := os.Getenv("MYSQL_PORT")
		mySQLDatabase := os.Getenv("MYSQL_DB")
		mySQLUser := os.Getenv("MYSQL_USER")
		mySQLPassword := os.Getenv("MYSQL_PASS")
		mySQLTimeout, err := strconv.Atoi(getEnvOrDefault("DB_TIMEOUT", "3000"))
		if err != nil {
			mySQLTimeout = 3000
		}
		repo, err := mysql.NewMySQLRepository("mysql", mySQLHost, mySQLPort, mySQLDatabase, mySQLUser, mySQLPassword, mySQLTimeout)
		if err != nil {
			log.Fatal("Could not create mysql repo")
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
	case "google":
		authProvider, err := provider.NewGoogleAuthProvider()
		if err != nil {
			log.Fatal("Can not create authentication provider")
		}
		return authProvider
	}

	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func parseDuration(envKey string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(envKey); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

func parseInt(envKey string, defaultVal int) int {
	if val := os.Getenv(envKey); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
