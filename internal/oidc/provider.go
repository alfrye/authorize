package oidc

import (
	"time"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/session"
	"github.com/alfrye/authorize/internal/token"
	"github.com/gorilla/mux"
)

type Config struct {
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AuthCodeTTL     time.Duration
	SessionMaxAge   time.Duration
	CacheSizeMB     int64
	DBPath          string
}

type Provider struct {
	config   Config
	repo     authorize.AuthorizeRepository
	signer   *token.Signer
	verifier *token.Verifier
	sessions *session.Store
	router   *mux.Router
}

func NewProvider(cfg Config, repo authorize.AuthorizeRepository) (*Provider, error) {
	signer, err := token.NewSigner(cfg.Issuer)
	if err != nil {
		return nil, err
	}
	verifier := token.NewVerifier(signer.PublicKey(), cfg.Issuer)
	sessions, err := session.NewStore(repo, cfg.SessionMaxAge, cfg.CacheSizeMB)
	if err != nil {
		return nil, err
	}
	p := &Provider{
		config:   cfg,
		repo:     repo,
		signer:   signer,
		verifier: verifier,
		sessions: sessions,
		router:   mux.NewRouter(),
	}
	p.setupRoutes()
	return p, nil
}

func (p *Provider) Signer() *token.Signer {
	return p.signer
}

func (p *Provider) Sessions() *session.Store {
	return p.sessions
}

func (p *Provider) Router() *mux.Router {
	return p.router
}
