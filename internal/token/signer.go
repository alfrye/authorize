package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Signer struct {
	privateKey *ecdsa.PrivateKey
	keyID      string
	issuer     string
}

func NewSigner(issuer string) (*Signer, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}
	return &Signer{
		privateKey: privateKey,
		keyID:      uuid.New().String(),
		issuer:     issuer,
	}, nil
}

func (s *Signer) SignIDToken(subject string, audience []string, nonce string, authTime time.Time, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   subject,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		Nonce:    nonce,
		AuthTime: authTime,
		ACR:      "urn:acr:basic",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = s.keyID
	return token.SignedString(s.privateKey)
}

func (s *Signer) SignAccessToken(subject string, audience []string, scopes string, roles []string, permissions []string, sessionID string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   subject,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		Roles:       roles,
		Permissions: permissions,
		SID:         sessionID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = s.keyID
	return token.SignedString(s.privateKey)
}

func (s *Signer) PublicKey() *ecdsa.PublicKey {
	return &s.privateKey.PublicKey
}

func (s *Signer) KeyID() string {
	return s.keyID
}

func (s *Signer) Issuer() string {
	return s.issuer
}
