package oidc

import (
	"encoding/json"
	"net/http"

	"github.com/alfrye/authorize/internal/token"
	"github.com/golang-jwt/jwt/v5"
)

func (p *Provider) handleIntrospect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeIntrospectError(w)
		return
	}

	tokenString := r.FormValue("token")
	if tokenString == "" {
		writeIntrospectError(w)
		return
	}

	parsedToken, err := p.verifier.VerifyToken(tokenString)
	if err != nil {
		writeIntrospectInactive(w)
		return
	}

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	resp := token.IntrospectionResponse{
		Active:    true,
		TokenType: "Bearer",
	}

	if sub, _ := claims["sub"].(string); sub != "" {
		resp.Sub = sub
	}
	if aud, ok := claims["aud"].([]interface{}); ok && len(aud) > 0 {
		if cid, _ := aud[0].(string); cid != "" {
			resp.ClientID = cid
		}
	}
	if exp, ok := claims["exp"].(float64); ok {
		resp.Exp = int64(exp)
	}
	if iat, ok := claims["iat"].(float64); ok {
		resp.Iat = int64(iat)
	}
	if iss, _ := claims["iss"].(string); iss != "" {
		resp.Iss = iss
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func writeIntrospectInactive(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"active": false})
}

func writeIntrospectError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]bool{"active": false})
}
