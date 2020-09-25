package authorize

import (
	"context"
	"net/http"

	"github.com/alfrye/authorize/internal/models"
)

type (

	// AuthProvider defines the provider for authentication
	AuthProvider interface {
		Login(username string) (string, error)
		Register(users models.Users) error
		GetOAuthClient(code string, ctx context.Context) (*http.Client, error)
		ProcessUserData(data []byte) (models.Users, error)
		GetName() string
	}
	// OAuthProvider defines the provider for OAuth Authentication
	OAuthProvider interface {
		GetOAuthClient(code string, ctx context.Context) (*http.Client, error)
		ProcessUserData(data []byte) (models.Users, error)
	}
)
