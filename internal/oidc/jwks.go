package oidc

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
)

func (p *Provider) handleJWKS(w http.ResponseWriter, r *http.Request) {
	pubKey := p.signer.PublicKey()
	jwk := map[string]interface{}{
		"kty": "EC",
		"kid": p.signer.KeyID(),
		"crv": "P-256",
		"x":   encodeBigInt(pubKey.X),
		"y":   encodeBigInt(pubKey.Y),
		"alg": "ES256",
		"use": "sig",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"keys": []map[string]interface{}{jwk},
	})
}

func encodeBigInt(n *big.Int) string {
	data := n.Bytes()
	if len(data) < 32 {
		padding := make([]byte, 32-len(data))
		data = append(padding, data...)
	}
	return base64.RawURLEncoding.EncodeToString(data)
}
