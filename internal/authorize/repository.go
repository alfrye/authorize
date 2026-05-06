package authorize

import "github.com/alfrye/authorize/internal/models"

type AuthorizeRepository interface {
	GetUser(username string) (models.Users, error)
	GetAllUsers() ([]models.Users, error)
	CreateUser(models.Users) error
	GetClient(clientID string) (models.Client, error)
	CreateClient(client models.Client) error
	SaveAuthCode(models.AuthCode) error
	GetAuthCode(code string) (models.AuthCode, error)
	DeleteAuthCode(code string) error
	CreateSession(models.Session) error
	GetSession(sessionID string) (models.Session, error)
	DeleteSession(sessionID string) error
	SaveRefreshToken(models.RefreshToken) error
	GetRefreshToken(token string) (models.RefreshToken, error)
	DeleteRefreshToken(token string) error
}
