package oidc

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alfrye/authorize/internal/models"
	"github.com/google/uuid"
)

func (p *Provider) handleToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeTokenError(w, "invalid_request", "could not parse form")
		return
	}

	grantType := r.FormValue("grant_type")
	switch grantType {
	case "authorization_code":
		p.handleAuthorizationCode(w, r)
	case "refresh_token":
		p.handleRefreshToken(w, r)
	default:
		writeTokenError(w, "unsupported_grant_type", "grant_type must be authorization_code or refresh_token")
	}
}

func (p *Provider) handleAuthorizationCode(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	codeVerifier := r.FormValue("code_verifier")
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")

	authCode, err := p.repo.GetAuthCode(code)
	if err != nil || time.Now().After(authCode.ExpiresAt) {
		writeTokenError(w, "invalid_grant", "invalid or expired authorization code")
		return
	}
	if authCode.ClientID != clientID || authCode.RedirectURI != redirectURI {
		writeTokenError(w, "invalid_grant", "client_id or redirect_uri mismatch")
		return
	}
	if !verifyPKCE(codeVerifier, authCode.CodeChallenge) {
		writeTokenError(w, "invalid_grant", "code_verifier mismatch")
		return
	}
	if err := p.repo.DeleteAuthCode(code); err != nil {
		writeTokenError(w, "server_error", "failed to delete auth code")
		return
	}

	session, err := p.sessions.Create(authCode.UserID, authCode.ClientID)
	if err != nil {
		writeTokenError(w, "server_error", "failed to create session")
		return
	}

	idToken, err := p.signer.SignIDToken(authCode.UserID, []string{clientID}, authCode.Nonce, time.Now(), p.config.AccessTokenTTL)
	if err != nil {
		writeTokenError(w, "server_error", "failed to sign ID token")
		return
	}

	accessToken, err := p.signer.SignAccessToken(authCode.UserID, []string{clientID}, authCode.Scopes, nil, nil, session.ID, p.config.AccessTokenTTL)
	if err != nil {
		writeTokenError(w, "server_error", "failed to sign access token")
		return
	}

	refreshToken := uuid.New().String()
	if err := p.repo.SaveRefreshToken(models.RefreshToken{
		Token:     refreshToken,
		SessionID: session.ID,
		ExpiresAt: time.Now().Add(p.config.RefreshTokenTTL),
	}); err != nil {
		writeTokenError(w, "server_error", "failed to save refresh token")
		return
	}

	writeTokenResponse(w, accessToken, idToken, refreshToken, int(p.config.AccessTokenTTL.Seconds()), authCode.Scopes)
}

func (p *Provider) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refresh_token")
	clientID := r.FormValue("client_id")

	rt, err := p.repo.GetRefreshToken(refreshToken)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		writeTokenError(w, "invalid_grant", "invalid or expired refresh token")
		return
	}

	session, err := p.sessions.Get(rt.SessionID)
	if err != nil {
		writeTokenError(w, "invalid_grant", "session not found")
		return
	}

	if err := p.repo.DeleteRefreshToken(refreshToken); err != nil {
		writeTokenError(w, "server_error", "failed to delete refresh token")
		return
	}

	idToken, _ := p.signer.SignIDToken(session.UserID, []string{clientID}, "", time.Now(), p.config.AccessTokenTTL)
	accessToken, _ := p.signer.SignAccessToken(session.UserID, []string{clientID}, "openid profile email", nil, nil, session.ID, p.config.AccessTokenTTL)

	newRefreshToken := uuid.New().String()
	p.repo.SaveRefreshToken(models.RefreshToken{
		Token:     newRefreshToken,
		SessionID: session.ID,
		ExpiresAt: time.Now().Add(p.config.RefreshTokenTTL),
	})

	writeTokenResponse(w, accessToken, idToken, newRefreshToken, int(p.config.AccessTokenTTL.Seconds()), "openid profile email")
}

func verifyPKCE(codeVerifier, codeChallenge string) bool {
	hash := sha256.Sum256([]byte(codeVerifier))
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	return encoded == codeChallenge
}

func writeTokenResponse(w http.ResponseWriter, accessToken, idToken, refreshToken string, expiresIn int, scope string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
		"id_token":      idToken,
		"refresh_token": refreshToken,
		"scope":         scope,
	})
}

func writeTokenError(w http.ResponseWriter, errorType, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":             errorType,
		"error_description": description,
	})
}
