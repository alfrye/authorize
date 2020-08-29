package authorize

type (

	// AuthService defines the provider for authentication
	AuthService struct {
		AuthRepository AuthorizeRepository
		AuthProvider   AuthProvider
	}
)

// NewAuthService Instaniates an Auth service
func NewAuthService(repo AuthorizeRepository, provider AuthProvider) AuthService {
	return AuthService{
		AuthRepository: repo,
		AuthProvider:   provider,
	}

}
