package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (authorize.AuthorizeRepository, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}
	dsn := dbPath
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}
	db.SetMaxOpenConns(4)
	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL; PRAGMA cache_size=-20000; PRAGMA temp_store=MEMORY;"); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return &sqliteRepository{db: db}, nil
}

func runMigrations(db *sql.DB) error {
	migration, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration: %w", err)
	}
	_, err = db.Exec(string(migration))
	return err
}

func (r *sqliteRepository) GetUser(username string) (models.Users, error) {
	var u models.Users
	var createdAt time.Time
	err := r.db.QueryRow(
		"SELECT id, email, first_name, last_name, avatar_url, created_at FROM users WHERE email = ?",
		username,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.AvatarURL, &createdAt)
	u.CreatedAt = createdAt
	return u, err
}

func (r *sqliteRepository) GetAllUsers() ([]models.Users, error) {
	rows, err := r.db.Query("SELECT id, email, first_name, last_name, avatar_url FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.Users
	for rows.Next() {
		var u models.Users
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *sqliteRepository) CreateUser(user models.Users) error {
	_, err := r.db.Exec(
		"INSERT INTO users (id, email, password_hash, first_name, last_name, avatar_url) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.AvatarURL,
	)
	return err
}

func (r *sqliteRepository) GetClient(clientID string) (models.Client, error) {
	var c models.Client
	var redirectURIs, scopes, grantTypes string
	err := r.db.QueryRow(
		"SELECT id, name, secret_hash, redirect_uris, scopes, grant_types FROM clients WHERE id = ?",
		clientID,
	).Scan(&c.ID, &c.Name, &c.SecretHash, &redirectURIs, &scopes, &grantTypes)
	if err != nil {
		return c, err
	}
	json.Unmarshal([]byte(redirectURIs), &c.RedirectURIs)
	json.Unmarshal([]byte(scopes), &c.Scopes)
	json.Unmarshal([]byte(grantTypes), &c.GrantTypes)
	return c, nil
}

func (r *sqliteRepository) CreateClient(client models.Client) error {
	redirectURIs, _ := json.Marshal(client.RedirectURIs)
	scopes, _ := json.Marshal(client.Scopes)
	grantTypes, _ := json.Marshal(client.GrantTypes)
	_, err := r.db.Exec(
		"INSERT INTO clients (id, name, secret_hash, redirect_uris, scopes, grant_types) VALUES (?, ?, ?, ?, ?, ?)",
		client.ID, client.Name, client.SecretHash, string(redirectURIs), string(scopes), string(grantTypes),
	)
	return err
}

func (r *sqliteRepository) SaveAuthCode(code models.AuthCode) error {
	_, err := r.db.Exec(
		"INSERT INTO auth_codes (code, client_id, user_id, redirect_uri, code_challenge, code_challenge_method, scopes, nonce, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		code.Code, code.ClientID, code.UserID, code.RedirectURI, code.CodeChallenge, code.CodeChallengeMethod, code.Scopes, code.Nonce, code.ExpiresAt,
	)
	return err
}

func (r *sqliteRepository) GetAuthCode(code string) (models.AuthCode, error) {
	var c models.AuthCode
	var createdAt time.Time
	err := r.db.QueryRow(
		"SELECT client_id, user_id, redirect_uri, code_challenge, code_challenge_method, scopes, nonce, expires_at FROM auth_codes WHERE code = ?",
		code,
	).Scan(&c.ClientID, &c.UserID, &c.RedirectURI, &c.CodeChallenge, &c.CodeChallengeMethod, &c.Scopes, &c.Nonce, &createdAt)
	c.Code = code
	c.ExpiresAt = createdAt
	return c, err
}

func (r *sqliteRepository) DeleteAuthCode(code string) error {
	_, err := r.db.Exec("DELETE FROM auth_codes WHERE code = ?", code)
	return err
}

func (r *sqliteRepository) CreateSession(s models.Session) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (id, user_id, client_id, created_at, expires_at) VALUES (?, ?, ?, ?, ?)",
		s.ID, s.UserID, s.ClientID, s.CreatedAt, s.ExpiresAt,
	)
	return err
}

func (r *sqliteRepository) GetSession(sessionID string) (models.Session, error) {
	var s models.Session
	var createdAt, expiresAt time.Time
	err := r.db.QueryRow(
		"SELECT user_id, client_id, created_at, expires_at FROM sessions WHERE id = ?",
		sessionID,
	).Scan(&s.UserID, &s.ClientID, &createdAt, &expiresAt)
	s.ID = sessionID
	s.CreatedAt = createdAt
	s.ExpiresAt = expiresAt
	return s, err
}

func (r *sqliteRepository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

func (r *sqliteRepository) SaveRefreshToken(rt models.RefreshToken) error {
	_, err := r.db.Exec(
		"INSERT INTO refresh_tokens (token, session_id, expires_at) VALUES (?, ?, ?)",
		rt.Token, rt.SessionID, rt.ExpiresAt,
	)
	return err
}

func (r *sqliteRepository) GetRefreshToken(token string) (models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.QueryRow(
		"SELECT session_id, expires_at FROM refresh_tokens WHERE token = ?",
		token,
	).Scan(&rt.SessionID, &rt.ExpiresAt)
	rt.Token = token
	return rt, err
}

func (r *sqliteRepository) DeleteRefreshToken(token string) error {
	_, err := r.db.Exec("DELETE FROM refresh_tokens WHERE token = ?", token)
	return err
}
