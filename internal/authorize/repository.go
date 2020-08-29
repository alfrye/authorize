package authorize

import (
	"github.com/alfrye/authorize/internal/models"
)

type (

	// AuthorizeRepository defines repository
	AuthorizeRepository interface {
		GetUser(username string) (models.Users, error)
		CreateUser(models.Users) error
	}
)
