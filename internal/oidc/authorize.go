package oidc

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alfrye/authorize/internal/models"
	"github.com/google/uuid"
)

func (p *Provider) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.handleAuthorizeGet(w, r)
	case http.MethodPost:
		p.handleAuthorizePost(w, r)
	}
}

func (p *Provider) handleAuthorizeGet(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	scope := r.URL.Query().Get("scope")
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")
	state := r.URL.Query().Get("state")

	if err := p.validateAuthorizeRequest(clientID, redirectURI, responseType, scope, codeChallenge, codeChallengeMethod); err != nil {
		httpErrorRedirect(w, redirectURI, state, err.Error())
		return
	}
}

func (p *Provider) handleAuthorizePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	clientID := r.FormValue("client_id")
	redirectURI := r.FormValue("redirect_uri")
	codeChallenge := r.FormValue("code_challenge")
	codeChallengeMethod := r.FormValue("code_challenge_method")
	state := r.FormValue("state")
	nonce := r.FormValue("nonce")
	userID := r.FormValue("user_id")
	scope := r.FormValue("scope")

	code := uuid.New().String()
	authCode := models.AuthCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		Scopes:              scope,
		Nonce:               nonce,
		ExpiresAt:           time.Now().Add(p.config.AuthCodeTTL),
	}
	if err := p.repo.SaveAuthCode(authCode); err != nil {
		httpErrorRedirect(w, redirectURI, state, err.Error())
		return
	}

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (p *Provider) validateAuthorizeRequest(clientID, redirectURI, responseType, scope, codeChallenge, codeChallengeMethod string) error {
	if responseType != "code" {
		return fmt.Errorf("unsupported response_type")
	}
	if codeChallenge == "" {
		return fmt.Errorf("code_challenge is required (PKCE)")
	}
	if codeChallengeMethod != "" && codeChallengeMethod != "S256" {
		return fmt.Errorf("unsupported code_challenge_method")
	}
	_, err := p.repo.GetClient(clientID)
	if err != nil {
		return fmt.Errorf("invalid client_id")
	}
	return nil
}

func httpErrorRedirect(w http.ResponseWriter, redirectURI, state, errorMsg string) {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return
	}
	q := u.Query()
	q.Set("error", "invalid_request")
	q.Set("error_description", errorMsg)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, nil, u.String(), http.StatusFound)
}
