package session

import "time"

type Session struct {
	ID        string
	UserID    string
	ClientID  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (s *Session) IsValid() bool {
	return time.Now().Before(s.ExpiresAt)
}
