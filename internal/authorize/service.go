package authorize

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alfrye/authorize/internal/models"
	"github.com/dgrijalva/jwt-go"
)

type (

	// AuthService defines the provider for authentication
	AuthService struct {
		AuthRepository AuthorizeRepository
		AuthProvider   AuthProvider
		//OAuthProvider  OAuthProvider
	}

	// CustomClaims defines the claims for the jtw token
	CustomClaims struct {
		Username string `json:"Username"`
		jwt.StandardClaims
	}
)

var (
	session = map[string]string{}
	key     = []byte("secret")
)

// NewAuthService Instaniates an Auth service
func NewAuthService(repo AuthorizeRepository, provider AuthProvider) AuthService {
	return AuthService{
		AuthRepository: repo,
		AuthProvider:   provider,
	}

}

//CreateSession creates a session and sends back a cookie
func (auth AuthService) CreateSession(u models.Users, w http.ResponseWriter) error {

	// Creates a session for the user
	// call generate token
	token := auth.GenerateToken(u)

	// Generate Cookie

	cookie := http.Cookie{Name: "auth", Value: token, Path: "/"}

	http.SetCookie(w, &cookie)
	// 	fmt.Printf("UserName:%s", result)
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("Login Succeeded"))
	session[u.Name] = token
	return nil
}

// GenerateToken generates he jwt token for the user
func (auth AuthService) GenerateToken(u models.Users) string {
	//	key := []byte("alan")
	claims := CustomClaims{
		Username: u.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func (auth AuthService) RegisterUser(u models.Users) error {

	// persist data for users
	//?? what db mongo, key/value, cockroach, sqllite rdbms

	err := auth.AuthRepository.CreateUser(u)
	if err != nil {
		log.Println("Unable to register user in database")
		return err
	}
	log.Println("Registered new user.......")
	return nil
	// Creates a session for the user
}

//ParseToken parse the jwt token
func (auth AuthService) ParseToken(t string) (string, error) {

	token, err := jwt.ParseWithClaims(t, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("different algorothm used")
		}

		return key, nil
	})

	if err != nil {
		return "", errors.New("Could not parse token with claims")
	}

	if !token.Valid {
		return "", errors.New("Token is not valid")
	}

	return token.Claims.(*CustomClaims).Username, nil
}
