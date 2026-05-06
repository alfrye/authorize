package authorize

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alfrye/authorize/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	AuthService struct {
		AuthRepository AuthorizeRepository
		AuthProvider   AuthProvider
		signingKey     *ecdsa.PrivateKey
		verifyingKey   *ecdsa.PublicKey
		keyID          string
		issuer         string
	}

	CustomClaims struct {
		Username string `json:"Username"`
		jwt.RegisteredClaims
	}
)

var session = map[string]string{}

func NewAuthService(repo AuthorizeRepository, provider AuthProvider) AuthService {
	signingKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate ECDSA key: %v", err)
	}
	return AuthService{
		AuthRepository: repo,
		AuthProvider:   provider,
		signingKey:     signingKey,
		verifyingKey:   &signingKey.PublicKey,
		keyID:          uuid.New().String(),
		issuer:         "https://localhost:9010",
	}
}

func (auth *AuthService) SetIssuer(issuer string) {
	auth.issuer = issuer
}

func (auth *AuthService) GetSigningKey() *ecdsa.PrivateKey {
	return auth.signingKey
}

func (auth *AuthService) GetVerifyingKey() *ecdsa.PublicKey {
	return auth.verifyingKey
}

func (auth *AuthService) GetKeyID() string {
	return auth.keyID
}

func (auth *AuthService) CreateSession(u models.Users, w http.ResponseWriter) error {
	token := auth.GenerateToken(u)
	cookie := http.Cookie{Name: "auth", Value: token, Path: "/", HttpOnly: true, Secure: true}
	http.SetCookie(w, &cookie)
	session[u.Name] = token
	return nil
}

func (auth *AuthService) GenerateToken(u models.Users) string {
	now := time.Now()
	claims := CustomClaims{
		Username: u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    auth.issuer,
			Subject:   u.ID,
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = auth.keyID
	tokenString, err := token.SignedString(auth.signingKey)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func (auth AuthService) RegisterUser(u models.Users) error {
	err := auth.AuthRepository.CreateUser(u)
	if err != nil {
		log.Println("Unable to register user in database")
		return err
	}
	log.Println("Registered new user.......")
	return nil
}

func (auth AuthService) ParseToken(t string) (string, error) {
	token, err := jwt.ParseWithClaims(t, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return auth.verifyingKey, nil
	}, jwt.WithIssuer(auth.issuer))

	if err != nil {
		return "", errors.New("Could not parse token with claims")
	}

	if !token.Valid {
		return "", errors.New("Token is not valid")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.Username, nil
	}

	return "", errors.New("Invalid token claims")
}
