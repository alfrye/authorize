package authorize

import "github.com/alfrye/authorize/internal/models"

type (

	// AuthProvider defines the provider for authentication
	AuthProvider interface {
		Login(username string) error
		Register(users models.Users) error
	}
)
