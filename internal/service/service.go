package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/alfrye/authorize/internal/authorize"
	"github.com/alfrye/authorize/internal/endpoints"
	"github.com/alfrye/authorize/internal/models"
)

type (
	// GoKitService implements the go-kit service interface
	GoKitService struct {
		authService authorize.AuthService
	}
)

// NewGoKitService creates a new go-kit service
func NewGoKitService(authService authorize.AuthService) endpoints.AuthService {
	return &GoKitService{
		authService: authService,
	}
}

// Login implements the go-kit service interface
func (s *GoKitService) Login(ctx context.Context, user models.Users) (string, error) {
	url, err := s.authService.AuthProvider.Login("")
	if err != nil {
		log.Printf("Could not get OAuth provider redirect URL: %v", err)
		return "", err
	}
	return url, nil
}

// Register implements the go-kit service interface
func (s *GoKitService) Register(ctx context.Context, user models.Users) error {
	return s.authService.RegisterUser(user)
}

// GoogleCallback implements the go-kit service interface
func (s *GoKitService) GoogleCallback(ctx context.Context, code, state string) (models.Users, error) {
	if state == "" {
		return models.Users{}, errors.New("state parameter is required")
	}
	
	if code == "" {
		return models.Users{}, errors.New("code parameter is required")
	}

	clt, err := s.authService.AuthProvider.GetOAuthClient(code, ctx)
	if err != nil {
		return models.Users{}, fmt.Errorf("error getting OAuth client: %w", err)
	}
	
	resp, err := clt.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return models.Users{}, fmt.Errorf("error getting user info from Google: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Users{}, fmt.Errorf("error reading response body: %w", err)
	}

	user, err := s.authService.AuthProvider.ProcessUserData(data)
	if err != nil {
		return models.Users{}, fmt.Errorf("error processing user data: %w", err)
	}

	// Check if user already exists
	existingUser, err := s.authService.AuthRepository.GetUser(user.Name)
	if err != nil {
		log.Printf("Error retrieving user from database: %v", err)
	}
	
	if (models.Users{}) == existingUser {
		// Register new user
		if err := s.authService.RegisterUser(user); err != nil {
			return models.Users{}, fmt.Errorf("error registering new user: %w", err)
		}
		existingUser = user
	}

	return existingUser, nil
}

// Serve implements the go-kit service interface
func (s *GoKitService) Serve(ctx context.Context, token string) (models.Users, error) {
	if token == "" {
		return models.Users{}, errors.New("token is required")
	}

	username, err := s.authService.ParseToken(token)
	if err != nil {
		return models.Users{}, fmt.Errorf("error parsing token: %w", err)
	}

	user, err := s.authService.AuthRepository.GetUser(username)
	if err != nil {
		return models.Users{}, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}