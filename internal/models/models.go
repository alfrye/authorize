package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Users in the system
type Users struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CustomClaims defines the claims for the jtw token
type CustomClaims struct {
	Username string `json:"Username"`
	jwt.StandardClaims
}

// GenerateToken generates he jwt token for the user
func (u Users) GenerateToken() string {
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
