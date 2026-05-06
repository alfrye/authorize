package token

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	publicKey *ecdsa.PublicKey
	issuer    string
}

func NewVerifier(publicKey *ecdsa.PublicKey, issuer string) *Verifier {
	return &Verifier{
		publicKey: publicKey,
		issuer:    issuer,
	}
}

func (v *Verifier) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.publicKey, nil
	}, jwt.WithIssuer(v.issuer))
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return token, nil
}
