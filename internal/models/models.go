package models

import "time"

type Users struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	AvatarURL string    `json:"avatarurl"`
	UserID    string    `json:"userid"`
	CreatedAt time.Time `json:"created_at"`
}

type Client struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	SecretHash   string   `json:"-"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       []string `json:"scopes"`
	GrantTypes   []string `json:"grant_types"`
}

type Session struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ClientID    string    `json:"client_id"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	MFAVerified bool      `json:"mfa_verified"`
}

type AuthCode struct {
	Code                string    `json:"-"`
	ClientID            string    `json:"client_id"`
	UserID              string    `json:"user_id"`
	RedirectURI         string    `json:"redirect_uri"`
	CodeChallenge       string    `json:"-"`
	CodeChallengeMethod string    `json:"code_challenge_method"`
	Scopes              string    `json:"scopes"`
	Nonce               string    `json:"nonce,omitempty"`
	ExpiresAt           time.Time `json:"expires_at"`
}

type RefreshToken struct {
	Token     string    `json:"-"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
