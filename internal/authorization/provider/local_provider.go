package provider

import (
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

// Login implementation of the login method
func (p *localProvider) Login(username string) error {
	return nil
}

// RegisterUsers is an implementation of the Register users
func (p *localProvider) Register(users models.Users) error {
	return nil
}
