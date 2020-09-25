package provider

import (
	"context"
	"net/http"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/models"
)

type localProvider struct {
	providerName string
}

// This will provide the implementation of an authorization provider

// NewLocalAuthProvider creates a instance of a local auth provider
func NewLocalAuthProvider() (authorize.AuthProvider, error) {
	ap := &localProvider{
		providerName: "local",
	}

	return ap, nil

}

func (p *localProvider) GetName() string {
	return p.providerName
}

// Login implementation of the login method
func (p *localProvider) Login(username string) (string, error) {

	return "", nil
}

// RegisterUsers is an implementation of the Register users
func (p *localProvider) Register(users models.Users) error {
	return nil
}

func (p *localProvider) GetOAuthClient(code string, ctx context.Context) (*http.Client, error) {
	return nil, nil
}

func (p *localProvider) ProcessUserData(data []byte) (models.Users, error) {
	return models.Users{}, nil
}
