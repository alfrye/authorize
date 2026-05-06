package session

import (
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/google/uuid"
	"github.com/alfrye/authorize/internal/models"
)

type Store struct {
	cache  *ristretto.Cache
	repo   SessionRepository
	maxAge time.Duration
}

type SessionRepository interface {
	CreateSession(s models.Session) error
	GetSession(id string) (models.Session, error)
	DeleteSession(id string) error
}

func NewStore(repo SessionRepository, maxAge time.Duration, cacheSizeMB int64) (*Store, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     cacheSizeMB << 20,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	return &Store{
		cache:  cache,
		repo:   repo,
		maxAge: maxAge,
	}, nil
}

func (s *Store) Create(userID, clientID string) (*Session, error) {
	session := &Session{
		ID:        generateID(),
		UserID:    userID,
		ClientID:  clientID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.maxAge),
	}
	modelSession := models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		ClientID:  session.ClientID,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}
	if err := s.repo.CreateSession(modelSession); err != nil {
		return nil, err
	}
	s.cache.Set(session.ID, session, 1)
	return session, nil
}

func (s *Store) Get(id string) (*Session, error) {
	if val, found := s.cache.Get(id); found {
		return val.(*Session), nil
	}
	modelSession, err := s.repo.GetSession(id)
	if err != nil {
		return nil, err
	}
	session := &Session{
		ID:        modelSession.ID,
		UserID:    modelSession.UserID,
		ClientID:  modelSession.ClientID,
		CreatedAt: modelSession.CreatedAt,
		ExpiresAt: modelSession.ExpiresAt,
	}
	s.cache.Set(session.ID, session, 1)
	return session, nil
}

func (s *Store) Delete(id string) error {
	s.cache.Del(id)
	return s.repo.DeleteSession(id)
}

func generateID() string {
	return fmt.Sprintf("sess_%s", uuid.New().String())
}
